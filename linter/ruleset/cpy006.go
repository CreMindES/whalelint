package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	Utils "github.com/cremindes/whalelint/utils"
)

var _ = NewRule("CPY006", "COPY --from value should not be the same as the stage.", "", ValError,
	ValidateCpy006)

func ValidateCpy006(stage instructions.Stage) RuleValidationResult {
	result := RuleValidationResult{isViolated: false, LocationRange: LocationRange{}}

	for _, command := range stage.Commands {
		if copyCommand, ok := command.(*instructions.CopyCommand); ok {
			if copyCommand.From == stage.Name || copyCommand.From == stage.BaseName ||
				Utils.MatchDockerImageNames(copyCommand.From, stage.BaseName) {
				result.SetViolated()
				result.LocationRange = ParseLocationFromRawParser(copyCommand.From, copyCommand.Location())
			}
		}
	}

	return result
}
