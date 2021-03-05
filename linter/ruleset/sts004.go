package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"robpike.io/filter"
)

// STS -> Stage Single.
var _ = NewRule("STS004", "There should only be 1 CMD and/or ENTRYPOINT command.", "TODO",
	ValWarning, ValidateSts004)

func ValidateSts004(stage instructions.Stage) RuleValidationResult {
	result := RuleValidationResult{isViolated: false, LocationRange: CopyLocationRange(stage.Location)}

	cmdCommands := filter.Choose(stage.Commands, func(c instructions.Command) bool { return c.Name() == "cmd" })
	entrypointCommands := filter.Choose(stage.Commands, func(c instructions.Command) bool {
		return c.Name() == "entrypoint"
	})

	if commandSlice, ok := cmdCommands.([]instructions.Command); ok {
		if len(commandSlice) > 1 {
			result.SetViolated()
			// message
			result.message += "More than 1 Entrypoint command."
			// location
			if cmdCommand, ok := commandSlice[1].(*instructions.CmdCommand); ok {
				str := cmdCommand.String()[:len("CMD")] // in case the command is lowercase
				result.LocationRange = ParseLocationFromRawParser(str, cmdCommand.Location())
			}
		}
	}

	if commandSlice, ok := entrypointCommands.([]instructions.Command); ok {
		if len(commandSlice) > 1 {
			result.SetViolated()
			// message
			result.message += "More than 1 Entrypoint command."
			// location
			if entrypointCommand, ok := commandSlice[1].(*instructions.EntrypointCommand); ok {
				str := entrypointCommand.String()[:len("ENTRYPOINT")] // in case the command is lowercase
				result.LocationRange = ParseLocationFromRawParser(str, entrypointCommand.Location())
			}
		}
	}

	return result
}
