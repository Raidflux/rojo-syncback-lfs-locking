package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func IsLfsInstalled() error {
	cmd := exec.Command("git", "lfs", "version")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("git-lfs is not installed: %v", err)
	}

	if !strings.Contains(string(output), "git-lfs") {
		return fmt.Errorf("git-lfs is not installed: %s %v", output, err)
	}

	return nil
}

func LfsUnlock(file string) error {
	cmd := exec.Command("git", "lfs", "unlock", file)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("git-lfs unlock failed: %s %v", output, err)
	}

	return nil
}

func LfsLock(file string) error {
	cmd := exec.Command("git", "lfs", "lock", file)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("git-lfs lock failed: %s %v", output, err)
	}

	return nil
}

func LfsListFiles() ([]string, error) {
	lsFilesCmd := exec.Command("git", "ls-files", "--recurse-submodules")
	attrsCmd := exec.Command("git", "check-attr", "--stdin", "lockable")

	var lsFilesOut, lsFilesErr bytes.Buffer
	lsFilesCmd.Stdout = &lsFilesOut
	lsFilesCmd.Stderr = &lsFilesErr

	if err := lsFilesCmd.Start(); err != nil {
		return nil, err
	}

	if err := lsFilesCmd.Wait(); err != nil {
		return nil, errors.New(lsFilesErr.String())
	}

	var attrsOut, attrsErr bytes.Buffer
	attrsCmd.Stdout = &attrsOut
	attrsCmd.Stderr = &attrsErr
	attrsCmd.Stdin = &lsFilesOut

	if err := attrsCmd.Start(); err != nil {
		return nil, err
	}

	if err := attrsCmd.Wait(); err != nil {
		return nil, errors.New(attrsErr.String())
	}

	attrsData := attrsOut.String()
	lines := strings.Split(strings.TrimSpace(attrsData), "\n")
	var lockableFiles []string

	for _, line := range lines {
		parts := strings.Split(line, ": lockable: ")
		if len(parts) == 2 && parts[1] == "set" {
			lockableFiles = append(lockableFiles, parts[0])
		}
	}

	return lockableFiles, nil
}

func LfsListLockedFiles() ([]string, error) {
	cmd := exec.Command("git", "lfs", "locks")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("git-lfs locks failed: %s %v", output, err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var lockedFiles []string

	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			lockedFiles = append(lockedFiles, parts[1])
		}
	}

	return lockedFiles, nil
}

func LfsIsLockMine(file string) (bool, error) {
	cmd := exec.Command("git", "lfs", "locks", "--json", file)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false, fmt.Errorf("git-lfs locks failed: %s %v", output, err)
	}

	if strings.Contains(string(output), "error") {
		return false, nil
	}

	return true, nil
}