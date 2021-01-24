package ruleset

import (
	"regexp"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

var _ = NewRule("CPY002", "COPY --chmod=XXXX where XXXX should be a valid permission set value.",
	ValError, ValidateCpy002)

// checks COPY --chmod option format for obvious errors
// --chmod=XXXX, where XXXX is a valid permission set value.
func ValidateCpy002(copyCommand *instructions.CopyCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(copyCommand),
	}

	locationOffset := 5

	regexpUnixPermission := regexp.MustCompile("^[0-7]{1,4}$")
	result.SetViolated(!regexpUnixPermission.MatchString(copyCommand.Chmod) && len(copyCommand.Chmod) != 0)
	result.LocationRange.start.charNumber = locationOffset

	if result.IsViolated() {
		result.message = "Invalid Unix permission value."
		unixPermissionValueIndex := strings.Index(copyCommand.String(), copyCommand.Chmod)
		result.LocationRange.start.charNumber = unixPermissionValueIndex
		result.LocationRange.end.charNumber = unixPermissionValueIndex + len(copyCommand.Chmod)
	}

	return result
}
