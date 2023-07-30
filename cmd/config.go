package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tywil04/slavartdl/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config [flags]",
	Short: "Returns the default config file location",
	RunE: func(cmd *cobra.Command, args []string) error {
		config.CreateConfigIfNotExist()

		userConfigDirectory, err := config.GetConfigPath()
		if err != nil {
			return err
		}

		fmt.Printf("The config file is located at: %s\n", userConfigDirectory)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
