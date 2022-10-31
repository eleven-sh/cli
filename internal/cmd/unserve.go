package cmd

import (
	"os"

	"github.com/eleven-sh/agent/config"
	"github.com/eleven-sh/eleven/features"
	"github.com/spf13/cobra"
)

func loadUnserveCmd(
	providerCmd *cobra.Command,
	provider *cloudProvider,
) {

	var unserveCmd = &cobra.Command{
		Use: "unserve <sandbox_name> <port>",

		Short: "Disallow TCP traffic on a port",

		Long: `Disallow TCP traffic on a port.

This command disallows TCP traffic on a port in a sandbox.

Once TCP traffic is disallowed, the port becomes unreachable from outside.`,

		Example: "  eleven " + provider.ShortName + " unserve eleven-api 8080",

		Args: cobra.ExactArgs(2),

		Run: func(cmd *cobra.Command, args []string) {
			envName := args[0]
			port := args[1]

			unserveInput := features.UnserveInput{
				EnvName:       envName,
				ReservedPorts: config.EnvReservedPorts,
				Port:          port,
			}

			unserve := provider.ProvideUnserveFeature()

			if err := unserve.Execute(unserveInput); err != nil {
				os.Exit(1)
			}
		},
	}

	providerCmd.AddCommand(unserveCmd)
}
