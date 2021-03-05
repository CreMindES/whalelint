package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	"github.com/cremindes/whalelint/utils"
)

// STS -> Stage Single.
var _ = NewRule("STS001", "Stage name should have an explicit tag..", "", ValWarning, ValidateSts001)

func ValidateSts001(stage instructions.Stage) RuleValidationResult {
	result := RuleValidationResult{isViolated: false, LocationRange: CopyLocationRange(stage.Location)}

	image, tag := utils.SplitKeyValue(stage.BaseName, ':')
	result.SetViolated(tag == "")

	if result.IsViolated() {
		result.message = "Image \"" + image + "\" should have an explicit tag."
		result.LocationRange = ParseLocationFromRawParser(stage.BaseName, stage.Location)
	}

	return result
}
