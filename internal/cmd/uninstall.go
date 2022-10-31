package cmd

import (
	"os"

	"github.com/eleven-sh/eleven/features"
	"github.com/spf13/cobra"
)

func loadUninstallCmd(
	providerCmd *cobra.Command,
	provider *cloudProvider,
) {

	var uninstallCmd = &cobra.Command{
		Use: "uninstall",

		Short: "Uninstall Eleven",

		Long: `Uninstall Eleven.

This command uninstall Eleven from your ` + provider.LongName + ` account.

All your sandboxes must be removed before running this command.`,

		Example: "  eleven " + provider.ShortName + " uninstall",

		Run: func(cmd *cobra.Command, args []string) {

			uninstallInput := features.UninstallInput{
				SuccessMessage:            provider.UninstallSuccessMessage,
				AlreadyUninstalledMessage: provider.UninstallAlreadyUninstalledMessage,
			}

			uninstall := provider.ProvideUninstallFeature()

			if err := uninstall.Execute(uninstallInput); err != nil {
				os.Exit(1)
			}
		},
	}

	providerCmd.AddCommand(uninstallCmd)
}
