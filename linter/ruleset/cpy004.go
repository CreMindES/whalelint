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

	sourceCount := len(copyCommand.SourcesAndDest.SourcePaths)
	// in case of CPY002 violation, the flag can end up in the sources list
	for _, src := range copyCommand.SourcesAndDest.SourcePaths {
		if strings.HasPrefix(src, "-") {
			sourceCount--
		}
	}

	if sourceCount > 1 {
		// is valid
		destination := copyCommand.SourcesAndDest.DestPath
		destinationLastChar := destination[len(destination)-1]
		result.SetViolated(destinationLastChar != '/')
		// location
		// note: prefixing the destination with a space in order to avoid the edge case, where the destination can be
		//       found in the source as well as a substring. This prefix need to be cut off, that's why the increment at
		//       the end.
		result.LocationRange = ParseLocationFromRawParser(" "+destination, copyCommand.Location())
		result.LocationRange.start.charNumber++
	}

	return result
}
