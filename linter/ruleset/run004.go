package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	Parser "github.com/cremindes/whalelint/parser"
)

var _ = NewRule(
	"RUN004",
	"Do not use sudo as it leads to unpredictable behavior. Use a tool like gosu to enforce root.", "",
	ValWarning,
	ValidateRun004)

func ValidateRun004(runCommand *instructions.RunCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(runCommand),
	}

	bashCommandList := Parser.ParseBashCommandList(runCommand.CmdLine)
	for _, bashCommand := range bashCommandList {
		if bashCommand.HasSudo() {
			result.SetViolated()
			result.LocationRange = ParseLocationFromRawParser("sudo", runCommand.Location())
		}
	}

	return result
}
