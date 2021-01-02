package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
	"robpike.io/filter"

	RuleSet "./linter/ruleset"

	"./linter"
)

func main() {
	log.Debug("We have", len(RuleSet.Get()), "ruleset.")

	/* CLI - TODO */
	filePath := os.Args[1]

	/* Parse Dockerfile */
	stageList, metaArgs, err := getDockerfileAst(filePath)
	if err != nil {
		log.Error("Cannot create Dockerfile AST from \"", filePath, "\".", err)
	}
	log.Debug("metaArgs |", metaArgs)

	/* Run Linter */
	var ruleValidationResultArray []RuleSet.RuleValidationResult = linter.Run(stageList)
	violations := filter.Choose(ruleValidationResultArray,
		                        func(x RuleSet.RuleValidationResult) bool { return x.IsViolated() } )

	/* JSON */
	resultJson, err := json.Marshal(violations)

	fmt.Println(string(resultJson))

	//log.WithFields(log.Fields{
	//	"getParameters": string(resultJson),
	//}).Info("Request received!")
}

func getDockerfileAst(filePathString string) (stages []instructions.Stage, metaArgs []instructions.ArgCommand,
	err error) {
	filePath := filepath.Clean(filePathString)

	fileHandle, err := os.Open(filePath)
	if err != nil {
		log.Error("Cannot open Dockerfile \"", filePath, "\".", err)
	}

	dockerfile, err := parser.Parse(fileHandle)
	if err != nil {
		log.Error("Cannot parser Dockerfile \"", filePath, "\"", err)
	}

	stageList, metaArgs, err := instructions.Parse(dockerfile.AST)
	if err != nil {
		log.Error(err)
	}

	return stageList, metaArgs, err
}
