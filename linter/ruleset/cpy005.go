package ruleset

import (
	"path"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
)

var _ = NewRule("CPY005", "Prefer ADD over COPY for extracting local archives into an image.", "",
	ValWarning, ValidateCpy005)

func ValidateCpy005(copyCommand *instructions.CopyCommand) RuleValidationResult {
	archiveExtensionList := []string{
		".7z", ".gz", ".lz", "lzo", "lzma", ".tar", ".tb2", ".tbz", ".tbz2", ".tgz",
		".tlz", ".tpz", ".txz", ".tZ", "zx", ".Z", ".zip",
	}

	result := RuleValidationResult{
		isViolated:    false,
		LocationRange: LocationRangeFromCommand(copyCommand),
	}

	fileExt := path.Ext(copyCommand.SourcesAndDest.SourcePaths[0])
	for _, archiveExt := range archiveExtensionList {
		if fileExt == archiveExt {
			result.SetViolated()
			result.LocationRange.start.charNumber = 0
			result.LocationRange.end.charNumber = len("COPY")
		}
	}

	return result
}
