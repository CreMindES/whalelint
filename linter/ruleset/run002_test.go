package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateRun002(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		commandStr  string
		isViolation bool
	}{
		{commandStr: "apt-get install vim=1.12.1", isViolation: false},
		{commandStr: "apt-get install vim", isViolation:  true},
		{commandStr: "apt install vim", isViolation:  true},
		{commandStr: "apt update && apt install vim", isViolation:  true},
		{commandStr: "DEBIAN_FRONTEND=noninteractive apt-get update", isViolation: false},
		{commandStr: "DEBIAN_FRONTEND=noninteractive apt-get install -y gedit vim=1.12.2", isViolation: true},
		{commandStr: "date", isViolation: false},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.commandStr, func(t *testing.T) {
			t.Parallel()

			// assemble command
			commandBody := instructions.ShellDependantCmdLine{CmdLine: []string{testCase.commandStr}, PrependShell: true}
			runCommandWithoutSudo := &instructions.RunCommand{ShellDependantCmdLine: commandBody}

			// test validation rule
			assert.Equal(t, testCase.isViolation, RuleSet.ValidateRun002(runCommandWithoutSudo).IsViolated())
		})
	}
}
