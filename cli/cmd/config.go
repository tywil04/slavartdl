package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tywil04/slavartdl/cli/internal/config"
)

var configCmd = &cobra.Command{
	Use:           "config [flags]",
	Short:         "Returns the default config file location",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		location, err := config.DefaultConfigLocation()
		if err != nil {
			return err
		}

		fmt.Printf("The config file is located at: %s\n", location)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
