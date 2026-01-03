package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	targetDomain string
	targetFile   string
	updateTools  bool
)

func init() {
	// Add Go bin and local bin to PATH, similar to the bash script
	homeDir, _ := os.UserHomeDir()
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		goPath = filepath.Join(homeDir, "go")
	}
	
	newPath := fmt.Sprintf("%s%c%s%c%s", 
		os.Getenv("PATH"), 
		os.PathListSeparator, 
		filepath.Join(goPath, "bin"), 
		os.PathListSeparator, 
		filepath.Join(homeDir, ".local", "bin"))
	
	os.Setenv("PATH", newPath)
}

func main() {
	flag.StringVar(&targetDomain, "d", "", "Target single domain")
	flag.StringVar(&targetDomain, "domain", "", "Target single domain")
	flag.StringVar(&targetFile, "f", "", "File containing list of domains")
	flag.StringVar(&targetFile, "file", "", "File containing list of domains")
	flag.BoolVar(&updateTools, "u", false, "Update all tools and Nuclei templates")
	flag.BoolVar(&updateTools, "update", false, "Update all tools and Nuclei templates")

	flag.Usage = func() {
		PrintBanner()
		fmt.Printf("Usage: %s [options]\n\n", os.Args[0])
		fmt.Println("Options:")
		flag.PrintDefaults()
		fmt.Println("\nExample:")
		fmt.Printf("  %s -d example.com\n", os.Args[0])
		fmt.Printf("  %s -f domains.txt\n", os.Args[0])
		fmt.Printf("  %s -u\n", os.Args[0])
	}

	flag.Parse()

	if updateTools {
		UpdateAllTools()
		os.Exit(0)
	}

	if targetDomain == "" && targetFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	if targetDomain != "" && targetFile != "" {
		LogError("Please specify either (-d) or (-f), not both.")
		os.Exit(1)
	}

	PrintBanner()
	CheckDependencies()

	outputName := ""
	if targetDomain != "" {
		outputName = targetDomain
	} else {
		outputName = strings.TrimSuffix(filepath.Base(targetFile), filepath.Ext(targetFile))
	}

	outputDir := outputName
	CreateDir(outputDir)

	fullURLsFile := filepath.Join(outputDir, "full-url-"+outputName+".txt")
	paramURLsFile := filepath.Join(outputDir, "param-url-"+outputName+".txt")
	jsURLsFile := filepath.Join(outputDir, "js-urls-"+outputName+".txt")
	nucleiResultsFile := filepath.Join(outputDir, "nuclei_results.txt")

	ctx := context.Background()

	LogInfo("Fetching URLs using multiple tools (gau, waybackurls, katana, paramspider, hakrawler, waymore) in parallel...")

	var commands []string
	if targetFile != "" {
		commands = []string{
			"gau --subs",
			"waybackurls",
			fmt.Sprintf("katana -list %s -d 5 -silent -jc -concurrency 50 -timeout 10", targetFile),
			fmt.Sprintf("paramspider -l %s -s", targetFile),
			"hakrawler",
		}
		// Note: waymore logic might need adjustment for file input if it doesn't support stdin easily in this way
	} else {
		commands = []string{
			fmt.Sprintf("echo %s | gau --subs", targetDomain),
			fmt.Sprintf("echo %s | waybackurls", targetDomain),
			fmt.Sprintf("katana -u https://%s -d 5 -silent -jc -concurrency 50 -timeout 10", targetDomain),
			fmt.Sprintf("paramspider -d %s -s > /dev/null 2>&1 && cat results/%s.txt", targetDomain, targetDomain),
			fmt.Sprintf("echo https://%s | hakrawler -d 2 -subs -u", targetDomain),
			fmt.Sprintf("waymore -i %s -mode U -oU /tmp/waymore_%s.txt > /dev/null 2>&1 && cat /tmp/waymore_%s.txt", targetDomain, targetDomain, targetDomain),
		}
	}

	var rawURLs []string
	if targetFile != "" {
		rawURLs = RunParallel(ctx, commands, targetFile)
	} else {
		rawURLs = RunParallel(ctx, commands, "")
	}

	// Deduplicate
	uniqueURLs := make(map[string]bool)
	for _, u := range rawURLs {
		if strings.TrimSpace(u) != "" {
			uniqueURLs[u] = true
		}
	}

	var filtered []string
	scopeRegex := ""
	if targetDomain != "" {
		scopeRegex = regexp.QuoteMeta(targetDomain)
	}

	for u := range uniqueURLs {
		if scopeRegex != "" {
			if matched, _ := regexp.MatchString(scopeRegex, u); !matched {
				continue
			}
		} else if targetFile != "" {
			// Basic filtering for file input: check if any domain in file is in URL
			// This is a bit slow in Go if we load file every time, so let's load once
			domains, _ := ReadLines(targetFile)
			found := false
			for _, d := range domains {
				if strings.Contains(u, d) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		filtered = append(filtered, u)
	}

	LogInfo(fmt.Sprintf("Checking for live URLs on %d discovered URLs...", len(filtered)))
	
	// Create temp file for httpx input
	tmpFile, _ := os.CreateTemp("", "fuzzer_urls_*.txt")
	defer os.Remove(tmpFile.Name())
	WriteLines(filtered, tmpFile.Name())

	liveURLs, _ := RunCommand(ctx, fmt.Sprintf("httpx -silent -threads 500 -rl 300 -timeout 5 < %s", tmpFile.Name()), nil)
	
	LogInfo("Categorizing live URLs...")
	
	var full, js, params []string
	full = liveURLs
	
	jsRegex := regexp.MustCompile(`\.js(\?|$)`)
	paramRegex := regexp.MustCompile(`\?[^=]+=.+$`)
	
	for _, u := range liveURLs {
		if jsRegex.MatchString(u) {
			js = append(js, u)
		}
	}
	
	// For parameters, use uro
	if len(liveURLs) > 0 {
		paramTmp, _ := os.CreateTemp("", "fuzzer_params_*.txt")
		defer os.Remove(paramTmp.Name())
		
		var liveWithParams []string
		for _, u := range liveURLs {
			if paramRegex.MatchString(u) {
				liveWithParams = append(liveWithParams, u)
			}
		}
		WriteLines(liveWithParams, paramTmp.Name())
		params, _ = RunCommand(ctx, fmt.Sprintf("uro < %s", paramTmp.Name()), nil)
	}

	WriteLines(full, fullURLsFile)
	WriteLines(js, jsURLsFile)
	WriteLines(params, paramURLsFile)

	LogStatus(fmt.Sprintf("Saved Full Live URLs: %s", fullURLsFile))
	LogStatus(fmt.Sprintf("Saved JS Live URLs: %s", jsURLsFile))
	LogStatus(fmt.Sprintf("Saved Parameter Live URLs: %s", paramURLsFile))

	if len(params) > 0 {
		LogInfo("Running nuclei for DAST scanning on parameter URLs...")
		nucleiCmd := fmt.Sprintf("nuclei -dast -retries 2 -silent -concurrency 50 -rate-limit 100 -o %s", nucleiResultsFile)
		
		pTmp, _ := os.CreateTemp("", "nuclei_input_*.txt")
		defer os.Remove(pTmp.Name())
		WriteLines(params, pTmp.Name())
		
		cmd := exec.Command("sh", "-c", nucleiCmd+" < "+pTmp.Name())
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	} else {
		LogInfo("No parameter URLs found for DAST scanning.")
	}

	LogInfo("Automation completed successfully!")
	
	// Final check on results
	if _, err := os.Stat(nucleiResultsFile); os.IsNotExist(err) || IsEmptyFile(nucleiResultsFile) {
		LogInfo("No vulnerable URLs found.")
	} else {
		LogInfo(fmt.Sprintf("Vulnerabilities were detected. Check %s for details.", nucleiResultsFile))
	}
}

func IsEmptyFile(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return true
	}
	return info.Size() == 0
}

func CheckDependencies() {
	tools := map[string]string{
		"gau":         "go install github.com/lc/gau/v2/cmd/gau@latest",
		"waybackurls": "go install github.com/tomnomnom/waybackurls@latest",
		"katana":      "go install github.com/projectdiscovery/katana/cmd/katana@latest",
		"httpx":       "go install github.com/projectdiscovery/httpx/cmd/httpx@latest",
		"nuclei":      "go install github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest",
		"uro":         "pip3 install uro --break-system-packages",
		"paramspider": "pip3 install git+https://github.com/devanshbatham/ParamSpider --break-system-packages",
		"waymore":     "pip3 install git+https://github.com/xnl-h4ck3r/waymore.git --break-system-packages",
		"hakrawler":   "go install github.com/hakluke/hakrawler@latest",
	}

	for tool, install := range tools {
		CheckAndInstall(tool, install)
	}
}

func UpdateAllTools() {
	LogInfo("Updating all tools and Nuclei templates...")
	
	updates := []struct {
		name string
		cmd  string
	}{
		{"gau", "go install github.com/lc/gau/v2/cmd/gau@latest"},
		{"waybackurls", "go install github.com/tomnomnom/waybackurls@latest"},
		{"katana", "go install github.com/projectdiscovery/katana/cmd/katana@latest"},
		{"httpx", "go install github.com/projectdiscovery/httpx/cmd/httpx@latest"},
		{"nuclei", "go install github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest"},
		{"nuclei templates", "nuclei -ut"},
		{"uro", "pip3 install --upgrade uro --break-system-packages"},
		{"paramspider", "pip3 install --upgrade git+https://github.com/devanshbatham/ParamSpider --break-system-packages"},
		{"waymore", "pip3 install --upgrade git+https://github.com/xnl-h4ck3r/waymore.git --break-system-packages"},
		{"hakrawler", "go install github.com/hakluke/hakrawler@latest"},
	}

	for _, up := range updates {
		LogStatus(fmt.Sprintf("Updating %s...", up.name))
		cmd := exec.Command("sh", "-c", up.cmd)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}

	LogInfo("All tools and templates updated successfully!")
}
