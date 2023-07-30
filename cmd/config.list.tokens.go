package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configListTokensCmd = &cobra.Command{
	Use:   "tokens [flags]",
	Short: "Lists stored session tokens",
	Run: func(cmd *cobra.Command, args []string) {
		sessionTokens := viper.GetStringSlice("divoltsessiontokens")

		for index, token := range sessionTokens {
			fmt.Printf("[%d]: %s\n", index, token)
		}
	},
}

func init() {
	configListCmd.AddCommand(configListTokensCmd)
}
