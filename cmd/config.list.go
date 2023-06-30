package cmd

import (
	"github.com/spf13/cobra"
)

var configListCmd = &cobra.Command{
	Use: "list",
}

func init() {
	configCmd.AddCommand(configListCmd)
}
