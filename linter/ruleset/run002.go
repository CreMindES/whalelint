package ruleset

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dustin/go-humanize/english"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	Parser "github.com/cremindes/whalelint/parser"
	Utils "github.com/cremindes/whalelint/utils"
)

var _ = NewRule("RUN002", "Consider pinning versions of packages", "", ValWarning, ValidateRun002)

// nolint:funlen
func ValidateRun002(runCommand *instructions.RunCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(runCommand),
	}

	printBodyWidthThreshold := 60

	packageWithoutVersionList := make([]string, 0)
	filterFunc := func(packageVersion string) bool {
		return len(packageVersion) == 0
	}
	filterFuncNpm := func(packageStr string) bool {
		return strings.ContainsRune(packageStr, '@')
	}
	filterFuncRubyGem := func(packageStr string) bool {
		return strings.ContainsRune(packageStr, ':')
	}
	filterFuncYum := func(packageStr string) bool {
		return strings.ContainsRune(packageStr, '-')
	}

	// For now all package installs are validated here. As soon as this becomes too long, complex or hard to read,
	// this will be divided up into separate rules.
	bashCommandList := Parser.ParseBashCommandList(runCommand.CmdLine)
	// nolint:wsl
	for _, bashCommand := range bashCommandList {
		if Parser.IsDebPackageInstall(bashCommand) {
			packageWithoutVersionList = Utils.FilterMapByValue(bashCommand.ArgMap(), filterFunc)
			result.message = "Debian"
		}
		if Parser.IsRpmPackageInstall(bashCommand) {
			packageWithoutVersionList = Utils.FilterMapByValue(bashCommand.ArgMap(), filterFuncYum)
			result.message = "RPM"
		}
		if Parser.IsApkPackageInstall(bashCommand) {
			packageWithoutVersionList = Utils.FilterMapByValue(bashCommand.ArgMap(), filterFunc)
			result.message = "Alpine"
		}
		if Parser.IsSusePackageInstall(bashCommand) {
			packageWithoutVersionList = Utils.FilterMapByValue(bashCommand.ArgMap(), filterFunc)
			result.message = "Suse"
		}
		if Parser.IsFedoraPackageInstall(bashCommand) {
			packageWithoutVersionList = Utils.FilterMapByValue(bashCommand.ArgMap(), filterFunc)
			result.message = "Fedora/RMP"
		}
		if Parser.IsPythonPackageInstall(bashCommand) {
			// case pip install -r [requirement.txt]
			if Utils.SliceContains(bashCommand.OptionKeyList(), []string{"-r", "--requirement"}) {
				continue
			}
			packageWithoutVersionList = Utils.FilterMapByValue(bashCommand.ArgMap(), filterFunc)
			result.message = "Python"
		}
		if Parser.IsNpmPackageInstall(bashCommand) {
			packageWithoutVersionList = Utils.FilterMapKeys(bashCommand.ArgMap(), filterFuncNpm)
			result.message = "NPM"
		}
		if Parser.IsRubyPackageInstall(bashCommand) {
			// TODO: check for -v flag
			packageWithoutVersionList = Utils.FilterMapKeys(bashCommand.ArgMap(), filterFuncRubyGem)
			result.message = "Ruby"
		}
	}

	result.SetViolated(len(packageWithoutVersionList) > 0)

	// Update location
	if result.isViolated && Parser.RawParser.IsInitialized() {
		packageLocationRangeSlice := ParseLocationSliceFromRawParser(packageWithoutVersionList, runCommand.Location())
		result.LocationRange = UnionOfLocationRanges(packageLocationRangeSlice)
	}

	// Assemble rule violation message
	if result.isViolated {
		sort.Strings(packageWithoutVersionList)
		packageWithoutVersionListStr := strings.Join(packageWithoutVersionList, ", ")

		if len(packageWithoutVersionListStr) < printBodyWidthThreshold {
			result.message = fmt.Sprintf("%s \"%s\" %s no version specified.",
				english.PluralWord(len(packageWithoutVersionList), "Package", ""),
				packageWithoutVersionListStr,
				english.PluralWord(len(packageWithoutVersionList), "has", "have"),
			)
		} else {
			result.message = fmt.Sprintf("%d %s %s no version specified.",
				len(packageWithoutVersionList),
				english.PluralWord(len(packageWithoutVersionList), "package", ""),
				english.PluralWord(len(packageWithoutVersionList), "has", "have"),
			)
		}
	}

	return result
}
