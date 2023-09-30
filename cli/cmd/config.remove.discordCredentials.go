package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/tywil04/slavartdl/cli/internal/config"
)

var configRemoveDiscordCredentialsCmd = &cobra.Command{
	Use:          "discordCredentials [flags] credentialIndex(s)",
	Short:        "Removes discord credential(s) using index shown by the list command",
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

		for index := range config.Open.DiscordLoginCredentials {
			for _, arg := range args {
				argNumber, err := strconv.Atoi(arg)
				if err == nil && argNumber == index {
					config.Open.DiscordLoginCredentials[index] = nil
				}
			}
		}

		resultingCredentials := []*config.ConfigCredential{}
		for _, token := range config.Open.DiscordLoginCredentials {
			if token != nil {
				resultingCredentials = append(resultingCredentials, token)
			}
		}

		config.Open.DiscordLoginCredentials = resultingCredentials

		return config.SaveConfig()
	},
}

func init() {
	flags := configRemoveDiscordCredentialsCmd.Flags()

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config\n[a custom config file must end in .json]")

	configRemoveCmd.AddCommand(configRemoveDiscordCredentialsCmd)
}
