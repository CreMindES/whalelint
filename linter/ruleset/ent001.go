package ruleset

import (
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

var _ = NewRule("ENT001", "Prefer JSON notation array format for CMD and ENTRYPOINT", "", ValWarning,
	ValidateEnt001)

func ValidateEnt001(entrypointCommand *instructions.EntrypointCommand) RuleValidationResult {
	// Get location, which also covers the case of multi line string
	locationRange := UnionOfLocationRanges(
		ParseLocationSliceFromRawParser(entrypointCommand.ShellDependantCmdLine.CmdLine, entrypointCommand.Location()),
	)

	// buildkit's instructions package handleJSONArgs parses CMD, ENTRYPOINT, SHELL and RUN commands.
	// In case of multiple str elements, if they are in JSON notation array format, our command's
	// ShellDependantCmdLine's CmdLine will be longer than 1.
	if len(entrypointCommand.ShellDependantCmdLine.CmdLine) > 1 || entrypointCommand.CmdLine == nil {
		return RuleValidationResult{
			isViolated:    false,
			LocationRange: locationRange,
		}
	}

	// In case there is only one str element, then format needs to be checked
	formatViolation := true
	str := entrypointCommand.String()[len(entrypointCommand.Name()):]
	str = strings.TrimSpace(str)

	if str[0] == '[' && str[len(str)-1] == ']' {
		innerStr := strings.TrimSpace(str[1 : len(str)-1])
		innerStrSlice := strings.Split(innerStr, ",")

		for _, s := range innerStrSlice {
			s = strings.TrimSpace(s)
			if strings.HasPrefix(s, "\"") && strings.HasSuffix(s, "\"") {
				// This point should not be reached, as the first check should be enough in this case.
				formatViolation = false
			}
		}
	}

	return RuleValidationResult{
		isViolated:    formatViolation,
		LocationRange: locationRange,
	}
}
