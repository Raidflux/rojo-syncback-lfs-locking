package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func PwdIsGitRepoRoot() bool {
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}

	// Run the git command to get the root directory of the Git repo
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	// Get the absolute path of the Git repo root
	repoRoot := strings.TrimSpace(string(out))

	return filepath.Clean(cwd) == filepath.Clean(repoRoot)
}
