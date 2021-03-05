package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	Parser "github.com/cremindes/whalelint/parser"
	Utils "github.com/cremindes/whalelint/utils"
)

var _ = NewRule("RUN009", "Pass -y|--yes|--assume-yes flag to apt-get in order to be headless.", "",
	ValWarning, ValidateRun009)

func ValidateRun009(runCommand *instructions.RunCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(runCommand),
	}

	packageManagerConfirmOptionMap := map[string][]string{
		"apt":     {"-y", "--yes", "--assume-yes"},
		"apt-get": {"-y", "--yes", "--assume-yes"},
		"dnf":     {"-y", "--assumeyes"},
		"yum":     {"-y", "--assumeyes"},
		"zypper":  {"-y", "--no-confirm", "-n", "--non-interactive"},
	}

	bashCommandList := Parser.ParseBashCommandList(runCommand.CmdLine)
	for _, bashCommand := range bashCommandList {
		if assumeYesSlice, ok := packageManagerConfirmOptionMap[bashCommand.Bin()]; ok {
			if !Utils.SliceContains(bashCommand.OptionKeyList(), assumeYesSlice) {
				result.SetViolated()
				result.LocationRange = ParseLocationFromRawParser(bashCommand.Bin(), runCommand.Location())
			}
		}
	}

	return result
}
