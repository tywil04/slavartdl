package cmd

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configRemoveTokensCmd = &cobra.Command{
	Use:          "tokens [flags] tokenIndex(s)",
	Short:        "Removes session token using index shown by the list command",
	SilenceUsage: true,
	Args:         cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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
	},
}

func init() {
	configRemoveCmd.AddCommand(configRemoveTokensCmd)
}
