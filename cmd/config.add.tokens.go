package cmd

import (
	"slavartdl/lib/config"

	"github.com/spf13/cobra"
)

var configAddTokensCmd = &cobra.Command{
	Use:   "tokens ...tokens",
	Short: "adds token to config",
	Long:  "adds token to config",
	RunE: func(cmd *cobra.Command, args []string) error {
		config.CreateConfigIfNotExist()

		config.Public.DivoltSessionTokens = append(config.Public.DivoltSessionTokens, args...)
		return config.WriteConfig()
	},
}

func init() {
	configAddCmd.AddCommand(configAddTokensCmd)
}
