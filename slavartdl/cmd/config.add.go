package cmd

import (
	"github.com/spf13/cobra"
)

var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds to config",
}

func init() {
	configCmd.AddCommand(configAddCmd)
}
