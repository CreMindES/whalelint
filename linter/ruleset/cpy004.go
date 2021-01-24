package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

var _ = NewRule("CPY004", "COPY with more than one source requires the destination to end with \"/\".",
	ValError, ValidateCpy004)

func ValidateCpy004(copyCommand *instructions.CopyCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(copyCommand),
	}

	if len(copyCommand.SourcesAndDest.Sources()) > 1 {
		// is valid
		destination := copyCommand.SourcesAndDest.Dest()
		destinationLastChar := destination[len(destination)-1]
		result.SetViolated(destinationLastChar != '/')

		// further narrow location
		lineLength := len(copyCommand.String())
		result.LocationRange.start.charNumber = lineLength - len(destination)
		result.LocationRange.end.charNumber = lineLength
	}

	return result
}
