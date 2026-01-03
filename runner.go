package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func RunCommand(ctx context.Context, command string, stdin *os.File) ([]string, error) {
	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	if stdin != nil {
		cmd.Stdin = stdin
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var results []string
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		results = append(results, scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		// Some tools might return non-zero exit codes even on success (e.g. grep)
		// Or if context was cancelled.
		return results, err
	}

	return results, nil
}

func RunParallel(ctx context.Context, commands []string, stdinPath string) []string {
	var wg sync.WaitGroup
	resultsChan := make(chan []string, len(commands))

	for _, cmdStr := range commands {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			var stdin *os.File
			if stdinPath != "" {
				var err error
				stdin, err = os.Open(stdinPath)
				if err != nil {
					LogError(fmt.Sprintf("Failed to open stdin file %s: %v", stdinPath, err))
					return
				}
				defer stdin.Close()
			}
			
			res, _ := RunCommand(ctx, c, stdin)
			resultsChan <- res
		}(cmdStr)
	}

	wg.Wait()
	close(resultsChan)

	var allResults []string
	for res := range resultsChan {
		allResults = append(allResults, res...)
	}

	return allResults
}
