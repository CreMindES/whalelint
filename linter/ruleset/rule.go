package ruleset

import (
	"encoding/json"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	log "github.com/sirupsen/logrus"
)

type Severity struct {
	level int
	name  string
}

var (
	Deprecation = Severity{3, "Deprecation"}
	Error       = Severity{0, "Error"}
	Info        = Severity{2, "Info"}
	Warning     = Severity{1, "Warning"}
)

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
		Severity:    rule.severity.name,
	})
}

func Get() []Rule {
	return all
}

var all []Rule