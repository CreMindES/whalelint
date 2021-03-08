package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

// nolint:funlen
func TestValidateRun002(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		CommandStr  string
		IsViolation bool
		ExampleName string
		DocsContext string
	}{
		{
			CommandStr:  "apt-get install vim=1.12.1",
			IsViolation: false,
			ExampleName: "Deb package install specific version.",
			DocsContext: "FROM ubuntu:20.04\nRUN {{ .CommandStr }}",
		},
		{
			CommandStr:  "apt-get install vim",
			IsViolation: true,
			ExampleName: "Deb package install.",
			DocsContext: "FROM ubuntu:20.04\nRUN {{ .CommandStr }}",
		},
		{
			CommandStr:  "apt install vim",
			IsViolation: true,
			ExampleName: "Deb package install with apt.",
			DocsContext: "FROM ubuntu:20.04\nRUN {{ .CommandStr }}",
		},
		{
			CommandStr:  "apt update && apt install vim",
			IsViolation: true,
			ExampleName: "Apt update and deb package install with apt.",
			DocsContext: "FROM ubuntu:20.04\nRUN {{ .CommandStr }}",
		},
		{
			CommandStr:  "DEBIAN_FRONTEND=noninteractive apt-get update",
			IsViolation: false,
			ExampleName: "deb package repository update, non-interactive env set.",
			DocsContext: "FROM ubuntu:20.04\nRUN {{ .CommandStr }}",
		},
		{
			CommandStr:  "DEBIAN_FRONTEND=noninteractive apt-get install -y gedit vim=1.12.2",
			IsViolation: true,
			ExampleName: "Multiple deb package install, with and without specific version, non-interactive env set.",
			DocsContext: "FROM ubuntu:20.04\nRUN {{ .CommandStr }}",
		},
		{
			CommandStr:  "pip install --no-cache-dir -r requirements.txt",
			IsViolation: false,
			ExampleName: "Install pip packages from requirements file.",
			DocsContext: "FROM ubuntu:20.04\nRUN {{ .CommandStr }}",
		},

		{
			CommandStr:  "date",
			IsViolation: false,
			ExampleName: "Unrelated command.",
			DocsContext: "FROM ubuntu:20.04\nRUN {{ .CommandStr }}",
		},
	}

	RuleSet.RegisterTestCaseDocs("RUN002", testCases)

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.ExampleName, func(t *testing.T) {
			t.Parallel()

			// assemble command
			commandBody := instructions.ShellDependantCmdLine{CmdLine: []string{testCase.CommandStr}, PrependShell: true}
			runCommandWithoutSudo := &instructions.RunCommand{ShellDependantCmdLine: commandBody}

			// test validation rule
			assert.Equal(t, testCase.IsViolation, RuleSet.ValidateRun002(runCommandWithoutSudo).IsViolated())
		})
	}
}
