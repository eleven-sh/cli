package interfaces

import (
	cliGitHub "github.com/eleven-sh/eleven/github"
	"github.com/google/go-github/v43/github"
)

type GitHubManager interface {
	GetAuthenticatedUser(accessToken string) (*cliGitHub.AuthenticatedUser, error)
	DoesRepositoryExist(accessToken, repositoryOwner, repositoryName string) (bool, error)

	CreateSSHKey(accessToken, keyPairName, publicKeyContent string) (*github.Key, error)
	RemoveSSHKey(accessToken string, sshKeyID int64) error

	IsNotFoundError(err error) bool
}
