package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tywil04/slavartdl/cli/internal/config"
)

var configRemoveDivoltTokensCmd = &cobra.Command{
	Use:          "divoltTokens [flags] tokenIndex(s)",
	Short:        "Removes divolt session token(s) using index shown by the list command",
	SilenceUsage: true,
	Args:         cobra.MinimumNArgs(1),
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

		for index := range config.Open.DivoltSessionTokens {
			for _, arg := range args {
				argNumber, err := strconv.Atoi(arg)
				if err == nil && argNumber == index {
					config.Open.DivoltSessionTokens[index] = "<DELETED>"
				}
			}
		}

		resultingSessionTokens := []string{}
		for _, token := range config.Open.DivoltSessionTokens {
			if token != "<DELETED>" {
				resultingSessionTokens = append(resultingSessionTokens, token)
			}
		}

		config.Open.DivoltSessionTokens = resultingSessionTokens

		return config.SaveConfig()
	},
}

func init() {
	flags := configRemoveDivoltTokensCmd.Flags()

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config\n[a custom config file must end in .json]")

	configRemoveCmd.AddCommand(configRemoveDivoltTokensCmd)
}
