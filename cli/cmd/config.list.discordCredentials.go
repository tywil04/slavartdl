package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tywil04/slavartdl/cli/internal/config"
)

var configListDiscordCredentialsCmd = &cobra.Command{
	Use:   "discordCredentials [flags]",
	Short: "Lists stored discord credentials",
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

		for index, credential := range config.Open.DiscordLoginCredentials {
			fmt.Printf("[%d]: Email = %s, Password = %s\n", index, credential.Email, credential.Password)
		}

		return nil
	},
}

func init() {
	flags := configListDiscordCredentialsCmd.Flags()

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config\n[a custom config file must end in .json]")

	configListCmd.AddCommand(configListDiscordCredentialsCmd)
}
