package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/tywil04/slavartdl/cli/internal/config"
)

var configAddDiscordTokensCmd = &cobra.Command{
	Use:          "discordTokens [flags] token(s)",
	Short:        "Adds discord session token(s) to config",
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

		for _, arg := range args {
			arg = strings.TrimSpace(arg)
			if arg != "" {
				config.Open.DiscordSessionTokens = append(config.Open.DiscordSessionTokens, arg)
			}
		}

		return config.SaveConfig()
	},
}

func init() {
	flags := configAddDiscordTokensCmd.Flags()

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config\n[a custom config file must end in .json]")

	configAddCmd.AddCommand(configAddDiscordTokensCmd)
}
