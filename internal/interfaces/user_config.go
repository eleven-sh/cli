package interfaces

import (
	"github.com/eleven-sh/cli/internal/config"
	"github.com/eleven-sh/eleven/github"
)

type UserConfigManager interface {
	GetString(key config.UserConfigKey) string
	GetBool(key config.UserConfigKey) bool
	Set(key config.UserConfigKey, value interface{})
	WriteConfig() error
	PopulateFromGitHubUser(githubUser *github.AuthenticatedUser)
}
