package ruleset

import (
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

var _ = NewRule("CPY004", "COPY with more than one source requires the destination to end with \"/\".", "",
	ValError, ValidateCpy004)

func ValidateCpy004(copyCommand *instructions.CopyCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(copyCommand),
	}

	sourceCount := len(copyCommand.SourcesAndDest.Sources())
	// in case of CPY002 violation, the flag can end up in the sources list
	for _, src := range copyCommand.SourcesAndDest {
		if strings.HasPrefix(src, "-") {
			sourceCount--
		}
	}

	if sourceCount > 1 {
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
