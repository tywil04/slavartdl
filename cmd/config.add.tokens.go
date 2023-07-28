package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tywil04/slavart/lib/config"
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
