package ruleset

import (
	"path/filepath"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

var _ = NewRule("WKD001", "WORKDIR should be an absolute path for clarity and reliability.", "", ValWarning,
	ValidateWkd001)

func ValidateWkd001(workdirCommand *instructions.WorkdirCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(workdirCommand),
	}

	if !filepath.IsAbs(workdirCommand.Path) {
		result.SetViolated()
		result.LocationRange = ParseLocationFromRawParser(workdirCommand.Path, workdirCommand.Location())
	}

	return result
}
