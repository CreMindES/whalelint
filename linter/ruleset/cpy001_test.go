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
			original: "COPY --chmod=7780 src src dst/",
			flags: []string{"--chmod=7780"},
			next: []string{"src", "src", "dst/"},
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
				Next: &parser.Node{
					Value: testCase.copyNode.next[0],
					Next: &parser.Node{
						Value: testCase.copyNode.next[1],
					},
				},
				Original: testCase.copyNode.original,
				Flags:    testCase.copyNode.flags,
			}

			pointer := node.Next
			for _, next := range testCase.copyNode.next {
				nodePtr := &parser.Node{Value: next} // nolint:exhaustivestruct
				pointer = nodePtr
				pointer = (*pointer).Next
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
