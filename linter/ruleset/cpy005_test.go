package ruleset_test

import (
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateCpy005(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		CommandParam string
		IsViolation  bool
		ExampleName  string
		DocsContext  string
	}{
		{
			CommandParam: "foo/bar /tmp/",
			IsViolation:  false,
			ExampleName:  "Standard COPY.",
			DocsContext:  "FROM golang:1.15\nCOPY {{ .CommandParam }}",
		},
		{
			CommandParam: "foo/bar.tar.gz /tmp/",
			IsViolation:  true,
			ExampleName:  "COPY \".tar.gz\"",
			DocsContext:  "FROM golang:1.15\nCOPY {{ .CommandParam }}",
		},
	}

	RuleSet.RegisterTestCaseDocs("CPY005", testCases)

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.ExampleName, func(t *testing.T) {
			t.Parallel()

			command := &instructions.CopyCommand{
				SourcesAndDest: strings.Fields(testCase.CommandParam),
				From:           "",
				Chown:          "",
				Chmod:          "",
			}

			assert.Equal(t, testCase.IsViolation, RuleSet.ValidateCpy005(command).IsViolated())
		})
	}
}
