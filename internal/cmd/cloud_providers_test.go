package cmd

import (
	"testing"

	"github.com/eleven-sh/cli/internal/globals"
	"github.com/spf13/cobra"
)

func TestPopulateCurrentCloudProviderGlobals(t *testing.T) {
	testCases := []struct {
		test                            string
		cloudProvider                   globals.CloudProvider
		baseCmd                         *cobra.Command
		baseCmdFlags                    map[string]string
		childrenCmd                     *cobra.Command
		expectedGlobalCloudProvider     globals.CloudProvider
		expectedGlobalCloudProviderArgs string
	}{
		{
			test:                            "with no base command",
			cloudProvider:                   globals.AWSCloudProvider,
			baseCmd:                         &cobra.Command{},
			expectedGlobalCloudProvider:     globals.AWSCloudProvider,
			expectedGlobalCloudProviderArgs: "",
		},

		{
			test:                            "with base command without flags",
			cloudProvider:                   globals.HetznerCloudProvider,
			baseCmd:                         &cobra.Command{},
			childrenCmd:                     &cobra.Command{},
			expectedGlobalCloudProvider:     globals.HetznerCloudProvider,
			expectedGlobalCloudProviderArgs: "",
		},

		{
			test:          "with base command and one flags",
			cloudProvider: globals.AWSCloudProvider,
			baseCmd:       &cobra.Command{},
			baseCmdFlags: map[string]string{
				"flag1": "flag1_value",
			},
			childrenCmd:                     &cobra.Command{},
			expectedGlobalCloudProvider:     globals.AWSCloudProvider,
			expectedGlobalCloudProviderArgs: "--flag1 flag1_value",
		},

		{
			test:          "with base command and three flags",
			cloudProvider: globals.AWSCloudProvider,
			baseCmd:       &cobra.Command{},
			baseCmdFlags: map[string]string{
				"flag1": "flag1_value",
				"flag4": "flag4_value",
				// Changed set to "false". See below.
				"flag6": "",
			},
			childrenCmd:                     &cobra.Command{},
			expectedGlobalCloudProvider:     globals.AWSCloudProvider,
			expectedGlobalCloudProviderArgs: "--flag1 flag1_value --flag4 flag4_value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			cmdUnderTest := tc.baseCmd

			if tc.childrenCmd != nil {
				tc.baseCmd.AddCommand(tc.childrenCmd)
				cmdUnderTest = tc.childrenCmd
			}

			if len(tc.baseCmdFlags) > 0 {
				for flagName, flagValue := range tc.baseCmdFlags {
					tc.baseCmd.Flags().String(flagName, flagValue, "")

					if len(flagValue) == 0 {
						continue
					}

					tc.baseCmd.Flags().Lookup(flagName).Changed = true
				}
			}

			populateCurrentCloudProviderGlobals(
				tc.cloudProvider,
				cmdUnderTest,
			)

			if tc.expectedGlobalCloudProvider != globals.CurrentCloudProvider {
				t.Fatalf(
					"expected global cloud provider to equal '%s', got '%s'",
					tc.expectedGlobalCloudProvider,
					globals.CurrentCloudProvider,
				)
			}

			if tc.expectedGlobalCloudProviderArgs != globals.CurrentCloudProviderArgs {
				t.Fatalf(
					"expected global cloud provider args to equal '%s', got '%s'",
					tc.expectedGlobalCloudProviderArgs,
					globals.CurrentCloudProviderArgs,
				)
			}
		})
	}
}
