package providers

import (
	"bufio"
	"context"
	"os"
	"os/exec"
)

// ExecCommand is a helper to run external binaries safely.
func ExecCommand(ctx context.Context, name string, args []string, stdinPath string) ([]string, error) {
	cmd := exec.CommandContext(ctx, name, args...)
	
	if stdinPath != "" {
		f, err := os.Open(stdinPath)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		cmd.Stdin = f
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
		return results, err
	}

	return results, nil
}
