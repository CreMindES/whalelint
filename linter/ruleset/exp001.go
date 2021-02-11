package ruleset

import (
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	Utils "github.com/cremindes/whalelint/utils"
)

var _ = NewRule("EXP001", "Expose a valid UNIX port.", "", ValWarning, ValidateExp001)

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
		} else {
			// port only format
			isPortValid := Utils.IsUnixPortValid(portStr)
			result.SetViolated(!isPortValid)
		}
	}

	// location
	result.LocationRange.start.charNumber = len("EXPOSE ")
	result.LocationRange.end.charNumber = len(exposeCommand.String())

	return result
}

func checkProtocolValue(protocol string) bool {
	return protocol == "tcp" || protocol == "udp"
}
