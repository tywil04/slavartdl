package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tywil04/slavartdl/cli/internal/config"
)

var configRemoveDiscordTokensCmd = &cobra.Command{
	Use:          "discordTokens [flags] tokenIndex(s)",
	Short:        "Removes discord session token(s) using index shown by the list command",
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

		for index := range config.Open.DiscordSessionTokens {
			for _, arg := range args {
				argNumber, err := strconv.Atoi(arg)
				if err == nil && argNumber == index {
					config.Open.DiscordSessionTokens[index] = "<DELETED>"
				}
			}
		}

		resultingSessionTokens := []string{}
		for _, token := range config.Open.DiscordSessionTokens {
			if token != "<DELETED>" {
				resultingSessionTokens = append(resultingSessionTokens, token)
			}
		}

		config.Open.DiscordSessionTokens = resultingSessionTokens

		return config.SaveConfig()
	},
}

func init() {
	flags := configRemoveDiscordTokensCmd.Flags()

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config\n[a custom config file must end in .json]")

	configRemoveCmd.AddCommand(configRemoveDiscordTokensCmd)
}
