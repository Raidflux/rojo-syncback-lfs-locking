package syncback

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

func IsSyncbackInstalled() error {
	cmd := exec.Command("rojo", "syncback", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("rojo-syncback is not installed: %s %v", output, err)
	}

	return nil
}

func RunSyncback(input string, dryRun bool) (string, error) {
	inputClean := filepath.Clean(input)
	inputFolder := filepath.Dir(inputClean)

	cmd := exec.Command("rojo", "syncback", "-y", "--input", inputClean, inputFolder)
	if dryRun {
		cmd.Args = append(cmd.Args, "--dry-run", "--list")
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", fmt.Errorf("rojo-syncback failed: %s %v", output, err)
	}

	return stripColorCodes(string(output)), nil
}

func ListSyncbackFiles(input string) ([]string, error) {
	output, err := RunSyncback(input, true)
	if err != nil {
		return nil, err
	}

	filesToWrite, filesToRemove := parseAndGetFilesToEditOrRemove(output)

	files := append(filesToWrite, filesToRemove...)
	sort.Strings(files)

	return files, nil
}

func parseAndGetFilesToEditOrRemove(output string) ([]string, []string) {
	var filesToWrite []string
	var filesToRemove []string

	scanner := bufio.NewScanner(strings.NewReader(output))
	isWritingSection := false
	isRemovingSection := false

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Writing files/directories:") {
			isWritingSection = true
			continue
		} else if strings.HasPrefix(line, "Removing files/directories:") {
			isRemovingSection = true
			continue
		} else if strings.HasPrefix(line, "Would write") {
			isWritingSection = false
			isRemovingSection = false
			continue
		}

		if isRemovingSection {
			filesToRemove = append(filesToRemove, strings.TrimSpace(line))
		} else if isWritingSection {
			filesToWrite = append(filesToWrite, strings.TrimSpace(line))
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}

	return filesToWrite, filesToRemove
}

func stripColorCodes(input string) string {
	re := regexp.MustCompile(`\x1b\[[;?0-9]*m`)
	return re.ReplaceAllString(input, "")
}
