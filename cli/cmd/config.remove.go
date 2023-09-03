package cmd

import (
	"github.com/spf13/cobra"
)

var configRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove from config",
}

func init() {
	configCmd.AddCommand(configRemoveCmd)
}
