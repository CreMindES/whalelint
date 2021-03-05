package lsp

import (
	Log "github.com/sirupsen/logrus"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

// VSCodeRangeFromLocationRange converts RuleSet.LocationRange to VS Code's Range Go equivalent type.
func VSCodeRangeFromLocationRange(lr RuleSet.LocationRange) Range {
	start := lr.Start()
	end := lr.End()

	// nolint: gofmt,gofumpt,goimports
	r := Range{
		Start: Position{
			Line:      float64(start.LineNumber()-1),
			Character: float64(start.CharNumber()),
		},
		End:   Position{
			Line:      float64(end.LineNumber()-1),
			Character: float64(end.CharNumber()),
		},
	}

	return r
}

// VSCodeSeverityFromSeverity convert RuleSet.Severity to VS Code's Severity Go equivalent type.
func VSCodeSeverityFromSeverity(s RuleSet.Severity) DiagnosticSeverity {
	switch s {
	case RuleSet.ValError:
		return SeverityError
	case RuleSet.ValWarning:
		return SeverityWarning
	case RuleSet.ValInfo:
		return SeverityInformation
	case RuleSet.ValDeprecation:
		return SeverityHint
	case RuleSet.ValUnknown:
		return SeverityInformation
	default:
		Log.Error("Invalid type")
	}

	return SeverityHint
}
