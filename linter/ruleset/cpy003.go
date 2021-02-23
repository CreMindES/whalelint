package ruleset

import (
	"regexp"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

var _ = NewRule("CPY003", "COPY chown flag should be in --chown=${USER}:${GROUP} format.", "",
	ValError, ValidateCpy003)

func ValidateCpy003(copyCommand *instructions.CopyCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(copyCommand),
	}

	regexpUserGroupPair := regexp.MustCompile(`^([a-zA-Z0-9]|[a-zA-Z0-9]:[a-zA-Z0-9]){1,}$`)
	result.SetViolated(!regexpUserGroupPair.MatchString(copyCommand.Chown) && len(copyCommand.Chown) != 0)

	if result.IsViolated() {
		result.message = "Invalid user and group pair"
		result.LocationRange = ParseLocationFromRawParser(copyCommand.Chown, copyCommand.Location())
	}

	return result
}
