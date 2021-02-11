package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

type mockCopyNode struct {
	Original string
	flags    []string
	next     []string
}

// nolint:funlen
func TestValidateCpy001(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		IsViolation bool
		ExampleName string
		CopyNode    mockCopyNode
		DocsContext string
	}{
		{
			IsViolation: false,
			ExampleName: "Proper `COPY` command with 1 `--chmod` flag.",
			CopyNode:    mockCopyNode{
				Original: "COPY --chmod=7780 src src2 dst/",
				flags:    []string{"--chmod=7780"},
				next:     []string{"src", "src2", "dst/"},
			},
			DocsContext: "{{ .CopyNode.Original }}",
		},
		{
			IsViolation: true,
			ExampleName: "`COPY` command with 1 `-chmod` flag.",
			CopyNode:    mockCopyNode{
				Original: "COPY -chmod=7780 src dst/",
				flags:    []string{},
				next:     []string{"-chmod=7780", "src", "dst/"},
			},
			DocsContext: "{{ .CopyNode.Original }}",
		},
		{
			IsViolation: true,
			ExampleName: "`COPY` command with 1 `chmod` flag.",
			CopyNode:    mockCopyNode{
				Original: "COPY chmod=7780 src dst/",
				flags:    []string{},
				next:     []string{"chmod=7780", "src", "dst/"},
			},
			DocsContext: "{{ .CopyNode.Original }}",
		},
		{
			IsViolation: true,
			ExampleName: "`COPY` command with 1 `-chown` and 1 `-chmod` flag.",
			CopyNode:    mockCopyNode{
				Original: "COPY -chown=user:user -chmod=7780 src dst/",
				flags:    []string{},
				next:     []string{"-chown=user:user", "-chmod=7780", "src", "dst/"},
			},
			DocsContext: "{{ .CopyNode.Original }}",
		},
		{
			IsViolation: false,
			ExampleName: "Strange `COPY` command with 1 `--chmod` flag.",
			CopyNode:    mockCopyNode{
				Original: "COPY --chmod=7780 chmod chmod.bak/",
				flags:    []string{"--chmod=7780"},
				next:     []string{"chmod", "chmod.bak"},
			},
			DocsContext: "{{ .CopyNode.Original }}",
		},
		{
			IsViolation: true,
			ExampleName: "Strange `COPY` command with 1 `-chmod` flag.",
			CopyNode:    mockCopyNode{
				Original: "COPY -chmod=7780 chmod chmod.bak/",
				flags:    []string{},
				next:     []string{"-chmod=7780", "chmod", "chmod.bak"},
			},
			DocsContext: "{{ .CopyNode.Original }}",
		},
		// buildkit responses with Parser error for the following, as it ends up with a flag --chmod without value
		// {
		// 	IsViolation: true,
		// 	ExampleName: "Strange COPY command with 1 --chmod = flag.",
		// 	CopyNode:    mockCopyNode{
		// 		Original: "COPY --chmod = 7780 chmod chmod.bak/",
		// 		flags: []string{"--chmod"},
		// 		next: []string{"=", "7780", "chmod", "chmod.bak"},
		// 	},
		// 	DocsContext: "{{ .CopyNode.Original }}",
		// },
	}

	RuleSet.RegisterTestCaseDocs("CPY001", testCases)

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.ExampleName, func(t *testing.T) {
			t.Parallel()

			// nolint:exhaustivestruct
			node := &parser.Node{
				Value:    "copy",
				Original: testCase.CopyNode.Original,
				Flags:    testCase.CopyNode.flags,
			}

			// fill out Next values down the tree
			nextPointer := &node.Next
			for _, next := range testCase.CopyNode.next {
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

			assert.Equal(t, testCase.IsViolation, RuleSet.ValidateCpy001(command).IsViolated())
		})
	}
}
