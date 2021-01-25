package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

type mockCopyNode struct {
	original string
	flags    []string
	next     []string
}

// nolint:funlen
func TestValidateCpy001(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		isViolation bool
		name        string
		copyNode    mockCopyNode
	}{
		{isViolation: false, name: "Proper COPY command with 1 --chmod flag.", copyNode: mockCopyNode{
			original: "COPY --chmod=7780 src src2 dst/",
			flags: []string{"--chmod=7780"},
			next: []string{"src", "src2", "dst/"},
		}},
		{isViolation: true, name: "COPY command with 1 -chmod flag.", copyNode: mockCopyNode{
			original: "COPY -chmod=7780 src dst/",
			flags: []string{},
			next: []string{"-chmod=7780", "src", "dst/"},
		}},
		{isViolation: true, name: "COPY command with 1 chmod flag.", copyNode: mockCopyNode{
			original: "COPY chmod=7780 src dst/",
			flags: []string{},
			next: []string{"chmod=7780", "src", "dst/"},
		}},
		{isViolation: true, name: "COPY command with 1 -chown and 1 -chmod flag.", copyNode: mockCopyNode{
			original: "COPY -chown=user:user -chmod=7780 src dst/",
			flags: []string{},
			next: []string{"-chown=user:user", "-chmod=7780", "src", "dst/"},
		}},
		{isViolation: false, name: "Strange COPY command with 1 --chmod flag.", copyNode: mockCopyNode{
			original: "COPY --chmod=7780 chmod chmod.bak/",
			flags: []string{"--chmod=7780"},
			next: []string{"chmod", "chmod.bak"},
		}},
		{isViolation: true, name: "Strange COPY command with 1 -chmod flag.", copyNode: mockCopyNode{
			original: "COPY -chmod=7780 chmod chmod.bak/",
			flags: []string{},
			next: []string{"-chmod=7780", "chmod", "chmod.bak"},
		}},
		// buildkit responses with Parser error for the following, as it ends up with a flag --chmod without value
		// {isViolation: true, name: "Strange COPY command with 1 --chmod = flag.", copyNode: mockCopyNode{
		// 	original: "COPY --chmod = 7780 chmod chmod.bak/",
		// 	flags: []string{"--chmod"},
		// 	next: []string{"=", "7780", "chmod", "chmod.bak"},
		// }},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// nolint:exhaustivestruct
			node := &parser.Node{
				Value: "copy",
				Original: testCase.copyNode.original,
				Flags:    testCase.copyNode.flags,
			}

			// fill out Next values down the tree
			nextPointer := &node.Next
			for _, next := range testCase.copyNode.next {
				*nextPointer = &parser.Node{Value: next} // nolint:exhaustivestruct
				nextPointer = &(*nextPointer).Next
			}

			instruction, err := instructions.ParseInstruction(node)
			if err != nil {
				t.Error("cannot parse mock CopyCommand AST Node")
			}
			command, ok := instruction.(*instructions.CopyCommand)
			if !ok {
				t.Error("cannot type assert instruction to *instructions.CopyCommand")
			}

			assert.Equal(t, testCase.isViolation, RuleSet.ValidateCpy001(command).IsViolated())
		})
	}
}
