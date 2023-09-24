package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		sessionTokens := viper.GetStringSlice("discordsessiontokens")

		for _, arg := range args {
			arg = strings.TrimSpace(arg)
			if arg != "" && arg != "<DELETED>" {
				sessionTokens = append(sessionTokens, arg)
			}
		}

		viper.Set("discordsessiontokens", sessionTokens)

		return config.Offload()
	},
}

func init() {
	flags := configAddDiscordTokensCmd.Flags()

	flags.StringP("configPath", "C", "", "a directory that contains an override config.json file\nor a file which contains an override config\n[a custom config file must end in .json]")

	configAddCmd.AddCommand(configAddDiscordTokensCmd)
}
