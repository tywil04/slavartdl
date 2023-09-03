package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tywil04/slavartdl/cli/internal/config"
)

var configRemoveCredentialsCmd = &cobra.Command{
	Use:          "credentials [flags] credentialIndex(s)",
	Short:        "Removes credential(s) using index shown by the list command",
	SilenceUsage: true,
	Args:         cobra.MinimumNArgs(1),
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

		credentials := viper.Get("divoltlogincredentials")

		credentialsSlice, ok := credentials.([]any)
		if ok {
			for index := range credentialsSlice {
				for _, arg := range args {
					argNumber, err := strconv.Atoi(arg)
					if err == nil && argNumber == index {
						credentialsSlice[index] = "<DELETED>"
					}
				}
			}
		} else {
			return fmt.Errorf("an unknown error has occured")
		}

		resultingCredentials := []any{}
		for _, token := range credentialsSlice {
			if token != "<DELETED>" {
				resultingCredentials = append(resultingCredentials, token)
			}
		}

		viper.Set("divoltlogincredentials", resultingCredentials)

		return config.Offload()
	},
}

func init() {
	flags := configRemoveCredentialsCmd.Flags()

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config\n[a custom config file must end in .json]")

	configRemoveCmd.AddCommand(configRemoveCredentialsCmd)
}
