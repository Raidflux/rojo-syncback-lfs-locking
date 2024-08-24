package cli

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var watchInput string

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "syncback-lfs once and then every time the input file changes",
	Run: func(cmd *cobra.Command, args []string) {
		Watch(cmd.Context(), watchInput)
	},
}

func init() {
	watchCmd.Flags().StringVarP(&watchInput, "input", "i", "", "The input file to pass to Rojo syncback")
	watchCmd.MarkFlagRequired("input")

	rootCmd.AddCommand(watchCmd)
}

func Watch(ctx context.Context, input string) {
	fmt.Println("Press Ctrl+C to exit")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Has(fsnotify.Create) {
					if event.Name == filepath.Base(input) {
						SyncbackLFS(ctx, input)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(filepath.Dir(input))
	if err != nil {
		log.Fatal(err)
	}

	<-ctx.Done()
}
