package ruleset

import (
	"regexp"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

var _ = NewRule(
	"CPY001", "Flag format validation | COPY --[chmod|chown|from]=... srcList... dest|destDir",
	``+"`COPY`"+` command
- flags [`+"`chmod`"+`|`+"`chown`"+`|`+"`from`"+`] are preceded by two dashes.
- `+"`chmod`"+` should have a valid Linux permission value.
- `+"`chown`"+` should be in `+"`user:group`"+` format.`,
	ValError, ValidateCpy001)

// checks COPY options format for obvious errors
// --[option]=...
func ValidateCpy001(copyCommand *instructions.CopyCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(copyCommand),
	}

	// can't use lookahead, lookbehind as it is not supported in Go's regexp.
	regexpWrongNumberOfDashViolation := regexp.MustCompile(`[^-](|-|-{3,})(chmod|chown|from)[ ]{0,1}=`)
	if regexpWrongNumberOfDashViolation.MatchString(copyCommand.String()) {
		result.SetViolated()
		result.message = "Flags must be prefixed with exactly two dashes."

		wrongFlagStr := regexpWrongNumberOfDashViolation.FindString(copyCommand.SourcesAndDest[0])
		result.LocationRange = ParseLocationFromRawParser(wrongFlagStr, copyCommand.Location())
	}

	// TODO: support invalid flag. Note: it might need contribution to buildkit.

	return result
}
