package cli

import (
	"fmt"
	"path/filepath"

	"github.com/Raidflux/rojo-syncback-lfs-locking/pkg/git"
	"github.com/Raidflux/rojo-syncback-lfs-locking/pkg/syncback"
	"github.com/ncruces/zenity"
)

func SyncbackLFS(input string) {
	lfsFiles, err := git.LfsListFiles()
	if err != nil {
		fmt.Println(err)
	}

	syncbackFiles, err := syncback.ListSyncbackFiles(input)
	if err != nil {
		fmt.Println(err)
	}

	filesThatNeedLocking := []string{}
	for _, file := range syncbackFiles {
		for _, lfsFile := range lfsFiles {
			cleanFile := filepath.Clean(file)
			clearLfsFile := filepath.Clean(lfsFile)
			if cleanFile == clearLfsFile {
				filesThatNeedLocking = append(filesThatNeedLocking, cleanFile+cleanFile)
			}
		}
	}

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
	)

	if err != nil {
		fmt.Println(err)
	}
}
