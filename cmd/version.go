package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tywil04/slavartdl/internal/update"
)

var versionCmd = &cobra.Command{
	Use:           "version [flags]",
	Short:         "Returns the version",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", update.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
