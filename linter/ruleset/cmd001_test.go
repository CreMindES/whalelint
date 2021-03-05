// nolint:dupl
package ruleset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateCmd001(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExampleName string
		IsViolation bool
		CmdStr      string
		DocsContext string
	}{
		{
			ExampleName: "Proper CMD command in exec JSON format.",
			IsViolation: false,
			CmdStr:      "[\"/bin/bash\", \"date\"]",
			DocsContext: "FROM golang 1.16\nCMD {{ .CmdStr }}",
		},
		{
			ExampleName: "Proper CMD command in shell format.",
			IsViolation: true,
			CmdStr:      "/bin/bash date",
			DocsContext: "FROM golang 1.16\nCMD {{ .CmdStr }}",
		},
		{
			ExampleName: "Proper CMD command in invalid format with 2 args.",
			IsViolation: true,
			CmdStr:      "[/bin/bash date]",
			DocsContext: "FROM golang 1.16\nCMD {{ .CmdStr }}",
		},
		{
			ExampleName: "Proper CMD command in shell format.",
			IsViolation: true,
			CmdStr:      "date",
			DocsContext: "FROM golang 1.16\nCMD {{ .CmdStr }}",
		},
		{
			ExampleName: "Proper CMD command in invalid format with 1 arg.",
			IsViolation: true,
			CmdStr:      "[date]",
			DocsContext: "FROM golang 1.16\nCMD {{ .CmdStr }}",
		},
		{
			ExampleName: "Proper CMD command in exec JSON format, but missing a comma.",
			IsViolation: true,
			CmdStr:      "[\"/bin/bash\" \"date\"",
			DocsContext: "FROM golang 1.16\nCMD {{ .CmdStr }}",
		},
	}

	RuleSet.RegisterTestCaseDocs("ENT001", testCases)

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.ExampleName, func(t *testing.T) {
			t.Parallel()

			cmdCommand, err := RuleSet.NewCmdCommand(testCase.CmdStr, 2)
			assert.Nil(t, err)

			assert.Equal(t, testCase.IsViolation, RuleSet.ValidateCmd001(cmdCommand).IsViolated())
		})
	}
}
