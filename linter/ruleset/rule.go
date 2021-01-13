package ruleset

import (
	"encoding/json"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	log "github.com/sirupsen/logrus"
)

// Severity type represents a severity, with an int level and a String function.
type Severity int

const (
	ValError Severity = iota
	ValDeprecation
	ValInfo
	ValWarning
)

// Severity.String() converts the raw Severity into a string.
func (severity Severity) String() string {
	switch severity {
	case ValDeprecation:
		return "Deprecation"
	case ValError:
		return "Error"
	case ValInfo:
		return "Info"
	case ValWarning:
		return "Warning"
	default:
		return "Unknown"
	}
}

type Rule struct {
	id             string
	description    string
	severity       Severity
	validationFunc func(command instructions.Command) RuleValidationResult
}

func NewRule(name string, description string, severity Severity,
	validationFunc func(command instructions.Command) RuleValidationResult) bool {
	log.Trace("NewRule called")

	all = append(all, Rule{
		name,
		description,
		severity,
		validationFunc,
	})

	log.Trace("New rule,", all[len(all)-1].Id(), "added.")

	return true
}

func (rule *Rule) Id() string {
	return rule.id
}

func (rule *Rule) Severity() Severity {
	return rule.severity
}

func (rule *Rule) Description() string {
	return rule.description
}

func (rule *Rule) ValidationFunc(command instructions.Command) RuleValidationResult {
	return rule.validationFunc(command)
}

func (rule Rule) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Id           string
		Description  string
		Severity     string
	} {
		Id:          rule.id,
		Description: rule.description,
		Severity:    rule.severity.String(),
	})
}

func Get() []Rule {
	return all
}

var all []Rule