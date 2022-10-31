package cmd

import (
	"os"

	"github.com/eleven-sh/cli/internal/dependencies"
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/features"
	"github.com/spf13/cobra"
)

func loadInitCmd(
	providerCmd *cobra.Command,
	provider *cloudProvider,
) {

	var instanceType string
	var runtimes []string
	var repositories []string

	var initCmd = &cobra.Command{
		Use: "init <sandbox_name>",

		Short: "Initialize a sandbox",

		Long: `Initialize a new sandbox.

To choose the type of instance that will run the sandbox, use the "--instance-type" flag.

To install some runtimes in the sandbox, use the "--runtimes" flag.

To clone some GitHub repositories in the sandbox, use the "--repositories" flag.`,

		Example: `  eleven ` + provider.ShortName + ` init eleven-api
  eleven ` + provider.ShortName + ` init eleven-api --instance-type ` + provider.ExampleInstanceType + ` --runtimes node@18.7.0,docker --repositories repo,organization/repo`,

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {

			elevenViewableErrorBuilder := dependencies.ProvideElevenViewableErrorBuilder()
			baseView := dependencies.ProvideBaseView()

			checkForRepositoryExistence := true
			repositoryResolver := dependencies.ProvideEnvRepositoriesResolver()
			repositories, err := repositoryResolver.Resolve(
				repositories,
				checkForRepositoryExistence,
			)

			if err != nil {
				baseView.ShowErrorViewWithStartingNewLine(
					elevenViewableErrorBuilder.Build(
						err,
					),
				)

				os.Exit(1)
			}

			envName := args[0]

			sshConfig := dependencies.ProvideSSHConfigManager()
			sshCfgDupHostCt, err := sshConfig.CountEntriesWithHostPrefix(
				entities.BuildInitialLocalSSHCfgHostnameForEnv(envName),
			)

			if err != nil {
				baseView.ShowErrorViewWithStartingNewLine(
					elevenViewableErrorBuilder.Build(
						err,
					),
				)

				os.Exit(1)
			}

			initInput := features.InitInput{
				InstanceType:         instanceType,
				EnvName:              envName,
				LocalSSHCfgDupHostCt: sshCfgDupHostCt,
				Repositories:         repositories,
				Runtimes:             runtimes,
			}

			init := provider.ProvideInitFeature()

			if err := init.Execute(initInput); err != nil {
				os.Exit(1)
			}
		},
	}

	initCmd.Flags().StringVar(
		&instanceType,
		"instance-type",
		provider.DefaultInstanceType,
		"the type of instance that will run the sandbox",
	)

	initCmd.Flags().StringSliceVar(
		&runtimes,
		"runtimes",
		[]string{},
		"the runtimes to install in the sandbox",
	)

	initCmd.Flags().StringSliceVar(
		&repositories,
		"repositories",
		[]string{},
		"the repositories to clone in the sandbox",
	)

	providerCmd.AddCommand(initCmd)
}
