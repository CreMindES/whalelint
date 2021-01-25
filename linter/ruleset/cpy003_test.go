package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateCpy003(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		chownValue  string
		isViolation bool
		name        string
	}{ // valid examples are from https://docs.docker.com/engine/reference/builder/#copy
		{chownValue: "55:mygroup", isViolation: false, name: "COPY with chown=55:mygroup"},
		{chownValue: "bin"       , isViolation: false, name: "COPY with chown=bin"       },
		{chownValue: "1"         , isViolation: false, name: "COPY with chown=1"         },
		{chownValue: "10:11"     , isViolation: false, name: "COPY with chown=10:11"     },
		{chownValue: "10;11"     , isViolation:  true, name: "COPY with chown=10;11"     },
		{chownValue: "10,11"     , isViolation:  true, name: "COPY with chown=10,11"     },
		{chownValue: "$$"        , isViolation:  true, name: "COPY with chown=$$"        },
		{chownValue: "55:11,22"  , isViolation:  true, name: "COPY with chown=55:11,22"  },
		{chownValue: "55:11 22"  , isViolation:  true, name: "COPY with chown=55:11 22"  },
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			command := &instructions.CopyCommand{
				SourcesAndDest: []string{},
				From:           "",
				Chown:          testCase.chownValue,
				Chmod:          "",
			}

			assert.Equal(t, testCase.isViolation, RuleSet.ValidateCpy003(command).IsViolated())
		})
	}
}
