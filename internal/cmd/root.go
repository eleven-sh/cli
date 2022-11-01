package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"

	"github.com/eleven-sh/cli/internal/config"
	"github.com/eleven-sh/cli/internal/dependencies"
	"github.com/eleven-sh/cli/internal/exceptions"
	"github.com/eleven-sh/cli/internal/system"
	"github.com/eleven-sh/eleven/github"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "eleven",

	Short: "Code sandboxes running in your cloud provider account",

	Long: `Eleven - Code sandboxes running in your cloud provider account

To begin, run the command "eleven login" to connect your GitHub account.	

From there, the most common workflow is:

  - eleven <cloud_provider> init <sandbox_name>   : to initialize a new sandbox

  - eleven <cloud_provider> edit <sandbox_name>   : to connect your preferred editor to a sandbox

  - eleven <cloud_provider> remove <sandbox_name> : to remove an existing sandbox`,

	TraverseChildren: true,

	Version: "v0.0.3",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(
		ensureElevenCLIRequirements,
		initializeElevenCLIConfig,
		ensureGitHubAccessTokenValidity,
	)
}

func ensureElevenCLIRequirements() {
	missingRequirements := []string{}

	sshCommand := "ssh"
	_, err := exec.LookPath(sshCommand)

	if err != nil {
		missingRequirements = append(
			missingRequirements,
			fmt.Sprintf(
				"OpenSSH client (looked for an \"%s\" command)",
				sshCommand,
			),
		)
	}

	if len(missingRequirements) > 0 {
		elevenViewableErrorBuilder := dependencies.ProvideElevenViewableErrorBuilder()
		baseView := dependencies.ProvideBaseView()

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
}

func initializeElevenCLIConfig() {
	configDir := system.UserConfigDir()
	configDirPerms := fs.FileMode(0700)

	// Ensure configuration dir exists
	err := os.MkdirAll(
		configDir,
		configDirPerms,
	)
	cobra.CheckErr(err)

	configFilePath := system.UserConfigFilePath()
	configFilePerms := fs.FileMode(0600)

	// Ensure configuration file exists
	f, err := os.OpenFile(
		configFilePath,
		os.O_CREATE,
		configFilePerms,
	)
	cobra.CheckErr(err)
	defer f.Close()

	viper.SetConfigFile(configFilePath)
	cobra.CheckErr(viper.ReadInConfig())
}

// ensureGitHubAccessTokenValidity ensures that
// the github access token has not been
// revoked by user
func ensureGitHubAccessTokenValidity() {
	userConfig := config.NewUserConfig()
	userIsLoggedIn := userConfig.GetBool(config.UserConfigKeyUserIsLoggedIn)

	if !userIsLoggedIn {
		return
	}

	gitHubService := github.NewService()

	githubUser, err := gitHubService.GetAuthenticatedUser(
		userConfig.GetString(
			config.UserConfigKeyGitHubAccessToken,
		),
	)

	if err != nil &&
		gitHubService.IsInvalidAccessTokenError(err) { // User has revoked access token

		userIsLoggedIn = false

		userConfig.Set(
			config.UserConfigKeyUserIsLoggedIn,
			userIsLoggedIn,
		)

		// Error is swallowed here to
		// not confuse user with unexpected error
		_ = userConfig.WriteConfig()
	}

	if err == nil {
		// Update config with updated values from GitHub
		userConfig.PopulateFromGitHubUser(githubUser)

		// Error is swallowed here to
		// not confuse user with unexpected error
		_ = userConfig.WriteConfig()
	}
}

func ensureUserIsLoggedIn() {
	userConfig := config.NewUserConfig()
	userIsLoggedIn := userConfig.GetBool(config.UserConfigKeyUserIsLoggedIn)

	if !userIsLoggedIn {
		elevenViewableErrorBuilder := dependencies.ProvideElevenViewableErrorBuilder()
		baseView := dependencies.ProvideBaseView()

		baseView.ShowErrorViewWithStartingNewLine(
			elevenViewableErrorBuilder.Build(
				exceptions.ErrUserNotLoggedIn,
			),
		)

		os.Exit(1)
	}
}
