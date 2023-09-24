package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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
		configPathRel, err := flags.GetString("configPath")
		if err != nil {
			return fmt.Errorf("unknown error when getting '--configPath'")
		}

		configPath, err := filepath.Abs(configPathRel)
		if err != nil {
			return fmt.Errorf("failed to resolve relative 'configPath' into absolute path")
		}

		// load config
		if err := config.Load(configPathRel == "", configPath); err != nil {
			return err
		}

		if len(args) != 2 {
			return fmt.Errorf("not enough arguments provided")
		}

		credentials := viper.Get("divoltlogincredentials")

		credentialsSlice, ok := credentials.([]any)
		if ok {
			credentialsSlice = append(credentialsSlice, map[string]string{
				"email":    args[0],
				"password": args[1],
			})
		} else {
			return fmt.Errorf("an unknown error has occurred")
		}

		viper.Set("divoltlogincredentials", credentialsSlice)

		return config.Offload()
	},
}

func init() {
	flags := configAddDivoltCredentialCmd.Flags()

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config\n[a custom config file must end in .json]")

	configAddCmd.AddCommand(configAddDivoltCredentialCmd)
}
