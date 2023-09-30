package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tywil04/slavartdl/cli/internal/config"
)

var configAddDivoltCredentialCmd = &cobra.Command{
	Use:          "divoltCredential [flags] email password",
	Short:        "Adds divolt credential token to config",
	SilenceUsage: true,
	Args:         cobra.ExactArgs(2),
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

		if len(args) != 2 {
			return fmt.Errorf("not enough arguments provided")
		}

		config.Open.DivoltLoginCredentials = append(config.Open.DivoltLoginCredentials, &config.ConfigCredential{
			Email:    args[0],
			Password: args[1],
		})

		return config.SaveConfig()
	},
}

func init() {
	flags := configAddDivoltCredentialCmd.Flags()

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config\n[a custom config file must end in .json]")

	configAddCmd.AddCommand(configAddDivoltCredentialCmd)
}
