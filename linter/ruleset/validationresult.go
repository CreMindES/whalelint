package ruleset

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type RuleValidationResult struct {
	rule          Rule
	isViolated    bool
	message       string
	LocationRange LocationRange
}

func (ruleValidationResult *RuleValidationResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule          Rule          `json:"Rule"`
		IsViolated    bool          `json:"IsViolated"`
		Message       string        `json:"Message"`
		LocationRange LocationRange `json:"LocationRange"`
	}{
		Rule:          ruleValidationResult.rule,
		IsViolated:    ruleValidationResult.isViolated,
		Message:       ruleValidationResult.message,
		LocationRange: ruleValidationResult.LocationRange,
	})
}

func (ruleValidationResult RuleValidationResult) IsViolated() bool {
	return ruleValidationResult.isViolated
}

func (ruleValidationResult *RuleValidationResult) SetViolated(params ...bool) {
	switch len(params) {
	case 0:
		ruleValidationResult.isViolated = true
	case 1:
		ruleValidationResult.isViolated = params[0]
	default:
		log.Error("Invalid params to RuleValidationResult::SetViolated")
	}
}

func (ruleValidationResult *RuleValidationResult) Location() LocationRange {
	return ruleValidationResult.LocationRange
}

func (ruleValidationResult *RuleValidationResult) SetLocation(startLineNumber, startCharNumber,
	endLineNumber, endCharNumber int) {
	ruleValidationResult.LocationRange.start.lineNumber = startLineNumber
	ruleValidationResult.LocationRange.start.charNumber = startCharNumber
	ruleValidationResult.LocationRange.end.lineNumber = endLineNumber
	ruleValidationResult.LocationRange.end.charNumber = endCharNumber
}

func (ruleValidationResult *RuleValidationResult) SetLocationRangeFrom(locationRange LocationRange) {
	ruleValidationResult.LocationRange = locationRange
}

func (ruleValidationResult *RuleValidationResult) SetRule(rule Rule) {
	ruleValidationResult.rule = rule
}
