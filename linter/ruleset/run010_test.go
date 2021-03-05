package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateRun010(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		ExampleName string
		CommandStr  string
		IsViolation bool
		DocsContext string
	}{
		{ExampleName: "", IsViolation: false, DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "apt-get -y --no-install-recommends install vim" },
		{ExampleName: "", IsViolation: false, DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "apt-get --yes --no-install-recommends install vim"},
		{CommandStr: "apt-get --assume-yes --no-install-recommends install vim", IsViolation: false},
		{CommandStr: "apt-get install -y vim", IsViolation:  true},
		{CommandStr: "DEBIAN_FRONTEND=noninteractive apt-get update", IsViolation: false},
		{CommandStr: "DEBIAN_FRONTEND=noninteractive apt     update", IsViolation: false},
		{CommandStr: "date", IsViolation: false},
		{CommandStr: "date; date", IsViolation: false},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.CommandStr, func(t *testing.T) {
			t.Parallel()

			// assemble command
			commandBody := instructions.ShellDependantCmdLine{CmdLine: []string{testCase.CommandStr}, PrependShell: true}
			runCommandWithoutSudo := &instructions.RunCommand{ShellDependantCmdLine: commandBody}

			// test validation rule
			assert.Equal(t, testCase.IsViolation, RuleSet.ValidateRun010(runCommandWithoutSudo).IsViolated())
		})
	}
}
