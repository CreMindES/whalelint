package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	"github.com/cremindes/whalelint/utils"
)

// STS -> Stage Single.
var _ = NewRule("STS002", "Stage name \"latest\" is prone to future errors.", "TODO", ValWarning, ValidateSts002)

func ValidateSts002(stage instructions.Stage) RuleValidationResult {
	result := RuleValidationResult{isViolated: false, LocationRange: CopyLocationRange(stage.Location)}

	image, tag := utils.SplitKeyValue(stage.BaseName, ':')
	result.SetViolated(tag == "latest") // as latest is the default

	if result.IsViolated() {
		result.message = "Image \"" + image + "\" should not use \"latest\" as tag."
		result.LocationRange = ParseLocationFromRawParser(stage.BaseName, stage.Location)
	}

	return result
}
