package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var syncbackLfsInput string

var rootCmd = &cobra.Command{
	Use:   "syncback-lfs",
	Short: "Wrapper for Rojo syncback to add file watching and LFS file locking support into fully managed rojo workflows with the goal to create a more professional development workflow",
	Run: func(cmd *cobra.Command, args []string) {
		SyncbackLFS(cmd.Context(), syncbackLfsInput)
	},
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&syncbackLfsInput, "input", "i", "", "The input file to pass to Rojo syncback")
	rootCmd.MarkFlagRequired("input")
}
