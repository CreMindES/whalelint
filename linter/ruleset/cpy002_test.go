package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateCpy002(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		ChmodValue  string
		IsViolation bool
		ExampleName string
		DocsContext string
	}{
		{ChmodValue: "7440", IsViolation: false, ExampleName: "COPY with chmod=7440",
			DocsContext: "FROM golang 1.15\nCOPY --chmod={{ .ChmodValue }} src dst"},
		{ChmodValue: "644", IsViolation: false, ExampleName: "COPY with chmod=644",
			DocsContext: "FROM golang 1.15\nCOPY --chmod={{ .ChmodValue }} src dst"},
		{ChmodValue: "88", IsViolation:  true, ExampleName: "COPY with chmod=88",
			DocsContext: "FROM golang 1.15\nCOPY --chmod={{ .ChmodValue }} src dst"},
		{ChmodValue: "7780", IsViolation:  true, ExampleName: "COPY with chmod=7780",
			DocsContext: "FROM golang 1.15\nCOPY --chmod={{ .ChmodValue }} src dst"},
	}

	RuleSet.RegisterTestCaseDocs("CPY002", testCases)

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.ExampleName, func(t *testing.T) {
			t.Parallel()

			command := &instructions.CopyCommand{
				SourcesAndDest: instructions.SourcesAndDest{},
				From:           "",
				Chown:          "",
				Chmod:          testCase.ChmodValue,
			}

			assert.Equal(t, testCase.IsViolation, RuleSet.ValidateCpy002(command).IsViolated())
		})
	}
}
