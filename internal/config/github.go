package config

var (
	GitHubOAuthClientID = "523ada8718f092d34c40"

	GitHubOAuthCLIToAPIURLPath = "/github/oauth/callback"
	GitHubOAuthCLIToAPIURL     = "http://127.0.0.1:8080" + GitHubOAuthCLIToAPIURLPath

	GitHubOAuthAPIToCLIURLPath = "/github/oauth/callback"

	GitHubOAuthScopes = []string{
		"read:user",
		"user:email",
		"repo",
		"admin:public_key",
	}
)
