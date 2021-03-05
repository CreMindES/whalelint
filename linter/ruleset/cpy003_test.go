package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateCpy003(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		ChownValue  string
		IsViolation bool
		ExampleName string
		DocsContext string
	}{ // valid examples are from https://docs.docker.com/engine/reference/builder/#copy
		{ChownValue: "55:mygroup", IsViolation: false, ExampleName: "COPY with chown=55:mygroup",
		 	DocsContext: "FROM golang 1.15\\nCOPY --chown={{ .ChownValue }} src dst"},
		{ChownValue: "bin"       , IsViolation: false, ExampleName: "COPY with chown=bin",
		 	DocsContext: "FROM golang 1.15\\nCOPY --chown={{ .ChownValue }} src dst"},
		{ChownValue: "1"         , IsViolation: false, ExampleName: "COPY with chown=1",
		 	DocsContext: "FROM golang 1.15\\nCOPY --chown={{ .ChownValue }} src dst"},
		{ChownValue: "10:11"     , IsViolation: false, ExampleName: "COPY with chown=10:11",
			DocsContext: "FROM golang 1.15\\nCOPY --chown={{ .ChownValue }} src dst"},
		{ChownValue: "10;11"     , IsViolation:  true, ExampleName: "COPY with chown=10;11",
			DocsContext: "FROM golang 1.15\\nCOPY --chown={{ .ChownValue }} src dst"},
		{ChownValue: "10,11"     , IsViolation:  true, ExampleName: "COPY with chown=10,11",
			DocsContext: "FROM golang 1.15\\nCOPY --chown={{ .ChownValue }} src dst"},
		{ChownValue: "$$"        , IsViolation:  true, ExampleName: "COPY with chown=$$",
			DocsContext: "FROM golang 1.15\\nCOPY --chown={{ .ChownValue }} src dst"},
		{ChownValue: "55:11,22"  , IsViolation:  true, ExampleName: "COPY with chown=55:11,22",
			DocsContext: "FROM golang 1.15\\nCOPY --chown={{ .ChownValue }} src dst"},
		{ChownValue: "55:11 22"  , IsViolation:  true, ExampleName: "COPY with chown=55:11 22",
			DocsContext: "FROM golang 1.15\\nCOPY --chown={{ .ChownValue }} src dst"},
	}

	RuleSet.RegisterTestCaseDocs("CPY003", testCases)

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.ExampleName, func(t *testing.T) {
			t.Parallel()

			command := &instructions.CopyCommand{
				SourcesAndDest: []string{},
				From:           "",
				Chown:          testCase.ChownValue,
				Chmod:          "",
			}

			assert.Equal(t, testCase.IsViolation, RuleSet.ValidateCpy003(command).IsViolated())
		})
	}
}
