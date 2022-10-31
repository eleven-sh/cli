package cmd

import (
	"os"

	"github.com/eleven-sh/cli/internal/dependencies"
	"github.com/eleven-sh/cli/internal/system"
	"github.com/eleven-sh/eleven/features"
	"github.com/spf13/cobra"
)

func loadRemoveCmd(
	providerCmd *cobra.Command,
	provider *cloudProvider,
) {

	var forceRemove bool

	var removeCmd = &cobra.Command{
		Use: "remove <sandbox_name>",

		Short: "Remove a sandbox",

		Long: `Remove an existing sandbox.

The sandbox will be PERMANENTLY removed along with ALL your un-pushed work.
		
There is no going back, so please be sure before running this command.`,

		Example: `  eleven ` + provider.ShortName + ` remove eleven-api
  eleven ` + provider.ShortName + ` remove eleven-api --force`,

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			envName := args[0]

			removeInput := features.RemoveInput{
				EnvName:       envName,
				PreRemoveHook: dependencies.ProvidePreRemoveHook(),
				ForceRemove:   forceRemove,
				ConfirmRemove: func() (bool, error) {
					logger := system.NewLogger()
					return system.AskForConfirmation(
						logger,
						os.Stdin,
						"All your un-pushed work will be lost.",
					)
				},
			}

			remove := provider.ProvideRemoveFeature()

			if err := remove.Execute(removeInput); err != nil {
				os.Exit(1)
			}
		},
	}

	removeCmd.Flags().BoolVar(
		&forceRemove,
		"force",
		false,
		"remove without confirmation",
	)

	providerCmd.AddCommand(removeCmd)
}
