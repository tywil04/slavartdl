package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:           "version [flags]",
	Short:         "Returns the version",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version: v1.1.4")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
