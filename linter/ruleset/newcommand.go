package ruleset

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
)

func NewRunCommand(cmd string, locationRange LocationRange) *instructions.RunCommand {
	start, end := locationRange.Start(), locationRange.End()

	node := &parser.Node{ // nolint:exhaustivestruct
		Value:     "run",
		Original:  "RUN " + cmd,
		StartLine: start.LineNumber(),
		EndLine:   end.LineNumber(),
	}

	node.Next = &parser.Node{Value: cmd} // nolint:exhaustivestruct

	instruction, err := instructions.ParseInstruction(node)
	if err != nil {
		log.Error("cannot parse mock CopyCommand AST Node")
	}

	command, ok := instruction.(*instructions.RunCommand)
	if !ok {
		log.Error("cannot type assert instruction to *instructions.CopyCommand")
	}

	return command
}

var errInvalidValue = errors.New("invalid value")

func NewEntrypointCommand(str string, lineNumber int) (*instructions.EntrypointCommand, error) {
	if lineNumber < 1 {
		return nil, errInvalidValue
	}

	buf := bytes.Buffer{}

	buf.WriteString("FROM golang:1.16\n")
	for i := 1; i < lineNumber-1; i++ { // nolint:wsl
		buf.WriteString("# Padding ...\n")
	}
	buf.WriteString("ENTRYPOINT " + str + "\n")

	reader := bytes.NewReader(buf.Bytes())

	stageList, err := parseMockDockerfile(reader)
	if err != nil {
		return nil, err
	}

	command := stageList[0].Commands[0]

	entrypointCommand, ok := command.(*instructions.EntrypointCommand)
	if !ok {
		return nil, fmt.Errorf("RuleSet test helper | %w", err)
	}

	return entrypointCommand, nil
}

func NewCmdCommand(str string, lineNumber int) (*instructions.CmdCommand, error) {
	if lineNumber < 1 {
		return nil, errInvalidValue
	}

	buf := bytes.Buffer{}

	buf.WriteString("FROM golang:1.16\n")
	for i := 1; i < lineNumber-1; i++ { // nolint:wsl
		buf.WriteString("# Padding ...\n")
	}
	buf.WriteString("CMD " + str + "\n")

	reader := bytes.NewReader(buf.Bytes())

	stageList, err := parseMockDockerfile(reader)
	if err != nil {
		return nil, err
	}

	command := stageList[0].Commands[0]

	cmdCommand, ok := command.(*instructions.CmdCommand)
	if !ok {
		return nil, fmt.Errorf("RuleSet test helper | %w", err)
	}

	return cmdCommand, nil
}

func NewMaintainerCommand(str string) (*instructions.MaintainerCommand, error) {
	buf := bytes.Buffer{}
	buf.WriteString("FROM golang:1.16\n")
	buf.WriteString("MAINTAINER " + str + "\n")

	reader := bytes.NewReader(buf.Bytes())

	stageList, err := parseMockDockerfile(reader)
	if err != nil {
		return nil, err
	}

	command := stageList[0].Commands[0]

	maintainerCommand, ok := command.(*instructions.MaintainerCommand)
	if !ok {
		return nil, fmt.Errorf("RuleSet test helper | %w", err)
	}

	return maintainerCommand, nil
}

func parseMockDockerfile(reader io.Reader) ([]instructions.Stage, error) {
	dockerfile, err := parser.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("dockerfile parse | %w", err)
	}

	stageList, _, err := instructions.Parse(dockerfile.AST)
	if err != nil {
		return nil, fmt.Errorf("dockerfile stage parse | %w", err)
	}

	return stageList, nil
}
