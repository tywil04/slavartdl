package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "slavartdl",
	Short: "SlavartDL:\nUtilitiy to download from SlavArt Divolt server",
}

func init() {
	rootCmd.DisableAutoGenTag = true
}

func Execute() error {
	// hide "completion" command
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	return rootCmd.Execute()
}
