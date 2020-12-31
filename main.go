package main

import (
	"os"
	"path/filepath"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"

    RuleSet "./linter/ruleset"
)

func main() {
	log.Debug("We have", len(RuleSet.Get()), "ruleset.")

	filePath := "./Dockerfile"

	stageList, metaArgs, err := getDockerfileAst(filePath)
	if err != nil {
		log.Error("Cannot create Dockerfile AST from \"", filePath, "\".", err)
	}

	log.Debug("metaArgs |", metaArgs)

	// OK, now we have the AST of the Dockerfile
	for _, stage := range stageList {
		for _, command := range stage.Commands {
			for _, rule := range RuleSet.Get() {
				if !rule.ValidationFunc(command) {
					log.Debug("Failed on rule", rule.Name)
				}
			}
		}
	}
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
