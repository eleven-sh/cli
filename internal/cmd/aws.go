package cmd

import (
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/eleven-sh/cli/internal/dependencies"
	"github.com/eleven-sh/cli/internal/globals"
	"github.com/eleven-sh/eleven/features"
	"github.com/spf13/cobra"
)

func getAWSCloudProvider() *cloudProvider {
	var profile string
	var region string

	var credentialsFilePath string
	var configFilePath string

	var aws = cloudProvider{
		ShortName:  "aws",
		LongName:   "Amazon Web Services",
		GlobalName: globals.AWSCloudProvider,

		DefaultInstanceType: "t2.medium",
		ExampleInstanceType: "m4.large",

		AddFlagsToBaseCmd: func(baseCmd *cobra.Command) {
			baseCmd.Flags().StringVar(
				&profile,
				"profile",
				"",
				"the configuration profile to use to access your AWS account",
			)

			baseCmd.Flags().StringVar(
				&region,
				"region",
				"",
				"the region to use to access your AWS account",
			)

			credentialsFilePath = config.DefaultSharedCredentialsFilename()
			configFilePath = config.DefaultSharedConfigFilename()
		},

		ProvideInitFeature: func() cloudProviderFeature[features.InitInput] {
			return dependencies.ProvideAWSInitFeature(
				region,
				profile,
				credentialsFilePath,
				configFilePath,
			)
		},

		ProvideEditFeature: func() cloudProviderFeature[features.EditInput] {
			return dependencies.ProvideAWSEditFeature(
				region,
				profile,
				credentialsFilePath,
				configFilePath,
			)
		},

		ProvideRemoveFeature: func() cloudProviderFeature[features.RemoveInput] {
			return dependencies.ProvideAWSRemoveFeature(
				region,
				profile,
				credentialsFilePath,
				configFilePath,
			)
		},

		ProvideServeFeature: func() cloudProviderFeature[features.ServeInput] {
			return dependencies.ProvideAWSServeFeature(
				region,
				profile,
				credentialsFilePath,
				configFilePath,
			)
		},

		ProvideUnserveFeature: func() cloudProviderFeature[features.UnserveInput] {
			return dependencies.ProvideAWSUnserveFeature(
				region,
				profile,
				credentialsFilePath,
				configFilePath,
			)
		},

		UninstallSuccessMessage:            "Eleven has been uninstalled from this region on this AWS account.",
		UninstallAlreadyUninstalledMessage: "Eleven is already uninstalled in this region on this AWS account.",

		ProvideUninstallFeature: func() cloudProviderFeature[features.UninstallInput] {
			return dependencies.ProvideAWSUninstallFeature(
				region,
				profile,
				credentialsFilePath,
				configFilePath,
			)
		},
	}

	return &aws
}
