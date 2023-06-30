package cmd

import (
	"github.com/spf13/cobra"
)

var configAddCmd = &cobra.Command{
	Use: "add",
}

func init() {
	configCmd.AddCommand(configAddCmd)
}
