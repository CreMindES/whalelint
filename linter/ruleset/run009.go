package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	Parser "github.com/cremindes/whalelint/parser"
	Utils "github.com/cremindes/whalelint/utils"
)

var _ = NewRule("RUN009", "Pass assume yes flag to package manager in order to be headless.", "",
	ValWarning, ValidateRun009)

func ValidateRun009(runCommand *instructions.RunCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(runCommand),
	}

	packageManagerConfirmOptionMap := map[string]struct {
		subcommandSlice []string
		assumeYesSlice  []string
	}{
		"apt":     {[]string{"install", "remove", "purge"}, []string{"-y", "--yes", "--assume-yes"}},
		"apt-get": {[]string{"install", "remove", "purge"}, []string{"-y", "--yes", "--assume-yes"}},
		"dnf":     {[]string{"install", "remove", "downgrade"}, []string{"-y", "--assumeyes"}},
		"yum":     {[]string{"install", "remove", "downgrade"}, []string{"-y", "--assumeyes"}},
		"zypper": {
			[]string{"install", "in", "remove", "rm"},
			[]string{"-y", "--no-confirm", "-n", "--non-interactive"},
		},
	}

	bashCommandList := Parser.ParseBashCommandList(runCommand.CmdLine)
	for _, bashCommand := range bashCommandList {
		if len(bashCommand.SubCommand()) == 0 {
			continue
		}

		if pmMap, ok := packageManagerConfirmOptionMap[bashCommand.Bin()]; ok {
			if Utils.SliceContains(pmMap.subcommandSlice, bashCommand.SubCommand()) &&
				!Utils.SliceContains(bashCommand.OptionKeyList(), pmMap.assumeYesSlice) {
				result.SetViolated()
				result.LocationRange = ParseLocationFromRawParser(bashCommand.Bin(), runCommand.Location())
			}
		}
	}

	return result
}
