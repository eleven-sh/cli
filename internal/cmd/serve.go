package cmd

import (
	"os"

	"github.com/eleven-sh/agent/config"
	"github.com/eleven-sh/cli/internal/dependencies"
	"github.com/eleven-sh/eleven/features"
	"github.com/spf13/cobra"
)

func loadServeCmd(
	providerCmd *cobra.Command,
	provider *cloudProvider,
) {

	var binding string

	var serveCmd = &cobra.Command{
		Use: "serve <sandbox_name> <port>",

		Short: "Allow TCP traffic on a port",

		Long: `Allow TCP traffic on a port.

This command allows TCP traffic on a port in a sandbox.

Once TCP traffic is allowed, the port becomes reachable from outside.

To reach the port through a domain name (via HTTP(S)), use the "--as" flag.`,

		Example: `  eleven ` + provider.ShortName + ` serve eleven-api 8080
  eleven ` + provider.ShortName + ` serve eleven-api 8080 --as api.eleven.sh`,

		Args: cobra.ExactArgs(2),

		Run: func(cmd *cobra.Command, args []string) {
			envName := args[0]
			port := args[1]

			serveInput := features.ServeInput{
				EnvName:                   envName,
				ReservedPorts:             config.EnvReservedPorts,
				Port:                      port,
				PortBinding:               binding,
				DomainReachabilityChecker: dependencies.ProvideDomainReachabilityChecker(),
			}

			serve := provider.ProvideServeFeature()

			if err := serve.Execute(serveInput); err != nil {
				os.Exit(1)
			}
		},
	}

	serveCmd.Flags().StringVar(
		&binding,
		"as",
		"",
		"the domain name to use to reach the port",
	)

	providerCmd.AddCommand(serveCmd)
}
