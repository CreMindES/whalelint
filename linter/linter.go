package linter

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	log "github.com/sirupsen/logrus"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

var MainLinter Linter // nolint:gochecknoglobals

// Linter
// TODO: add config Config.
type Linter struct{}

// nolint:nestif, funlen, gocognit
/* Validate each Dockerfile AST entry against rules in ruleset package. */
func (l *Linter) Run(stageList []instructions.Stage) []RuleSet.RuleValidationResult {
	ruleValidationResultArray := make([]RuleSet.RuleValidationResult, 0)

	if len(stageList) == 0 {
		return ruleValidationResultArray
	}

	// Call Dockerfile AST level validators
	stageListRuleSet := RuleSet.GetRulesForAstElement(stageList)
	for _, rule := range stageListRuleSet {
		validationResult := rule.Validate(stageList)
		ruleValidationResultArray = append(ruleValidationResultArray, validationResult)
	}

	// Get rules for stage elements
	stageRuleSet := RuleSet.GetRulesForAstElement(stageList[0])
	// Go over the stages
	for _, stage := range stageList {
		// Call Dockerfile stage level validators
		for _, rule := range stageRuleSet {
			validationResult := rule.Validate(stage)
			ruleValidationResultArray = append(ruleValidationResultArray, validationResult)
		}

		argMap := make(map[string]string)

		for _, command := range stage.Commands {
			// Call Dockerfile Command level validators, but first filter them by type
			if argCommand, ok := command.(*instructions.ArgCommand); ok {
				for _, rule := range RuleSet.GetRulesForAstElement(argCommand) {
					validationResult := rule.Validate(argCommand)
					ruleValidationResultArray = append(ruleValidationResultArray, validationResult)
				}

				// TODO: better solution than this workaround
				// register stage arg value
				value := *argCommand.Args[0].Value
				argMap[argCommand.Args[0].Key] = value[1 : len(value)-1]
			} else if cmdCommand, ok := command.(*instructions.CmdCommand); ok {
				for _, rule := range RuleSet.GetRulesForAstElement(cmdCommand) {
					validationResult := rule.Validate(cmdCommand)
					ruleValidationResultArray = append(ruleValidationResultArray, validationResult)
				}
			} else if copyCommand, ok := command.(*instructions.CopyCommand); ok {
				for _, rule := range RuleSet.GetRulesForAstElement(copyCommand) {
					validationResult := rule.Validate(copyCommand)
					ruleValidationResultArray = append(ruleValidationResultArray, validationResult)
				}
			} else if entrypointCommand, ok := command.(*instructions.EntrypointCommand); ok {
				for _, rule := range RuleSet.GetRulesForAstElement(entrypointCommand) {
					validationResult := rule.Validate(entrypointCommand)
					ruleValidationResultArray = append(ruleValidationResultArray, validationResult)
				}
			} else if exposeCommand, ok := command.(*instructions.ExposeCommand); ok {
				for _, rule := range RuleSet.GetRulesForAstElement(exposeCommand) {
					ResolveSliceFromArgMap(exposeCommand.Ports, argMap)

					validationResult := rule.Validate(exposeCommand)
					ruleValidationResultArray = append(ruleValidationResultArray, validationResult)
				}
			} else if labelCommand, ok := command.(*instructions.LabelCommand); ok {
				for _, rule := range RuleSet.GetRulesForAstElement(labelCommand) {
					validationResult := rule.Validate(labelCommand)
					ruleValidationResultArray = append(ruleValidationResultArray, validationResult)
				}
			} else if runCommand, ok := command.(*instructions.RunCommand); ok {
				for _, rule := range RuleSet.GetRulesForAstElement(runCommand) {
					validationResult := rule.Validate(runCommand)
					ruleValidationResultArray = append(ruleValidationResultArray, validationResult)
				}
			} else if shellCommand, ok := command.(*instructions.ShellCommand); ok {
				for _, rule := range RuleSet.GetRulesForAstElement(shellCommand) {
					validationResult := rule.Validate(shellCommand)
					ruleValidationResultArray = append(ruleValidationResultArray, validationResult)
				}
			} else if userCommand, ok := command.(*instructions.UserCommand); ok {
				for _, rule := range RuleSet.GetRulesForAstElement(userCommand) {
					validationResult := rule.Validate(userCommand)
					ruleValidationResultArray = append(ruleValidationResultArray, validationResult)
				}
			} else if workdirCommand, ok := command.(*instructions.WorkdirCommand); ok {
				for _, rule := range RuleSet.GetRulesForAstElement(workdirCommand) {
					validationResult := rule.Validate(workdirCommand)
					ruleValidationResultArray = append(ruleValidationResultArray, validationResult)
				}
			} else if maintainerCommand, ok := command.(*instructions.MaintainerCommand); ok {
				for _, rule := range RuleSet.GetRulesForAstElement(maintainerCommand) {
					validationResult := rule.Validate(maintainerCommand)
					ruleValidationResultArray = append(ruleValidationResultArray, validationResult)
				}
			} else {
				log.Error("Unhandled Command!")
			}
		}
	}

	return ruleValidationResultArray
}

func ResolveSliceFromArgMap(strSlice []string, argMap map[string]string) {
	for i, str := range strSlice {
		strSlice[i] = ResolveValueFromArgMap(str, argMap)
	}
}

func ResolveValueFromArgMap(str string, argMap map[string]string) string {
	for key, value := range argMap {
		if "$"+key == str || "${"+key+"}" == str {
			return value
		}
	}

	return str
}
