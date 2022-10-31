package entities

import (
	"github.com/eleven-sh/agent/proto"
	"github.com/eleven-sh/cli/internal/config"
	"github.com/eleven-sh/cli/internal/interfaces"
	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/eleven/github"
)

type EnvRepositoriesResolver struct {
	logger     interfaces.Logger
	userConfig interfaces.UserConfigManager
	github     interfaces.GitHubManager
}

func NewEnvRepositoriesResolver(
	logger interfaces.Logger,
	userConfig interfaces.UserConfigManager,
	github interfaces.GitHubManager,
) EnvRepositoriesResolver {

	return EnvRepositoriesResolver{
		logger:     logger,
		userConfig: userConfig,
		github:     github,
	}
}

func (e EnvRepositoriesResolver) Resolve(
	repositories []string,
	checkForRepositoryExistence bool,
) ([]entities.EnvRepository, error) {

	githubUsername := e.userConfig.GetString(
		config.UserConfigKeyGitHubUsername,
	)

	resolvedRepos, err := resolveRepositories(
		repositories,
		githubUsername,
	)

	if err != nil {
		return nil, err
	}

	if checkForRepositoryExistence {

		githubAccessToken := e.userConfig.GetString(
			config.UserConfigKeyGitHubAccessToken,
		)

		for _, repository := range resolvedRepos {
			repoExists, err := e.github.DoesRepositoryExist(
				githubAccessToken,
				repository.Owner,
				repository.Name,
			)

			if err != nil {
				return nil, err
			}

			if !repoExists {
				return nil, entities.ErrEnvRepositoryNotFound{
					RepoOwner: repository.Owner,
					RepoName:  repository.Name,
				}
			}
		}
	}

	return resolvedRepos, nil
}

func resolveRepositories(
	repositories []string,
	githubUsername string,
) ([]entities.EnvRepository, error) {

	resolvedRepos := []entities.EnvRepository{}
	alreadyResolvedRepos := map[string]bool{}

	for _, repositoryName := range repositories {
		parsedRepoName, err := github.ParseRepositoryName(
			repositoryName,
			githubUsername,
		)

		if err != nil {
			// If the repository name is invalid, we are sure
			// that the repository doesn't exist
			return nil, entities.ErrEnvRepositoryNotFound{
				RepoOwner: githubUsername,
				RepoName:  repositoryName,
			}
		}

		repoUniqueKey := parsedRepoName.Owner + parsedRepoName.Name

		if _, alreadyResolved := alreadyResolvedRepos[repoUniqueKey]; alreadyResolved {
			return nil, entities.ErrEnvDuplicatedRepositories{
				RepoOwner: parsedRepoName.Owner,
				RepoName:  parsedRepoName.Name,
			}
		}

		alreadyResolvedRepos[repoUniqueKey] = true

		resolvedRepos = append(resolvedRepos, entities.EnvRepository{
			Owner:         parsedRepoName.Owner,
			ExplicitOwner: parsedRepoName.ExplicitOwner,

			Name: parsedRepoName.Name,

			GitURL:     github.BuildGitURL(parsedRepoName),
			GitHTTPURL: github.BuildGitHTTPURL(parsedRepoName),
		})
	}

	return resolvedRepos, nil
}

func BuildProtoRepositoriesFromEnv(
	env *entities.Env,
) []*proto.EnvRepository {

	repositories := []*proto.EnvRepository{}

	for _, repository := range env.Repositories {
		repositories = append(repositories, &proto.EnvRepository{
			Name:  repository.Name,
			Owner: repository.Owner,
		})
	}

	return repositories
}
