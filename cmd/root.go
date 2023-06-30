package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "slavartdl",
	Short: "slavartdl",
	Long:  "slavartdl",
}

func init() {
	rootCmd.DisableAutoGenTag = true
}

func Execute() error {
	return rootCmd.Execute()
}
