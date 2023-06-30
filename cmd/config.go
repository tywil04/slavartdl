package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:       "config",
	Short:     "get the config file location",
	Long:      "get the config file location",
	Args:      cobra.ExactArgs(0),
	ValidArgs: []string{"url"},
	RunE: func(cmd *cobra.Command, args []string) error {
		userConfigDirectory, err := os.UserConfigDir()
		if err != nil {
			return err
		}

		fmt.Printf("The config file is located at:\n%s/slavart/config.yaml\n", userConfigDirectory)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
