package ruleset_test

import (
	"strings"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateCpy005(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		commandParam string
		isViolation  bool
		name         string
	}{
		{commandParam: "foo/bar /tmp/"       , isViolation: false, name: "Standard COPY."  },
		{commandParam: "foo/bar.tar.gz /tmp/", isViolation:  true, name: "COPY \".tar.gz\""},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			command := &instructions.CopyCommand{
				SourcesAndDest: strings.Fields(testCase.commandParam),
				From:           "",
				Chown:          "",
				Chmod:          "",
			}

			assert.Equal(t, testCase.isViolation, RuleSet.ValidateCpy005(command).IsViolated())
		})
	}
}
