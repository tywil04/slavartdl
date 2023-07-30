package cmd

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tywil04/slavartdl/internal/config"
)

var configRemoveTokensCmd = &cobra.Command{
	Use:          "tokens [flags] tokenIndex(s)",
	Short:        "Removes session token using index shown by the list command",
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
			return fmt.Errorf("failed to load config")
		}

		sessionTokens := viper.GetStringSlice("divoltsessiontokens")

		for index := range sessionTokens {
			for _, arg := range args {
				argNumber, err := strconv.Atoi(arg)
				if err == nil && argNumber == index {
					sessionTokens = append(
						sessionTokens[:index],
						sessionTokens[index+1:]...,
					)
				}
			}
		}

		viper.Set("divoltsessiontokens", sessionTokens)

		return nil
	},
}

func init() {
	flags := configRemoveTokensCmd.Flags()

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config")

	configRemoveCmd.AddCommand(configRemoveTokensCmd)
}
