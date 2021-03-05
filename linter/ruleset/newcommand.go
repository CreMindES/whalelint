package ruleset

import (
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
