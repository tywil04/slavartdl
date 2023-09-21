package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tywil04/slavartdl/cli/internal/config"
)

var configListCredentialsCmd = &cobra.Command{
	Use:   "credentials [flags]",
	Short: "Lists stored credentials",
	Args:  cobra.ExactArgs(0),
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
			for index, slice := range credentialsSlice {
				sliceMap, ok := slice.(map[string]any)
				if ok {
					fmt.Printf("[%d]: Email = %s, Password = %s\n", index, sliceMap["email"], sliceMap["password"])
				} else {
					return fmt.Errorf("an unknown error has occurred")
				}
			}
		} else {
			return fmt.Errorf("an unknown error has occurred")
		}

		return nil
	},
}

func init() {
	flags := configListCredentialsCmd.Flags()

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config\n[a custom config file must end in .json]")

	configListCmd.AddCommand(configListCredentialsCmd)
}
