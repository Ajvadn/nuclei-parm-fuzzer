package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/Ajvadn/parm-fuzzer/pkg/engine"
	"github.com/Ajvadn/parm-fuzzer/pkg/providers"
	"github.com/Ajvadn/parm-fuzzer/pkg/ui"
	"github.com/Ajvadn/parm-fuzzer/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	targetDomain string
	targetFile   string
	concurrency  int
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "parm-fuzzer",
		Short: "Aggressive Parameter & JS Discovery Engine v2.0 Enterprise",
		Long:  ui.Banner + "\nA high-performance URL discovery and parameter gathering tool optimized for speed and security.",
		Run: func(cmd *cobra.Command, args []string) {
			if targetDomain == "" && targetFile == "" {
				cmd.Help()
				return
			}
			startDiscovery()
		},
	}

	rootCmd.Flags().StringVarP(&targetDomain, "domain", "d", "", "Target single domain")
	rootCmd.Flags().StringVarP(&targetFile, "file", "f", "", "File containing list of domains")
	rootCmd.Flags().IntVarP(&concurrency, "concurrency", "c", 10, "Number of concurrent providers")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func startDiscovery() {
	ui.PrintBanner()

	// Validation
	if targetDomain != "" {
		domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
		if !domainRegex.MatchString(targetDomain) {
			ui.Log("ERROR", "Invalid domain format")
			os.Exit(1)
		}
	}

	e := engine.NewEngine(concurrency)
	providers.RegisterAll(e)

	target := ""
	isFile := false
	outputName := ""

	if targetDomain != "" {
		target = targetDomain
		outputName = utils.SanitizeFilename(target)
	} else {
		target = targetFile
		isFile = true
		outputName = utils.SanitizeFilename(strings.TrimSuffix(filepath.Base(target), filepath.Ext(target)))
	}

	ui.Log("INFO", fmt.Sprintf("Starting discovery for: %s", target))
	ui.NewProgressBar(len(e.Providers))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results := e.Run(ctx, target, isFile)

	var allURLs []string
	var allJS []string
	var wg sync.WaitGroup
	var mu sync.Mutex

	// JS Discovery Worker Pool
	jsChan := make(chan []string, len(e.Providers))
	wg.Add(1)
	go func() {
		defer wg.Done()
		for urls := range jsChan {
			// Extract JS from current batch
			batchJS := extractJS(urls)
			
			// Deep extraction via subjs
			if len(urls) > 0 {
				tmpFile, _ := os.CreateTemp("", "fuzzer_batch_*.txt")
				utils.WriteLines(urls, tmpFile.Name())
				deepJS, _ := providers.ExecCommand(ctx, "subjs", []string{"-c", "100"}, tmpFile.Name())
				os.Remove(tmpFile.Name())
				
				mu.Lock()
				allJS = append(allJS, batchJS...)
				allJS = append(allJS, deepJS...)
				mu.Unlock()
			}
		}
	}()

	// Collect results
	for res := range results {
		if ui.GlobalBar != nil {
			ui.GlobalBar.Add(1)
			mu.Lock()
			count := len(allURLs)
			mu.Unlock()
			ui.GlobalBar.Describe(fmt.Sprintf("[cyan][reset] Found %d URLs via %s...", count, res.ProviderName))
		}

		if res.Err != nil {
			ui.Log("WARN", fmt.Sprintf("Provider %s failed: %v", res.ProviderName, res.Err))
		} else {
			if len(res.URLs) > 0 {
				ui.Log("SUCCESS", fmt.Sprintf("Provider %s found %d URLs", res.ProviderName, len(res.URLs)))
				mu.Lock()
				allURLs = append(allURLs, res.URLs...)
				mu.Unlock()
				
				// Send to JS discovery worker
				jsChan <- res.URLs
			}
		}
	}
	close(jsChan)
	wg.Wait()

	uniqueURLs := utils.Unique(allURLs)
	uniqueJS := utils.Unique(allJS)

	ui.Log("INFO", fmt.Sprintf("Total unique URLs discovered: %d", len(uniqueURLs)))
	ui.Log("SUCCESS", fmt.Sprintf("Total unique JS files found: %d", len(uniqueJS)))

	if len(uniqueURLs) == 0 {
		ui.Log("WARN", "No URLs found. Exiting.")
		return
	}

	// Save results
	fullFile := "full-urls-" + outputName + ".txt"
	jsFile := "js-files-" + outputName + ".txt"
	utils.WriteLines(uniqueURLs, fullFile)
	utils.WriteLines(uniqueJS, jsFile)

	ui.Log("SUCCESS", "Full URL results saved to "+fullFile)
	ui.Log("SUCCESS", "JS file results saved to "+jsFile)
	ui.Log("INFO", "Redesign v2.0 Enterprise completed successfully!")
}

func extractJS(urls []string) []string {
	var js []string
	jsRegex := regexp.MustCompile(`\.js(\?|$)`)
	for _, u := range urls {
		if jsRegex.MatchString(u) {
			js = append(js, u)
		}
	}
	return js
}
