package git

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

type lfsFileLock struct {
	ID    string `json:"id"`
	Path  string `json:"path"`
	Owner struct {
		Name string `json:"name"`
	}

	LockedAt time.Time `json:"locked_at"`
}

type lfsLocksContainer struct {
	Ours   []lfsFileLock `json:"ours"`
	Theirs []lfsFileLock `json:"theirs"`
}

func LfsListLockedFiles() (*lfsLocksContainer, error) {
	cmd := exec.Command("git", "lfs", "locks", "--verify", "--json")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, fmt.Errorf("git-lfs locks failed: %s %v", output, err)
	}

	var locks lfsLocksContainer
	if err := json.Unmarshal(output, &locks); err != nil {
		return nil, err
	}

	for i := range locks.Ours {
		locks.Ours[i].Path = filepath.Clean(locks.Ours[i].Path)
	}

	for i := range locks.Theirs {
		locks.Theirs[i].Path = filepath.Clean(locks.Theirs[i].Path)
	}

	return &locks, nil
}
