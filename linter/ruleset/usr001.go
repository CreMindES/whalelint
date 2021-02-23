package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

var _ = NewRule("USR001", "Last USER should not be root.", "", ValWarning, ValidateUsr001)

func ValidateUsr001(stageList []instructions.Stage) RuleValidationResult {
	result := RuleValidationResult{isViolated: false, LocationRange: LocationRange{}}

	lastUser := ""
	lastUserLocationRange := LocationRange{}

	for _, stage := range stageList {
		for _, command := range stage.Commands {
			if userCommand, ok := command.(*instructions.UserCommand); ok {
				lastUser = userCommand.User
				lastUserLocationRange = ParseLocationFromRawParser(lastUser, userCommand.Location())
			}
		}
	}

	if lastUser == "root" {
		result.SetViolated()
		result.LocationRange = lastUserLocationRange
	}

	return result
}
