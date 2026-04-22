package providers

import (
	"context"
	"os"
	"github.com/Ajvadn/parm-fuzzer/pkg/engine"
	"github.com/Ajvadn/parm-fuzzer/pkg/utils"
)

type GauProvider struct{}

func (p *GauProvider) Name() string { return "gau" }

func (p *GauProvider) Execute(ctx context.Context, target string, isFile bool) ([]string, error) {
	args := []string{"--subs", "--threads", "100"}
	stdinPath := ""
	if isFile {
		stdinPath = target
	} else {
		tmp, _ := os.CreateTemp("", "gau_stdin_*.txt")
		defer os.Remove(tmp.Name())
		os.WriteFile(tmp.Name(), []byte(target), 0644)
		stdinPath = tmp.Name()
	}
	return ExecCommand(ctx, "gau", args, stdinPath)
}

type WaybackProvider struct{}

func (p *WaybackProvider) Name() string { return "waybackurls" }

func (p *WaybackProvider) Execute(ctx context.Context, target string, isFile bool) ([]string, error) {
	stdinPath := ""
	if isFile {
		stdinPath = target
	} else {
		tmp, _ := os.CreateTemp("", "wayback_stdin_*.txt")
		defer os.Remove(tmp.Name())
		os.WriteFile(tmp.Name(), []byte(target), 0644)
		stdinPath = tmp.Name()
	}
	return ExecCommand(ctx, "waybackurls", []string{}, stdinPath)
}

type KatanaProvider struct{}

func (p *KatanaProvider) Name() string { return "katana" }

func (p *KatanaProvider) Execute(ctx context.Context, target string, isFile bool) ([]string, error) {
	args := []string{"-d", "5", "-silent", "-jc", "-jsl", "-kf", "all", "-aff", "-pc", "-fsu", "-concurrency", "100"}
	if isFile {
		args = append(args, "-list", target)
	} else {
		args = append(args, "-u", "https://"+target)
	}
	return ExecCommand(ctx, "katana", args, "")
}

type ParamSpiderProvider struct{}

func (p *ParamSpiderProvider) Name() string { return "paramspider" }

func (p *ParamSpiderProvider) Execute(ctx context.Context, target string, isFile bool) ([]string, error) {
	var args []string
	if isFile {
		args = []string{"-l", target, "-s"}
	} else {
		args = []string{"-d", target, "-s"}
	}
	return ExecCommand(ctx, "paramspider", args, "")
}

type HakrawlerProvider struct{}

func (p *HakrawlerProvider) Name() string { return "hakrawler" }

func (p *HakrawlerProvider) Execute(ctx context.Context, target string, isFile bool) ([]string, error) {
	args := []string{"-d", "5", "-subs", "-u", "-wayback"}
	stdinPath := ""
	if isFile {
		stdinPath = target
	} else {
		tmp, _ := os.CreateTemp("", "hak_stdin_*.txt")
		defer os.Remove(tmp.Name())
		os.WriteFile(tmp.Name(), []byte("https://"+target), 0644)
		stdinPath = tmp.Name()
	}
	return ExecCommand(ctx, "hakrawler", args, stdinPath)
}

type WaymoreProvider struct{}

func (p *WaymoreProvider) Name() string { return "waymore" }

func (p *WaymoreProvider) Execute(ctx context.Context, target string, isFile bool) ([]string, error) {
	tmp, _ := os.CreateTemp("", "waymore_*.txt")
	defer os.Remove(tmp.Name())
	
	args := []string{"-i", target, "-mode", "U", "-oU", tmp.Name()}
	_, err := ExecCommand(ctx, "waymore", args, "")
	if err != nil {
		return nil, err
	}
	
	// Read results from temp file
	return utils.ReadLines(tmp.Name())
}

type GospiderProvider struct{}

func (p *GospiderProvider) Name() string { return "gospider" }

func (p *GospiderProvider) Execute(ctx context.Context, target string, isFile bool) ([]string, error) {
	var args []string
	if isFile {
		args = []string{"-S", target, "-c", "20", "-d", "5", "--other-source", "--subs", "--active"}
	} else {
		args = []string{"-s", "https://" + target, "-c", "20", "-d", "5", "--other-source", "--subs", "--active"}
	}
	return ExecCommand(ctx, "gospider", args, "")
}

func RegisterAll(e *engine.Engine) {
	e.AddProvider(&GauProvider{})
	e.AddProvider(&WaybackProvider{})
	e.AddProvider(&KatanaProvider{})
	e.AddProvider(&ParamSpiderProvider{})
	e.AddProvider(&HakrawlerProvider{})
	e.AddProvider(&WaymoreProvider{})
	e.AddProvider(&GospiderProvider{})
}
