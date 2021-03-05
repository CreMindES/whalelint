// nolint:dupl
package ruleset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateEnt001(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExampleName   string
		IsViolation   bool
		EntrypointStr string
		DocsContext   string
	}{
		{
			ExampleName:   "Proper ENTRYPOINT command in exec JSON format.",
			IsViolation:   false,
			EntrypointStr: "[\"/bin/bash\", \"date\"]",
			DocsContext:   "FROM golang 1.16\nENTRYPOINT {{ .EntrypointStr }}",
		},
		{
			ExampleName:   "Proper ENTRYPOINT command in shell format.",
			IsViolation:   true,
			EntrypointStr: "/bin/bash date",
			DocsContext:   "FROM golang 1.16\nENTRYPOINT {{ .EntrypointStr }}",
		},
		{
			ExampleName:   "Proper ENTRYPOINT command in invalid format with 2 args.",
			IsViolation:   true,
			EntrypointStr: "[/bin/bash date]",
			DocsContext:   "FROM golang 1.16\nENTRYPOINT {{ .EntrypointStr }}",
		},
		{
			ExampleName:   "Proper ENTRYPOINT command in shell format.",
			IsViolation:   true,
			EntrypointStr: "date",
			DocsContext:   "FROM golang 1.16\nENTRYPOINT {{ .EntrypointStr }}",
		},
		{
			ExampleName:   "Proper ENTRYPOINT command in invalid format with 1 arg.",
			IsViolation:   true,
			EntrypointStr: "[date]",
			DocsContext:   "FROM golang 1.16\nENTRYPOINT {{ .EntrypointStr }}",
		},
		{
			ExampleName:   "Proper ENTRYPOINT command in exec JSON format, but missing a comma.",
			IsViolation:   true,
			EntrypointStr: "[\"/bin/bash\" \"date\"",
			DocsContext:   "FROM golang 1.16\nENTRYPOINT {{ .EntrypointStr }}",
		},
	}

	RuleSet.RegisterTestCaseDocs("ENT001", testCases)

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.ExampleName, func(t *testing.T) {
			t.Parallel()

			entrypointCommand, err := RuleSet.NewEntrypointCommand(testCase.EntrypointStr, 2)
			assert.Nil(t, err)

			assert.Equal(t, testCase.IsViolation, RuleSet.ValidateEnt001(entrypointCommand).IsViolated())
		})
	}
}
