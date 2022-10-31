package cmd

import (
	"github.com/eleven-sh/cli/internal/dependencies"
	"github.com/eleven-sh/cli/internal/globals"
	"github.com/eleven-sh/cli/internal/system"
	"github.com/eleven-sh/eleven/features"
	"github.com/spf13/cobra"
)

func getHetznerCloudProvider() *cloudProvider {
	var context string
	var region string

	var hetzner = cloudProvider{
		ShortName:  "hetzner",
		LongName:   "Hetzner",
		GlobalName: globals.HetznerCloudProvider,

		DefaultInstanceType: "cx11",
		ExampleInstanceType: "cx21",

		AddFlagsToBaseCmd: func(baseCmd *cobra.Command) {
			baseCmd.Flags().StringVar(
				&context,
				"context",
				"",
				"the configuration context to use to access your Hetzner account",
			)

			baseCmd.Flags().StringVar(
				&region,
				"region",
				"",
				"the region to use to access your Hetzner account",
			)
		},

		ProvideInitFeature: func() cloudProviderFeature[features.InitInput] {
			return dependencies.ProvideHetznerInitFeature(
				system.UserConfigDir(),
				region,
				context,
			)
		},

		ProvideEditFeature: func() cloudProviderFeature[features.EditInput] {
			return dependencies.ProvideHetznerEditFeature(
				system.UserConfigDir(),
				region,
				context,
			)
		},

		ProvideRemoveFeature: func() cloudProviderFeature[features.RemoveInput] {
			return dependencies.ProvideHetznerRemoveFeature(
				system.UserConfigDir(),
				region,
				context,
			)
		},

		ProvideServeFeature: func() cloudProviderFeature[features.ServeInput] {
			return dependencies.ProvideHetznerServeFeature(
				system.UserConfigDir(),
				region,
				context,
			)
		},

		ProvideUnserveFeature: func() cloudProviderFeature[features.UnserveInput] {
			return dependencies.ProvideHetznerUnserveFeature(
				system.UserConfigDir(),
				region,
				context,
			)
		},

		UninstallSuccessMessage:            "Eleven has been uninstalled from this region on this Hetzner account.",
		UninstallAlreadyUninstalledMessage: "Eleven is already uninstalled in this region on this Hetzner account.",

		ProvideUninstallFeature: func() cloudProviderFeature[features.UninstallInput] {
			return dependencies.ProvideHetznerUninstallFeature(
				system.UserConfigDir(),
				region,
				context,
			)
		},
	}

	return &hetzner
}
