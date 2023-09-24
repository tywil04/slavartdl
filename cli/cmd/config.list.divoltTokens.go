package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tywil04/slavartdl/cli/internal/config"
)

var configListDivoltTokensCmd = &cobra.Command{
	Use:   "divoltTokens [flags]",
	Short: "Lists stored divolt session tokens",
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

		sessionTokens := viper.GetStringSlice("divoltsessiontokens")

		for index, token := range sessionTokens {
			fmt.Printf("[%d]: %s\n", index, token)
		}

		return nil
	},
}

func init() {
	flags := configListDivoltTokensCmd.Flags()

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config\n[a custom config file must end in .json]")

	configListCmd.AddCommand(configListDivoltTokensCmd)
}
