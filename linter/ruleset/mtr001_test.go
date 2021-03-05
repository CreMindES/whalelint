package ruleset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateMtr001(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		IsViolation   bool
		HasMaintainer bool
		ExampleName   string
		DocsContext   string
	}{
		{
			IsViolation:   true,
			ExampleName:   "Maintainer John Doe",
			HasMaintainer: true,
			DocsContext:   "`FROM` golang:1.16\n`MAINTAINER` John Doe <john.doe@example.com>",
		},
		{
			IsViolation:   false,
			ExampleName:   "No Maintainer",
			HasMaintainer: false,
			DocsContext:   "`FROM` golang:1.16",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.ExampleName, func(t *testing.T) {
			t.Parallel()

			if testCase.HasMaintainer {
				maintainerCommand, err := RuleSet.NewMaintainerCommand("John Doe <john.doe@example.com>")
				assert.Nil(t, err)

				assert.Equal(t, testCase.IsViolation, RuleSet.ValidateMtr001(maintainerCommand).IsViolated())
			} else {
				assert.Equal(t, testCase.IsViolation, testCase.HasMaintainer)
			}
		})
	}
}
