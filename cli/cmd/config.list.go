package cmd

import (
	"github.com/spf13/cobra"
)

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists from config",
}

func init() {
	configCmd.AddCommand(configListCmd)
}
