package hetzner

import (
	"errors"
	"testing"

	"github.com/eleven-sh/eleven/entities"
	"github.com/eleven-sh/hetzner-cloud-provider/userconfig"
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
			test:                       "with Hetzner error",
			passedError:                userconfig.ErrMissingRegionInEnv,
			expectedViewableErrorTitle: "Missing region",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			builder := NewHetznerViewableErrorBuilder()

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
