package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

// TestValidateRun004 tests RUN command validation with and without "sudo".
//
// Scenario: Shell Command validation
//
// GIVEN a shell command
// WHEN that command
//   - does not have sudo in it
//   - starts with sudo
//   - has sudo somewhere in it after multiple environment variable assigment
//   - has sudo in second command - TODO
// THEN this should
//   - pass
//   - fail
//   - fail
//   - fail
func TestValidateRun004(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		command   string
		violation bool
	}{
		{
			"echo 'ok'",
			false,
		},
		{
			"sudo echo 'not ok'",
			true,
		},
		{
			"myEnvVar=2 myOtherEnvVar=\"$test\" sudo echo 'not ok'",
			true,
		},
		{
			"date; sudo echo 'not ok'",
			true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.command, func(t *testing.T) {
			t.Parallel()

			commandBody := instructions.ShellDependantCmdLine{CmdLine: []string{testCase.command}, PrependShell: true}
			runCommandWithoutSudo := &instructions.RunCommand{ShellDependantCmdLine: commandBody}
			result := RuleSet.ValidateRun004(runCommandWithoutSudo).IsViolated()
			assert.Equal(t, result, testCase.violation)
		})
	}
}
