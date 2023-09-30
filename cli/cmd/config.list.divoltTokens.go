package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tywil04/slavartdl/cli/internal/config"
)

var configListDivoltTokensCmd = &cobra.Command{
	Use:   "divoltTokens [flags]",
	Short: "Lists stored divolt session tokens",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()

		// optional
		configPath, err := flags.GetString("configPath")
		if err != nil {
			return err
		}

		// load config
		if err := config.OpenConfig(configPath); err != nil {
			return err
		}

		for index, token := range config.Open.DivoltSessionTokens {
			fmt.Printf("[%d]: %s\n", index, token)
		}

		return nil
	},
}

func init() {
	flags := configListDivoltTokensCmd.Flags()

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config\n[a custom config file must end in .json]")

	configListCmd.AddCommand(configListDivoltTokensCmd)
}
