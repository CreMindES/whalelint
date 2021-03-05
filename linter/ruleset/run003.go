package ruleset

import (
	"regexp"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

var _ = NewRule("RUN003", "Operators \"&&, ||, |\" has no affect after semicolon.", "", ValError,
	ValidateRun003)

func ValidateRun003(runCommand *instructions.RunCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(runCommand),
	}

	regexpInvalidPattern := regexp.MustCompile(`;\s+[|&]{1,2}`)
	// TODO: consider using FindSubmatchIndex in order to support multiple locations
	if match := regexpInvalidPattern.FindString(runCommand.String()); len(match) > 0 {
		result.SetViolated()
		result.LocationRange = ParseLocationFromRawParser(match, runCommand.Location())
		result.message = "Probably not what you wanted: " + match
	}

	return result
}
