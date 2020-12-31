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
	Name           string
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

	log.Trace("New rule,", all[len(all)-1].Name, "added.")

	return true
}

func Get() []Rule {
	return all
}

var all []Rule
