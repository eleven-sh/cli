//go:build prod

package config

func init() {
	GitHubOAuthClientID = "eb67d11a0e1cf1824f09"

	GitHubOAuthCLIToAPIURL = "https://api.eleven.sh" + GitHubOAuthCLIToAPIURLPath
}
