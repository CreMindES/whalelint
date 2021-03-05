package ruleset_test

import (
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

// TestValidateCPY004 tests COPY command validation
//
// Scenario: COPY command with multiple sources and one destination
//
// GIVEN a copy command WHEN |  sources are  | destination is | THEN | this should be {}.
//                           |     src1      |     dest       |      |        VALID       |
//                           |     src1      |     dest/      |      |        VALID       |
//                           |   src1 src2   |     dest       |      |       INVALID      |
//                           |   src1 src2   |     dest/      |      |        VALID       |
//                           | -chmod=7 src1 |     dest/      |      |        VALID       |
func TestValidateCpy004(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		CommandParam string
		IsViolation  bool
		ExampleName  string
		DocsContext  string
	}{
		{"src1 dst1 " /**/, true,  "COPY src1 dst1", "FROM golang:1.15\nCOPY {{ .CommandParam }}"},
		{"src1      dst1/", true,  "COPY src1      dst1", "FROM golang:1.15\nCOPY {{ .CommandParam }}"},
		{"src1 src2 dst1 ", false, "COPY src1 src2 dst1", "FROM golang:1.15\nCOPY {{ .CommandParam }}"},
		{"src1 src2 dst1/", true,  "COPY src1 src2 dst1/", "FROM golang:1.15\nCOPY {{ .CommandParam }}"},
		{"-chmod=7 src2 dst1/", true, "COPY -chmod=7 src1 dst1/", "FROM golang:1.15\nCOPY {{ .CommandParam }}"},
	}

	RuleSet.RegisterTestCaseDocs("CPY004", testCases)

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.CommandParam, func(t *testing.T) {
			t.Parallel()

			command := &instructions.CopyCommand{
				SourcesAndDest: strings.Fields(testCase.CommandParam),
				From:           "",
				Chown:          "",
				Chmod:          "",
			}

			result := !RuleSet.ValidateCpy004(command).IsViolated()

			assert.Equal(t, testCase.IsViolation, result)
		})
	}
}
