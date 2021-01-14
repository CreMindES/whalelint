package main

// nolint
import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
	"robpike.io/filter"

	Linter  "github.com/cremindes/whalelint/linter"
	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func main() {
	log.Debug("We have", RuleSet.Get().Count(), "ruleset.")

	/* CLI - TODO */
	filePath := os.Args[1]

	/* Parse Dockerfile */
	stageList, metaArgs := getDockerfileAst(filePath)
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

func getDockerfileAst(filePathString string) (stages []instructions.Stage, metaArgs []instructions.ArgCommand) {
	filePath := filepath.Clean(filePathString)

	fileHandle, err := os.Open(filePath)
	if err != nil {
		log.Error("Cannot open Dockerfile \"", filePath, "\".", err)
	}

	dockerfile, err := parser.Parse(fileHandle)
	if err != nil {
		log.Error("Cannot parse Dockerfile \"", filePath, "\"", err)
	}

	stageList, metaArgs, err := instructions.Parse(dockerfile.AST)
	if err != nil {
		log.Error("Cannot create Dockerfile AST from \"", filePath, "\".", err)
	}

	return stageList, metaArgs
}

func printResultAsJSON(violations interface{}) {
	resultJSON, err := json.Marshal(violations)
	if err != nil {
		log.Error(err)
	}

	fmt.Println(string(resultJSON)) // nolint:forbidigo
}
