package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configAddTokensCmd = &cobra.Command{
	Use:          "tokens [flags] token(s)",
	Short:        "Adds session token to config",
	SilenceUsage: true,
	Args:         cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sessionTokens := viper.GetStringSlice("divoltsessiontokens")
		sessionTokens = append(sessionTokens, args...)
		viper.Set("divoltsessiontokens", sessionTokens)
	},
}

func init() {
	configAddCmd.AddCommand(configAddTokensCmd)
}
