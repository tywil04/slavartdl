package cmd

import (
	"fmt"
	"slavartdl/lib/config"

	"github.com/spf13/cobra"
)

var configListTokensCmd = &cobra.Command{
	Use:   "tokens",
	Short: "lists stored tokens",
	Long:  "lists stored tokens",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		config.CreateConfigIfNotExist()

		for index, token := range config.Public.DivoltSessionTokens {
			fmt.Printf("[%d]: %s\n", index, token)
		}
	},
}

func init() {
	configListCmd.AddCommand(configListTokensCmd)
}