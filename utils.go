package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CheckAndInstall(tool string, installCmd string) {
	_, err := exec.LookPath(tool)
	if err != nil {
		LogError(fmt.Sprintf("%s is not installed.", tool))
		fmt.Printf("Do you want to install %s? (y/n): ", tool)
		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))

		if choice == "y" {
			LogStatus(fmt.Sprintf("Installing %s...", tool))
			cmd := exec.Command("sh", "-c", installCmd)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err != nil {
				LogError(fmt.Sprintf("Failed to install %s: %v", tool, err))
				os.Exit(1)
			}
			LogStatus(fmt.Sprintf("%s installed successfully.", tool))
		} else {
			LogError(fmt.Sprintf("%s is required. Exiting.", tool))
			os.Exit(1)
		}
	}
}


func WriteLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}

func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
