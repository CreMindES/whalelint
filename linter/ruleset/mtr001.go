package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

var _ = NewRule("MTR001", "MAINTAINER is deprecated. Use a LABEL instead.", "",
	ValDeprecation, ValidateMtr001)

func ValidateMtr001(maintainerCommand *instructions.MaintainerCommand) RuleValidationResult {
	return RuleValidationResult{
		isViolated:    true,
		LocationRange: ParseLocationFromRawParser(maintainerCommand.String(), maintainerCommand.Location()),
	}
}
