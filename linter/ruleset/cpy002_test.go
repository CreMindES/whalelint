package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateCpy002(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		chmodValue   string
		isViolation  bool
		name         string
	}{
		{chmodValue: "7440", isViolation: false, name: "COPY with chmod=7440"},
		{chmodValue:  "644", isViolation: false, name: "COPY with chmod=644" },
		{chmodValue:   "88", isViolation:  true, name: "COPY with chmod=88"  },
		{chmodValue: "7780", isViolation:  true, name: "COPY with chmod=7780"},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			command := &instructions.CopyCommand{
				SourcesAndDest: []string{},
				From:           "",
				Chown:          "",
				Chmod:          testCase.chmodValue,
			}

			assert.Equal(t, testCase.isViolation, RuleSet.ValidateCpy002(command).IsViolated())
		})
	}
}
