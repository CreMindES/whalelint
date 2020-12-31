package ruleset

import (
	"path/filepath"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

var _ = NewRule("DL3000", "WORKDIR...", Error, ValidateDl3000)

func ValidateDl3000(command instructions.Command) bool {
	if workdirCommand, ok := command.(*instructions.WorkdirCommand); ok {
		return filepath.IsAbs(workdirCommand.Path)
	}

	return true
}
