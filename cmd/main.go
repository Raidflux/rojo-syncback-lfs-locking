package main

import (
	"fmt"
	"os"

	"github.com/Raidflux/rojo-syncback-lfs-locking/pkg/cli"
	"github.com/Raidflux/rojo-syncback-lfs-locking/pkg/git"
	"github.com/Raidflux/rojo-syncback-lfs-locking/pkg/syncback"
)

func main() {
	if err := git.IsLfsInstalled(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := syncback.IsSyncbackInstalled(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if isRepo := git.PwdIsGitRepoRoot(); !isRepo {
		fmt.Fprintln(os.Stderr, "Current directory is not a git repository")
		os.Exit(1)
	}


	cli.Execute()
}
