package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		cli.Execute(ctx)
	}()

	<-ctx.Done()
	wg.Wait()
}
