package linter

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	log "github.com/sirupsen/logrus"

	RuleSet "./ruleset"
)

/* Validate each Dockerfile AST entry against rules in ruleset package. */
func Run(stageList []instructions.Stage) []RuleSet.RuleValidationResult {
	var ruleValidationResultArray []RuleSet.RuleValidationResult

	// OK, now we have the AST of the Dockerfile
	for _, stage := range stageList {
		for _, command := range stage.Commands {
			for _, rule := range RuleSet.Get() {
				ruleValidationResult := rule.ValidationFunc(command)

				if ruleValidationResult.IsViolated() {
					log.Debug("Failed on rule", rule.Id(), "[", rule.Severity(), "]",
						"on line", ruleValidationResult.Location())
				}

				ruleValidationResult.SetRule(rule)
				ruleValidationResultArray = append(ruleValidationResultArray, ruleValidationResult)
			}
		}
	}

	return ruleValidationResultArray
}