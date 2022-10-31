package entities

import (
	"errors"
	"reflect"
	"testing"

	"github.com/eleven-sh/eleven/entities"
)

func TestResolveRepositories(t *testing.T) {
	testCases := []struct {
		test                  string
		repositories          []string
		githubUsername        string
		expectedResolvedRepos []entities.EnvRepository
		expectedError         error
	}{
		{
			test: "with empty repositories",

			repositories: []string{},

			githubUsername: "jeremylevy",

			expectedResolvedRepos: []entities.EnvRepository{},
		},

		{
			test: "with valid repositories",

			repositories: []string{
				"test",
				"eleven-sh/test",
				"https://github.com/recode-sh/test",
				"git@github.com:recode-sh/agent.git",
				"https://github.com/recode-sh/api/blob/master/src/BaseCommand.ts",
				"https://github.com/recode-sh/workspace.git",
			},

			githubUsername: "jeremylevy",

			expectedResolvedRepos: []entities.EnvRepository{
				{
					Name:          "test",
					Owner:         "jeremylevy",
					ExplicitOwner: false,
					GitURL:        "git@github.com:jeremylevy/test.git",
					GitHTTPURL:    "https://github.com/jeremylevy/test.git",
				},

				{
					Name:          "test",
					Owner:         "eleven-sh",
					ExplicitOwner: true,
					GitURL:        "git@github.com:eleven-sh/test.git",
					GitHTTPURL:    "https://github.com/eleven-sh/test.git",
				},

				{
					Name:          "test",
					Owner:         "recode-sh",
					ExplicitOwner: true,
					GitURL:        "git@github.com:recode-sh/test.git",
					GitHTTPURL:    "https://github.com/recode-sh/test.git",
				},

				{
					Name:          "agent",
					Owner:         "recode-sh",
					ExplicitOwner: true,
					GitURL:        "git@github.com:recode-sh/agent.git",
					GitHTTPURL:    "https://github.com/recode-sh/agent.git",
				},

				{
					Name:          "api",
					Owner:         "recode-sh",
					ExplicitOwner: true,
					GitURL:        "git@github.com:recode-sh/api.git",
					GitHTTPURL:    "https://github.com/recode-sh/api.git",
				},

				{
					Name:          "workspace",
					Owner:         "recode-sh",
					ExplicitOwner: true,
					GitURL:        "git@github.com:recode-sh/workspace.git",
					GitHTTPURL:    "https://github.com/recode-sh/workspace.git",
				},
			},
		},

		{
			test: "with duplicated repositories",

			repositories: []string{
				"test",
				"eleven-sh/api",
				"https://github.com/jeremylevy/test",
				"git@github.com:recode-sh/agent.git",
				"https://github.com/recode-sh/api/blob/master/src/BaseCommand.ts",
				"https://github.com/recode-sh/workspace.git",
			},

			githubUsername: "jeremylevy",

			expectedError: entities.ErrEnvDuplicatedRepositories{
				RepoName:  "test",
				RepoOwner: "jeremylevy",
			},
		},

		{
			test: "with invalid repositories",

			repositories: []string{
				"test",
				"eleven-sh/test",
				"/test",
				"https://github.com/recode-sh/test",
				"git@github.com:recode-sh/agent.git",
				"https://github.com/recode-sh/api/blob/master/src/BaseCommand.ts",
				"https://github.com/recode-sh/workspace.git",
			},

			githubUsername: "jeremylevy",

			expectedError: entities.ErrEnvRepositoryNotFound{
				RepoName:  "/test",
				RepoOwner: "jeremylevy",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			resolvedRepos, err := resolveRepositories(
				tc.repositories,
				tc.githubUsername,
			)

			if err != nil && tc.expectedError == nil {
				t.Fatalf("expected no error, got '%+v'", err)
			}

			if err == nil && tc.expectedError != nil {
				t.Fatalf(
					"expected error to equal '%+v', got nothing",
					tc.expectedError,
				)
			}

			if tc.expectedError != nil &&
				errors.As(
					tc.expectedError,
					&entities.ErrEnvDuplicatedRepositories{},
				) {

				typedExpectedError := tc.expectedError.(entities.ErrEnvDuplicatedRepositories)
				typedError, ok := err.(entities.ErrEnvDuplicatedRepositories)

				if !ok {
					t.Fatalf(
						"expected env duplicated repositories error, got '%+v'",
						err,
					)
				}

				if typedError.RepoOwner != typedExpectedError.RepoOwner {
					t.Fatalf(
						"expected error repository owner to equal '%s', got '%s'",
						typedExpectedError.RepoOwner,
						typedError.RepoOwner,
					)
				}

				if typedError.RepoName != typedExpectedError.RepoName {
					t.Fatalf(
						"expected error repository name to equal '%s', got '%s'",
						typedExpectedError.RepoName,
						typedError.RepoName,
					)
				}
			}

			if tc.expectedError != nil &&
				errors.As(
					tc.expectedError,
					&entities.ErrEnvRepositoryNotFound{},
				) {

				typedExpectedError := tc.expectedError.(entities.ErrEnvRepositoryNotFound)
				typedError, ok := err.(entities.ErrEnvRepositoryNotFound)

				if !ok {
					t.Fatalf(
						"expected env repository not found error, got '%+v'",
						err,
					)
				}

				if typedError.RepoOwner != typedExpectedError.RepoOwner {
					t.Fatalf(
						"expected error repository owner to equal '%s', got '%s'",
						typedExpectedError.RepoOwner,
						typedError.RepoOwner,
					)
				}

				if typedError.RepoName != typedExpectedError.RepoName {
					t.Fatalf(
						"expected error repository name to equal '%s', got '%s'",
						typedExpectedError.RepoName,
						typedError.RepoName,
					)
				}
			}

			if tc.expectedError != nil {
				return
			}

			if !reflect.DeepEqual(tc.expectedResolvedRepos, resolvedRepos) {
				t.Fatalf(
					"expected resolved repositories to equal '%+v', got '%+v'",
					tc.expectedResolvedRepos,
					resolvedRepos,
				)
			}
		})
	}
}
