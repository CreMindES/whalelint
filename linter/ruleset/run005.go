package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	Parser "github.com/cremindes/whalelint/parser"
)

var _ = NewRule("RUN005", "Do not upgrade or dist-upgrade the base image", "", ValError, ValidateRun005)

func ValidateRun005(runCommand *instructions.RunCommand) RuleValidationResult {
	notAdvisedPackageManagerCommandMap := map[string][]string{
		"apt":     {"upgrade", "dist-upgrade"},
		"apt-get": {"upgrade", "dist-upgrade"},
		"apk":     {"upgrade"},
		"dnf":     {"upgrade", "up", "distro-sync", "downgrade"},
		"yum":     {"upgrade", "distro-sync", "dsync", "downgrade"},
		"zypper":  {"update", "up", "dist-upgrade"},
	}

	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(runCommand),
	}

	bashCommandList := Parser.ParseBashCommandList(runCommand.CmdLine)
	for _, bashCommand := range bashCommandList {
		for packageManager, notAdvisedCommandSlice := range notAdvisedPackageManagerCommandMap {
			for _, notAdvisedCommand := range notAdvisedCommandSlice {
				if bashCommand.Bin() == packageManager && bashCommand.SubCommand() == notAdvisedCommand {
					result.SetViolated()
					result.LocationRange = ParseLocationFromRawParser(bashCommand.SubCommand(), runCommand.Location())
				}
			}
		}
	}

	return result
}
