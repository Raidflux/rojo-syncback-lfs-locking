package cli

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/Raidflux/rojo-syncback-lfs-locking/pkg/git"
	"github.com/Raidflux/rojo-syncback-lfs-locking/pkg/syncback"
	"github.com/ncruces/zenity"
)

func SyncbackLFS(ctx context.Context, input string) {
	lfsFiles, err := git.LfsListFiles()
	if err != nil {
		fmt.Println(err)
		return
	}

	syncbackFiles, err := syncback.ListSyncbackFiles(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	filesThatNeedLocking := []string{}
	for _, file := range syncbackFiles {
		for _, lfsFile := range lfsFiles {
			cleanFile := filepath.Clean(file)
			clearLfsFile := filepath.Clean(lfsFile)
			if cleanFile == clearLfsFile {
				filesThatNeedLocking = append(filesThatNeedLocking, cleanFile)
				break
			}
		}
	}

	lockedFiles, err := git.LfsListLockedFiles()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, theirFile := range lockedFiles.Theirs {
		for _, fileThatNeedsLocking := range filesThatNeedLocking {
			if theirFile.Path == fileThatNeedsLocking {
				zenity.Error(
					fmt.Sprintf("File %s is already locked by %s", theirFile.Path, theirFile.Owner.Name),
					zenity.Title("File already locked"),
					zenity.Context(ctx),
				)
				return
			}
		}
	}

	for _, ourFile := range lockedFiles.Ours {
		toRemove := -1
		for i, fileThatNeedsLocking := range filesThatNeedLocking {
			if ourFile.Path == fileThatNeedsLocking {
				toRemove = i
				break
			}
		}

		if toRemove != -1 {
			filesThatNeedLocking = append(filesThatNeedLocking[:toRemove], filesThatNeedLocking[toRemove+1:]...)
		}
	}

	if len(filesThatNeedLocking) > 0 {
		infoMessage := "Files that need locking:\n\n"

		for _, file := range filesThatNeedLocking {
			fmt.Println("Needs locking: ", file)
			infoMessage += file + "\n\n"
		}

		err = zenity.Question(
			infoMessage,
			zenity.Title("Files need to be locked"),
			zenity.OKLabel("Lock files"),
			zenity.Icon(zenity.WarningIcon),
			zenity.Context(ctx),
		)

		if err == nil {
			for _, file := range filesThatNeedLocking {
				if err := git.LfsLock(file); err != nil {
					fmt.Println(err)
					zenity.Error(
						err.Error(),
						zenity.Title("Failed to lock file"),
						zenity.Context(ctx),
					)
				}
			}

			zenity.Info(
				fmt.Sprintf("Successfully locked %d files", len(filesThatNeedLocking)),
				zenity.Title("Files locked"),
				zenity.Context(ctx),
			)
		}
	}

	output, err := syncback.RunSyncback(input, false)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(output)
}
