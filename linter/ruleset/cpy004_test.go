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

	testCases := []struct {
		commandParam string
		expected     bool
	}{
		{"src1 dst1 " /**/, true},
		{"src1      dst1/", true},
		{"src1 src2 dst1 ", false},
		{"src1 src2 dst1/", true},
		{"-chmod=7 src2 dst1/", true},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.commandParam, func(t *testing.T) {
			t.Parallel()

			command := &instructions.CopyCommand{
				SourcesAndDest: strings.Fields(testCase.commandParam),
				From:           "",
				Chown:          "",
				Chmod:          "",
			}

			result := !RuleSet.ValidateCpy004(command).IsViolated()

			assert.Equal(t, testCase.expected, result)
		})
	}
}
