package cmd

import (
	"os"

	"github.com/eleven-sh/cli/internal/dependencies"
	"github.com/eleven-sh/cli/internal/features"
	"github.com/spf13/cobra"
)

// loginCmd represents the "eleven login" command
var loginCmd = &cobra.Command{
	Use: "login",

	Short: "Connect a GitHub account",

	Long: `Connect a GitHub account.

This command connects a GitHub account to use with Eleven.

Eleven requires the following permissions:

  - "Public SSH keys" and "Repositories" to let you access your repositories from your sandboxes

  - "Personal user data" to configure Git

All your data (including the OAuth access token) will only be stored locally.`,

	Example: "  eleven login",

	Run: func(cmd *cobra.Command, args []string) {
		loginInput := features.LoginInput{}

		login := dependencies.ProvideLoginFeature()

		if err := login.Execute(loginInput); err != nil {
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
