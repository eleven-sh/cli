package aws

import (
	"errors"
	"testing"

	"github.com/eleven-sh/aws-cloud-provider/userconfig"
	"github.com/eleven-sh/eleven/entities"
)

func TestViewableErrorBuilder(t *testing.T) {
	testCases := []struct {
		test                       string
		passedError                error
		expectedViewableErrorTitle string
	}{
		{
			test:                       "with unknown error",
			passedError:                errors.New(""),
			expectedViewableErrorTitle: "Unknown error",
		},

		{
			test:                       "with Eleven error",
			passedError:                entities.ErrEditCreatingEnv{},
			expectedViewableErrorTitle: "Invalid sandbox state",
		},

		{
			test:                       "with AWS error",
			passedError:                userconfig.ErrMissingAccessKeyInEnv,
			expectedViewableErrorTitle: "Missing environment variable",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			builder := NewAWSViewableErrorBuilder()

			err := builder.Build(tc.passedError)

			if err.Title != tc.expectedViewableErrorTitle {
				t.Fatalf(
					"expected viewable error title to equal '%s', got '%s'",
					tc.expectedViewableErrorTitle,
					err.Title,
				)
			}
		})
	}
}
