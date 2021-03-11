package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"

	Parser "github.com/cremindes/whalelint/parser"
	Utils "github.com/cremindes/whalelint/utils"
)

var _ = NewRule("RUN009", "Pass assume yes flag to package manager in order to be headless.", "",
	ValWarning, ValidateRun009)

// nolint: funlen, nestif
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
	for bashCommandIdx, bashCommand := range bashCommandList {
		if len(bashCommand.SubCommand()) == 0 {
			continue
		}

		if pmMap, ok := packageManagerConfirmOptionMap[bashCommand.Bin()]; ok {
			if Utils.SliceContains(pmMap.subcommandSlice, bashCommand.SubCommand()) &&
				!Utils.SliceContains(bashCommand.OptionKeyList(), pmMap.assumeYesSlice) {
				result.SetViolated()

				// since multiple bash commands can have the same string pattern in them,
				// the location need to be further restricted to the current bash command.
				adjustedLocation := make([]parser.Range, 0)

				// temp workaround, till bash parser can work together with raw parser
				if bashCommandIdx > 0 {
					location := LocationRangeFromCommand(runCommand)
					if location.start.lineNumber == location.end.lineNumber {
						// count sum length of bash commands so far
						for _, bc := range bashCommandList[:bashCommandIdx] {
							location.Start().charNumber += len(bc.String())
						}

						adjustedLocation = []parser.Range{LocationRangeToBKRange(location)}
					} else {
						// very naive heuristics, sorry
						location.Start().lineNumber = location.Start().LineNumber() + bashCommandIdx
						adjustedLocation = []parser.Range{LocationRangeToBKRange(location)}
					}
				}

				if len(adjustedLocation) == 0 {
					result.LocationRange = ParseLocationFromRawParser(bashCommand.Bin(), runCommand.Location())
				} else {
					result.LocationRange = ParseLocationFromRawParser(bashCommand.Bin(), adjustedLocation)
				}
			}
		}
	}

	return result
}
