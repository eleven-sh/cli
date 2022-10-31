package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/eleven-sh/cli/internal/dependencies"
	"github.com/eleven-sh/cli/internal/exceptions"
	"github.com/eleven-sh/cli/internal/vscode"
	"github.com/eleven-sh/eleven/features"
	"github.com/spf13/cobra"
)

func loadEditCmd(
	providerCmd *cobra.Command,
	provider *cloudProvider,
) {

	var editCmd = &cobra.Command{
		Use: "edit <sandbox_name>",

		Short: "Connect your editor to a sandbox",

		Long: `Connect your editor to a sandbox.

This command connects your preferred editor to a sandbox.

Supported editor(s): Microsoft Visual Studio Code`,

		Example: "  eleven " + provider.ShortName + " edit eleven-api",

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {

			elevenViewableErrorBuilder := dependencies.ProvideElevenViewableErrorBuilder()
			baseView := dependencies.ProvideBaseView()

			missingRequirements := []string{}

			vscodeCLI := vscode.CLI{}
			_, err := vscodeCLI.LookupPath(runtime.GOOS)

			if vscodeCLINotFoundErr, ok := err.(vscode.ErrCLINotFound); ok {
				missingRequirements = append(
					missingRequirements,
					fmt.Sprintf(
						"Visual Studio Code (looked in \"%s)",
						strings.Join(vscodeCLINotFoundErr.VisitedPaths, "\", \"")+"\"",
					),
				)
			}

			if len(missingRequirements) > 0 {
				missingRequirementsErr := exceptions.ErrMissingRequirements{
					MissingRequirements: missingRequirements,
				}

				baseView.ShowErrorViewWithStartingNewLine(
					elevenViewableErrorBuilder.Build(
						missingRequirementsErr,
					),
				)

				os.Exit(1)
			}

			envName := args[0]

			editInput := features.EditInput{
				EnvName: envName,
			}

			edit := provider.ProvideEditFeature()

			if err := edit.Execute(editInput); err != nil {
				os.Exit(1)
			}
		},
	}

	providerCmd.AddCommand(editCmd)
}
