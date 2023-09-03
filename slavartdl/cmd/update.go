package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tywil04/slavartdl/slavartdl/internal/update"
)

var updateCmd = &cobra.Command{
	Use:          "update [flags]",
	Short:        "Updates this tool",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()

		// optional
		force, err := flags.GetBool("force")
		if err != nil {
			return fmt.Errorf("unknown error when getting '--force'")
		}

		version, err := update.Update(force)
		if err != nil {
			return err
		}

		if version != "" {
			fmt.Printf("Successfully updated to version %s!\n", version)
		} else {
			fmt.Println("All up to date, no updated required!")
		}

		return nil
	},
}

func init() {
	flags := updateCmd.Flags()

	flags.BoolP("force", "f", false, "forces update even if slavartdl is already on the latest version")

	rootCmd.AddCommand(updateCmd)
}
