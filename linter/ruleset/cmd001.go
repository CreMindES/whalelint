package ruleset

import (
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

var _ = NewRule("CMD001", "Prefer JSON notation array format for CMD and ENTRYPOINT", "", ValWarning,
	ValidateCmd001)

func ValidateCmd001(cmdCommand *instructions.CmdCommand) RuleValidationResult {
	argStr := cmdCommand.String()[len(cmdCommand.Name()):]
	argStr = strings.TrimSpace(argStr)
	lineNum := cmdCommand.Location()[0].Start.Line

	entrypointCommand, err := NewEntrypointCommand(argStr, lineNum)
	if err != nil {
		return RuleValidationResult{
			isViolated:    true,
			message:       "",
			LocationRange: LocationRangeFromCommand(cmdCommand),
		}
	}

	return ValidateEnt001(entrypointCommand)
}
