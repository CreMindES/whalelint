package main

// nolint
import (
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"robpike.io/filter"

	Linter "github.com/cremindes/whalelint/linter"
	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
	Utils "github.com/cremindes/whalelint/utils"
)

func main() {
	log.SetLevel(log.DebugLevel)
	log.Debug("We have", RuleSet.Get().Count(), "ruleset.")

	/* CLI - TODO */
	filePath := os.Args[1]

	/* Parse Dockerfile */
	stageList, metaArgs := Utils.GetDockerfileAst(filePath)
	if metaArgs != nil {
		log.Debug("metaArgs |", metaArgs)
	}

	/* Run Linter */
	ruleValidationResultArray := Linter.Run(stageList)
	violations := filter.Choose(ruleValidationResultArray,
		func(x RuleSet.RuleValidationResult) bool {
			return x.IsViolated()
		})

	/* Print result | TODO: cli dependent output */
	printResultAsJSON(violations)
}

func printResultAsJSON(violations interface{}) {
	resultJSON, err := json.Marshal(violations)
	if err != nil {
		log.Error(err)
	}

	fmt.Println(string(resultJSON)) // nolint:forbidigo
}
