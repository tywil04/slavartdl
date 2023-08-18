package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/tywil04/slavartdl/internal/update"
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

			// new timeout structure added to this version
			if strings.EqualFold(version, "v1.1.11") {
				fmt.Println("Please note, if you previously have modified 'downloadcmd.timeout.seconds' or 'downloadcmd.timeout.minutes' then your config will no longer work. 'downloadcmd.timeout.minutes' and its related flag '--timeoutMinutes' has been removed. 'downloadcmd.timeout.seconds' has been renamed to 'downloadcmd.timeout' and the related flag '--timeoutSeconds' has been replaced with '--timeout'. Having the support for adding additional seconds and minutes to the download commands timeout was redundant, now timeout is seconds only.")
			}
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
