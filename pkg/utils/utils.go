package utils

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

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

func Unique(lines []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range lines {
		if _, value := keys[entry]; !value && entry != "" {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func SanitizeFilename(name string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9.-]`)
	return re.ReplaceAllString(name, "_")
}
