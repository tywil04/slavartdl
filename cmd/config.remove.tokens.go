package cmd

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/tywil04/slavartdl/internal/config"
)

var configRemoveTokensCmd = &cobra.Command{
	Use:   "tokens [flags] tokenIndex(s)",
	Short: "Removes session token using index shown by the list command",
	RunE: func(cmd *cobra.Command, args []string) error {
		config.CreateConfigIfNotExist()

		for index := range config.Public.DivoltSessionTokens {
			for _, arg := range args {
				argNumber, err := strconv.Atoi(arg)
				if err == nil && argNumber == index {
					config.Public.DivoltSessionTokens = append(
						config.Public.DivoltSessionTokens[:index],
						config.Public.DivoltSessionTokens[index+1:]...,
					)
				}
			}
		}

		return config.WriteConfig()
	},
}

func init() {
	configRemoveCmd.AddCommand(configRemoveTokensCmd)
}
