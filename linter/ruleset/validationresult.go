package ruleset

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type RuleValidationResult struct {
	rule          *Rule
	isViolated    bool
	message       string
	LocationRange LocationRange
}

const FORCE = true // used for overriding the latched isViolated flag in SetViolated.

func (ruleValidationResult *RuleValidationResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Rule          *Rule         `json:"Rule"`
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

func (ruleValidationResult *RuleValidationResult) UnmarshalJSON(data []byte) error {
	rvr := struct {
		Rule          *Rule
		IsViolated    bool
		Message       string
		LocationRange LocationRange
	}{}

	err := json.Unmarshal(data, &rvr)
	if err != nil {
		return fmt.Errorf("failed to unmarshal RuleValidationResult: %w", err)
	}

	ruleValidationResult.rule = rvr.Rule
	ruleValidationResult.isViolated = rvr.IsViolated
	ruleValidationResult.message = rvr.Message
	ruleValidationResult.LocationRange = rvr.LocationRange

	return nil
}

func (ruleValidationResult RuleValidationResult) IsViolated() bool {
	return ruleValidationResult.isViolated
}

func (ruleValidationResult *RuleValidationResult) SetViolated(params ...bool) {
	// nolint:gomnd
	switch len(params) {
	case 0:
		ruleValidationResult.isViolated = true
	case 1:
		// isViolated is latched, i.e once set, it cannot be unset, unless passed a second param: FORCE
		if !ruleValidationResult.isViolated {
			ruleValidationResult.isViolated = params[0]
		}
	case 2:
		if params[1] == FORCE && ruleValidationResult.isViolated { // nolint:gosimple
			ruleValidationResult.isViolated = params[0]
		}
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

func (ruleValidationResult *RuleValidationResult) SetRule(rule *Rule) {
	ruleValidationResult.rule = rule
}
