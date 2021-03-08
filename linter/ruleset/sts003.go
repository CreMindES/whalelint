package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

// STS -> Stage Single.
var _ = NewRule("STS003", "Platform should be specified in build tool and not FROM.", "TODO",
	ValWarning, ValidateSts003)

func ValidateSts003(stage instructions.Stage) RuleValidationResult {
	result := RuleValidationResult{isViolated: false, LocationRange: BKRangeSliceToLocationRange(stage.Location)}

	result.SetViolated(len(stage.Platform) > 0)

	if result.IsViolated() {
		result.message = "Specifying platform at build tool level gives more flexibility."
		result.LocationRange = ParseLocationFromRawParser(stage.Platform, stage.Location)
	}

	return result
}
