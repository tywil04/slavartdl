package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tywil04/slavartdl/internal/update"
)

var updateCmd = &cobra.Command{
	Use:           "update [flags]",
	Short:         "Updates this tool",
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		newUpdate, err := update.Update()
		if err != nil {
			return err
		}

		if newUpdate {
			fmt.Printf("Successfully updated to version %s!\n", update.Version)
		} else {
			fmt.Println("All up to date, no updated required!")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
