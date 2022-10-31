package cmd

import (
	"github.com/eleven-sh/cli/internal/globals"
	"github.com/eleven-sh/eleven/features"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type cloudProviderFeature[T interface{}] interface {
	Execute(input T) error
}

type cloudProvider struct {
	LongName   string
	ShortName  string
	GlobalName globals.CloudProvider

	DefaultInstanceType string
	ExampleInstanceType string

	AddFlagsToBaseCmd func(*cobra.Command)

	ProvideInitFeature    func() cloudProviderFeature[features.InitInput]
	ProvideEditFeature    func() cloudProviderFeature[features.EditInput]
	ProvideRemoveFeature  func() cloudProviderFeature[features.RemoveInput]
	ProvideServeFeature   func() cloudProviderFeature[features.ServeInput]
	ProvideUnserveFeature func() cloudProviderFeature[features.UnserveInput]

	UninstallSuccessMessage            string
	UninstallAlreadyUninstalledMessage string
	ProvideUninstallFeature            func() cloudProviderFeature[features.UninstallInput]
}

func addCloudProviderCmdToRootCmd(
	rootCmd *cobra.Command,
	provider *cloudProvider,
) *cobra.Command {

	var cloudProviderCmd = &cobra.Command{
		Use: provider.ShortName,

		Short: "Use Eleven on " + provider.LongName,

		Long: `Use Eleven on ` + provider.LongName + `.
		
To begin, create your first sandbox using the command:
	
  eleven ` + provider.ShortName + ` init <sandbox_name>
	
Once created, you may want to connect your editor to it using the command: 
	
  eleven ` + provider.ShortName + ` edit <sandbox_name>
	
If you don't plan to use this sandbox again, you could remove it using the command:
		
  eleven ` + provider.ShortName + ` remove <sandbox_name>`,

		Example: `  eleven ` + provider.ShortName + ` init eleven-api --instance-type ` + provider.ExampleInstanceType + ` 
  eleven ` + provider.ShortName + ` edit eleven-api
  eleven ` + provider.ShortName + ` remove eleven-api`,

		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ensureUserIsLoggedIn()
			populateCurrentCloudProviderGlobals(provider.GlobalName, cmd)
		},
	}

	provider.AddFlagsToBaseCmd(cloudProviderCmd)

	rootCmd.AddCommand(cloudProviderCmd)

	return cloudProviderCmd
}

func populateCurrentCloudProviderGlobals(
	cloudProvider globals.CloudProvider,
	cloudProviderCmd *cobra.Command,
) {

	globals.CurrentCloudProvider = cloudProvider
	globals.CurrentCloudProviderArgs = ""

	// command ran without subcommand -> help displayed -> no args to parse
	if !cloudProviderCmd.HasParent() {
		return
	}

	cloudProviderCmd.Parent().Flags().VisitAll(func(f *pflag.Flag) {
		if !f.Changed {
			return
		}

		if len(globals.CurrentCloudProviderArgs) > 0 {
			globals.CurrentCloudProviderArgs += " "
		}

		globals.CurrentCloudProviderArgs += "--" + f.Name + " " + f.Value.String()
	})
}

func init() {
	availableCloudProviders := []*cloudProvider{
		getAWSCloudProvider(),
		getHetznerCloudProvider(),
	}

	for _, cloudProvider := range availableCloudProviders {
		cloudProviderCmd := addCloudProviderCmdToRootCmd(rootCmd, cloudProvider)

		loadInitCmd(cloudProviderCmd, cloudProvider)
		loadEditCmd(cloudProviderCmd, cloudProvider)

		loadServeCmd(cloudProviderCmd, cloudProvider)
		loadUnserveCmd(cloudProviderCmd, cloudProvider)

		loadRemoveCmd(cloudProviderCmd, cloudProvider)
		loadUninstallCmd(cloudProviderCmd, cloudProvider)
	}
}
