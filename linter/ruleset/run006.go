package ruleset

import (
	"regexp"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	"github.com/cremindes/whalelint/parser"
)

// TODO: revisit

var _ = NewRule("RUN006", "Clean cache after package manager operation.", "", ValWarning,
	ValidateRun006)

func ValidateRun006(runCommand *instructions.RunCommand) RuleValidationResult {
	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(runCommand),
	}

	if len(runCommand.CmdLine) == 0 {
		return result
	} else if len(runCommand.CmdLine[0]) == 0 {
		return result
	}

	// TODO: move to bash parsing/bash utils
	packageManagerInstallCleanMap := map[string]*regexp.Regexp{
		"apt":     regexp.MustCompile(`(apt clean|rm -rf /var/lib/apt/lists)`),
		"apt-get": regexp.MustCompile(`(apt-get clean|rm -rf /var/lib/apt/lists)`),
		"yum":     regexp.MustCompile(`yum clean all`),
		"apk":     regexp.MustCompile(`apk(.*)--no-cache`),
		"pip":     regexp.MustCompile(`pip(.*)--no-cache-dir`),
		"zypper":  regexp.MustCompile(`zypper (clean|-a)`),
		"dnf":     regexp.MustCompile(`dnf clean all`),
	}

	bashCommandList := parser.ParseBashCommandChain(runCommand).BashCommandList

	for i, bashCommand := range bashCommandList {
		packageManager := bashCommand.Bin()
		if pmRegexp, ok := packageManagerInstallCleanMap[packageManager]; ok {
			if parser.HasPackageUpdateCommand(packageManager, bashCommand) {
				for j := i; j < len(bashCommandList); j++ {
					if pmRegexp.MatchString(bashCommandList[j].String()) {
						return RuleValidationResult{
							isViolated:    false,
							LocationRange: LocationRangeFromCommand(runCommand),
						}
					}
				}

				return RuleValidationResult{
					isViolated:    true,
					LocationRange: ParseLocationFromRawParser(packageManager, runCommand.Location()),
				}
			}
		}
	}

	return result
}
