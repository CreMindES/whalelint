package ruleset

import (
	"regexp"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

var _ = NewRule("CPY001", "COPY format is --[chmod|chown|from]=... srcList... dest|destDir",
	ValError, ValidateCpy001)

// checks COPY options format for obvious errors
// --[option]=...
func ValidateCpy001(copyCommand *instructions.CopyCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(copyCommand),
	}

	// can't use lookahead, lookbehind as it is not supported in Go's regexp.
	regexpWrongNumberOfDashViolation := regexp.MustCompile(`\s+(|-|-{3,})(chmod|chown|from)[ ]{0,1}=`)
	if regexpWrongNumberOfDashViolation.MatchString(copyCommand.String()) {
		result.SetViolated()
		// update location
		indexSlice := regexpWrongNumberOfDashViolation.FindAllStringIndex(copyCommand.String(), 2)
		result.LocationRange.start.charNumber = indexSlice[0][0]
		result.LocationRange.end.charNumber = indexSlice[0][1]
	}

	return result
}
