package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	Parser "github.com/cremindes/whalelint/parser"
)

var _ = NewRule("RUN007", "Use 'WORKDIR' to switch to a directory.", "", ValWarning, ValidateRun007)

func ValidateRun007(runCommand *instructions.RunCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(runCommand),
	}

	bashCommandList := Parser.ParseBashCommandList(runCommand)
	if bashCommandList == nil {
		return result
	}
	// RUN command starts with "cd" right away
	if bashCommandList[0].Bin() == "cd" {
		result.SetViolated()
		result.LocationRange = ParseLocationFromRawParser(bashCommandList[0].Bin(), runCommand.Location())
	} else if len(bashCommandList) >= 2 { // nolint:gomnd
		// RUN command starts with mkdir and then followed by a cd
		if bashCommandList[0].Bin() == "mkdir" && bashCommandList[1].Bin() == "cd" {
			result.SetViolated()
			result.LocationRange = ParseLocationFromRawParser(bashCommandList[0].Bin(), runCommand.Location())
		}
	}

	return result
}
