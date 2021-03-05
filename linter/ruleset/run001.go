package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	log "github.com/sirupsen/logrus"

	Parser "github.com/cremindes/whalelint/parser"
)

var _ = NewRule("RUN001", "Some bash commands make no sense in an ordinary Docker container.", "", ValWarning,
	ValidateRun001)

func ValidateRun001(runCommand *instructions.RunCommand) RuleValidationResult {
	invalidCmdSet := []string{"free", "kill", "mount", "ps", "reboot", "service", "shutdown", "top"}

	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(runCommand),
	}
	result.LocationRange.end.charNumber += len(runCommand.ShellDependantCmdLine.CmdLine)

	bashCommandChain := Parser.ParseBashCommandChain(runCommand.CmdLine)

	for _, bashCommand := range bashCommandChain.BashCommandList {
		for _, invalidCmd := range invalidCmdSet {
			if bashCommand.Bin() == invalidCmd {
				result.SetViolated()
			}
		}
	}

	log.Trace("ValidateRun001 result:", result)

	return result
}
