package ruleset

import (
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	Utils "github.com/cremindes/whalelint/utils"
)

var _ = NewRule("EXP001", "Expose a valid UNIX port.", "", ValError, ValidateExp001)

func ValidateExp001(exposeCommand *instructions.ExposeCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(exposeCommand),
	}

	// Port should be in format `([0-9]+)(\/(tcp|udp)){0,1}`. As it's really simple, a primitive parser is used here
	// instead of regexp for performance and readability.
	for _, portStr := range exposeCommand.Ports {
		separatorIndex := strings.IndexRune(portStr, '/')
		if separatorIndex != -1 {
			// port:protocol format
			port := portStr[:separatorIndex]
			protocol := portStr[separatorIndex+1:]

			isPortValid := Utils.IsUnixPortValid(port)
			isProtocolValid := checkProtocolValue(protocol)

			result.SetViolated(!isPortValid || !isProtocolValid)
			result.LocationRange = ParseLocationFromRawParser(portStr, exposeCommand.Location())
		} else {
			// port only format
			isPortValid := Utils.IsUnixPortValid(portStr)
			result.SetViolated(!isPortValid)
			result.LocationRange = ParseLocationFromRawParser(portStr, exposeCommand.Location())
		}
	}

	return result
}

func checkProtocolValue(protocol string) bool {
	return protocol == "tcp" || protocol == "udp"
}
