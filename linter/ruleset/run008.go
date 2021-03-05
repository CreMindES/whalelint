package ruleset

import (
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	Parser "github.com/cremindes/whalelint/parser"
)

var _ = NewRule("RUN008", "Prefer apt-get over apt as the latter does not have a stable CLI.", "", ValWarning,
	ValidateRun008)

func ValidateRun008(runCommand *instructions.RunCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(runCommand),
	}

	bin := "apt"

	bashCommandList := Parser.ParseBashCommandList(runCommand.CmdLine)
	for _, bashCommand := range bashCommandList {
		if bashCommand.Bin() == bin {
			result.SetViolated()
			// location
			result.LocationRange.start.charNumber = strings.Index(bashCommand.String(), bin)
			result.LocationRange.end.charNumber = strings.LastIndex(bashCommand.String(), bin)
		}
	}

	return result
}
