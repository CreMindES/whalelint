package ruleset

import (
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
	ValidationFunc func(command instructions.Command) bool
}

func NewRule(name string, description string, severity Severity,
	validationFunc func(command instructions.Command) bool) bool {
	log.Trace("NewRule called")

	all = append(all, Rule{
		name,
		description,
		severity,
		validationFunc,
	})

	log.Trace("New rule,", all[len(all)-1].Id, "added.")

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

func Get() []Rule {
	return all
}

var all []Rule
