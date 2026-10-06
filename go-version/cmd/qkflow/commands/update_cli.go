package commands

import (
	"github.com/spf13/cobra"
	"github.com/yiminlab/quick-workflow/internal/updater"
)

var updateCLICmd = &cobra.Command{
	Use:   "update-cli",
	Short: "Update qkflow to the latest version",
	Long:  `Check for updates and install the latest version of qkflow.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := updater.ManualUpdate(Version); err != nil {
			cmd.PrintErrf("Failed to update: %v\n", err)
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCLICmd)
}

