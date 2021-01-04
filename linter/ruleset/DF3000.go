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
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(command),
	}

	if workdirCommand, ok := command.(*instructions.WorkdirCommand); ok {
		if filepath.IsAbs(workdirCommand.Path) == false {
			result.SetViolated()
			result.LocationRange.end.charNumber += len(workdirCommand.String())
		}
	}

	log.Trace("ValidateDl3000 result:", result)

	return result
}
