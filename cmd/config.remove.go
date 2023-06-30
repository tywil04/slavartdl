package cmd

import (
	"github.com/spf13/cobra"
)

var configRemoveCmd = &cobra.Command{
	Use: "remove",
}

func init() {
	configCmd.AddCommand(configRemoveCmd)
}
