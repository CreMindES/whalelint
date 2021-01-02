package ruleset

import (
	"path/filepath"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	log "github.com/sirupsen/logrus"
)

var _ = NewRule("DL3000", "WORKDIR should be an absolute path for clarity and reliability.", Warning,
	            ValidateDl3000)

func ValidateDl3000(command instructions.Command) RuleValidationResult {
	result := RuleValidationResult{
		isViolated: false,
	}

	if workdirCommand, ok := command.(*instructions.WorkdirCommand); ok {
		if filepath.IsAbs(workdirCommand.Path) == false {
			result.SetViolated()

			var lineString = workdirCommand.String()

			result.SetLocation(
				workdirCommand.Location()[0].Start.Line,
				workdirCommand.Location()[0].Start.Character,
				workdirCommand.Location()[0].End.Line,
				workdirCommand.Location()[0].End.Character + len(lineString))
		}
	}

	log.Trace("ValidateDl3000 result:", result)

	return result
}
