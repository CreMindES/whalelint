package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	Parser "github.com/cremindes/whalelint/parser"
	Utils "github.com/cremindes/whalelint/utils"
)

var _ = NewRule("RUN010", "Pass --no-install-recommends to avoid installing unnecessary packages.", "",
	ValWarning, ValidateRun010)

func ValidateRun010(runCommand *instructions.RunCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(runCommand),
	}

	binSlice := []string{"apt-get", "apt"}
	option := "--no-install-recommends"

	bashCommandList := Parser.ParseBashCommandList(runCommand.CmdLine)
	for _, bashCommand := range bashCommandList {
		if Utils.EqualsEither(bashCommand.Bin(), binSlice) && bashCommand.SubCommand() == "install" &&
			!Utils.SliceContains(bashCommand.OptionKeyList(), option) {
			result.SetViolated()
			result.LocationRange = ParseLocationFromRawParser(bashCommand.SubCommand(), runCommand.Location())
		}
	}

	return result
}
