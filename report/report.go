// Package report provides different formats, i.e views on rule validation results that are printed to the provided
// io.Writer.
package report

import (
	"encoding/json"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"
	"robpike.io/filter"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

// PrintResultAsJSON prints lint rule violations to writer in JSON format.
func PrintResultAsJSON(ruleValidationResultArray []RuleSet.RuleValidationResult, writer io.Writer) {
	// Filter lint rule violations
	violations := filter.Choose(ruleValidationResultArray,
		func(x RuleSet.RuleValidationResult) bool {
			return x.IsViolated()
		},
	)

	// Marshal result to JSON
	resultJSON, err := json.Marshal(violations)
	if err != nil {
		log.Error(err)
	}

	// Send result to io.Writer
	printToOutput(string(resultJSON), writer)
}

// VerbosityLevel represent the verbosity level used in the summary below.
// For details please see the Verbosity* const definitions.
type VerbosityLevel int

const (
	VerbosityShort  VerbosityLevel = iota // just a quick, one line summary
	VerbosityNormal                       // default verbosity level: position, rule ID and rule definition
	VerbosityHigh                         // normal verbosity extended with further information
)

// SummaryOption holds the global settings for the summary.
type SummaryOption struct {
	NoColor   bool
	Verbosity VerbosityLevel
}

// PrintOptions represents finer grain, printing specific options.
type PrintOptions struct {
	LineNumStrWidthTarget int
	PrintOptionsForSeverityMap
}

// PrintOptionsForSeverityMap is a mapping between severity levels, their serialized string and color function.
type PrintOptionsForSeverityMap map[RuleSet.Severity]struct {
	name    string
	colorFn func(...interface{}) string
}

// FindingsMap shorthand for map[RuleSet.Severity][]RuleSet.RuleValidationResult.
type FindingsMap map[RuleSet.Severity][]RuleSet.RuleValidationResult

// GroupFindings groups lint rule violations based on their Severity.
func GroupFindings(findings []RuleSet.RuleValidationResult) (FindingsMap, bool) {
	// RuleSet::Rule severity list
	severityLevelSlice := RuleSet.GetSeverityList()
	findingsMap := make(map[RuleSet.Severity][]RuleSet.RuleValidationResult, len(severityLevelSlice))
	hasViolations := false

	// group findings be severity
	for i := 0; i < len(severityLevelSlice); i++ {
		sev := severityLevelSlice[i]

		// filter findings for current severity
		sevFilterResult := filter.Choose(findings, func(x RuleSet.RuleValidationResult) bool {
			return x.IsViolated() && x.Severity() == sev
		})
		// assert back to []RuleSet.RuleValidationResult
		sevSlice, ok := sevFilterResult.([]RuleSet.RuleValidationResult)
		if !ok {
			return map[RuleSet.Severity][]RuleSet.RuleValidationResult{}, false
		}
		// sort findings by line number
		sort.Slice(sevSlice, func(i, j int) bool {
			return sevSlice[i].LocationRange.Start().LineNumber() < sevSlice[j].LocationRange.Start().LineNumber()
		})

		findingsMap[severityLevelSlice[i]] = sevSlice
		if len(sevSlice) > 0 {
			hasViolations = true
		}
	}

	return findingsMap, hasViolations
}

// AssembleSummaryHeader prepares the one line short summary header.
// Verbosity::short.
func AssembleSummaryHeader(findingsMap FindingsMap, hasViolation bool, printOptions PrintOptions,
	strBuilder *strings.Builder) {
	// Header | Start
	strBuilder.WriteString("WhaleLint summary: ")

	if !hasViolation {
		strBuilder.WriteString(color.New(color.FgGreen).SprintFunc()("Everything looks good."))

		return
	}

	// Header | Body
	hasPrev := false

	severityLevelSlice := RuleSet.GetSeverityList()
	for i := 0; i < len(severityLevelSlice); i++ {
		severity := severityLevelSlice[i]
		sevPrintOption := printOptions.PrintOptionsForSeverityMap[severity]

		// findings at severity level
		itemList := findingsMap[severity]
		if len(itemList) > 0 {
			printConditionally(", ", hasPrev, strBuilder)

			strBuilder.WriteString(sevPrintOption.colorFn(strconv.Itoa(len(itemList)) + " " + sevPrintOption.name))
			printConditionally(sevPrintOption.colorFn("s"), len(findingsMap[severity]) > 1, strBuilder)

			hasPrev = true
		}
	}

	strBuilder.WriteRune('\n')
}

// AssembleSummaryBody prepares the extension to Verbosity::short, be listing each lint rule violation.
// Format is 'LineNum | RuleID | RuleDefinition'/
// Verbosity::normal = Verbosity::short + this summary body.
func AssembleSummaryBody(findingsMap FindingsMap, printOptions PrintOptions, strBuilder *strings.Builder) {
	strBuilder.WriteRune('\n')

	severityLevelSlice := RuleSet.GetSeverityList()
	for i := 0; i < len(severityLevelSlice); i++ {
		severity := severityLevelSlice[i]
		itemList := findingsMap[severity]
		sevPrintOption := printOptions.PrintOptionsForSeverityMap[severity]

		if len(itemList) == 0 {
			continue
		}

		// print severity group name
		strBuilder.WriteString(sevPrintOption.colorFn(sevPrintOption.name))
		printConditionally(sevPrintOption.colorFn("s"), len(findingsMap[severity]) > 1, strBuilder)
		strBuilder.WriteString(":\n")

		// print individual violations in the following format:
		// Line nnn | RULE ID | RuleValidation.Message
		for _, violation := range itemList {
			lineNumber := strconv.Itoa(violation.Location().Start().LineNumber())
			lineNumber = printWithPadding(lineNumber, printOptions.LineNumStrWidthTarget, padBefore)
			strBuilder.WriteString("Line " + lineNumber + " | ")
			strBuilder.WriteString(sevPrintOption.colorFn(violation.RuleID()) + " | ")
			strBuilder.WriteString(violation.Message())
			strBuilder.WriteRune('\n')
		}

		strBuilder.WriteRune('\n')
	}
}

// PrintSummary prints the RuleValidationResult list's summary using the provided writer and options.
func PrintSummary(violations []RuleSet.RuleValidationResult, writer io.Writer, options SummaryOption) {
	// global color output option
	color.NoColor = options.NoColor

	findingsMap, hasViolation := GroupFindings(violations)

	// Helper print option map
	printOptionsForSeverityMap := PrintOptionsForSeverityMap{
		RuleSet.ValError:       {"Error", color.New(color.FgRed).SprintFunc()},
		RuleSet.ValWarning:     {"Warning", color.New(color.FgYellow).SprintFunc()},
		RuleSet.ValInfo:        {"Info", color.New(color.FgBlue).SprintFunc()},
		RuleSet.ValDeprecation: {"Deprecation", color.New(color.FgCyan).SprintFunc()},
		RuleSet.ValUnknown:     {"Unknown", color.New(color.FgWhite).SprintFunc()},
	}

	printOptions := PrintOptions{
		LineNumStrWidthTarget:      getMaxLine(violations),
		PrintOptionsForSeverityMap: printOptionsForSeverityMap,
	}

	// Main string builder
	strBuilder := &strings.Builder{}

	AssembleSummaryHeader(findingsMap, hasViolation, printOptions, strBuilder)

	// End of VerbosityShort summary
	if options.Verbosity == VerbosityShort { // Short Summary ends here
		printToOutput(strBuilder.String(), writer)

		return
	}

	// Detailed summary
	AssembleSummaryBody(findingsMap, printOptions, strBuilder)

	// Send to io.Writer
	printToOutput(strBuilder.String(), writer)
}

// getMaxLine returns the highest line locations of a given RuleValidationResult set.
func getMaxLine(findingList []RuleSet.RuleValidationResult) int {
	max := 0

	for _, finding := range findingList {
		if finding.Location().Start() != nil {
			findingLine := strconv.Itoa(finding.Location().Start().LineNumber())
			if max < len(findingLine) {
				max = len(findingLine)
			}
		}
	}

	return max
}

// Padding type represents the padding strategy.
type Padding int

const (
	padBefore Padding = iota // pad before string
	padAfter                 // pad after string
)

// printWithPadding prints with optional padding to reach a minimum target string length.
// example: print "hello" string with padding 6 and padBefore will result in "  hello".
func printWithPadding(str string, num int, padding Padding) string {
	strBuilder := strings.Builder{}

	if len(str) < num {
		// need padding
		if padding == padBefore {
			for i := 0; i < num-len(str); i++ {
				strBuilder.WriteRune(' ')
			}
		}

		strBuilder.WriteString(str)

		if padding == padAfter {
			for i := 0; i < num-len(str); i++ {
				strBuilder.WriteRune(' ')
			}
		}
	} else {
		return str
	}

	return strBuilder.String()
}

// printToOutput simple printer wrapper to reduce error handling boilerplate code.
func printToOutput(str string, writer io.Writer) {
	_, err := writer.Write([]byte(str))
	if err != nil {
		log.Error(err)
	}
}

// printConditionally simple printer wrapper to reduce boilerplate code.
func printConditionally(str string, doPrint bool, writer io.Writer) {
	if doPrint {
		_, err := writer.Write([]byte(str))
		if err != nil {
			log.Error(err)
		}
	}
}
