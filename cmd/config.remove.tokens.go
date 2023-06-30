package cmd

import (
	"slavartdl/lib/config"
	"strconv"

	"github.com/spf13/cobra"
)

var configRemoveTokensCmd = &cobra.Command{
	Use:   "tokens ...tokenIndexes",
	Short: "removes token using index shown when list command is used",
	Long:  "removes token using index shown when list command is used",
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
