package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateCpy006(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		stageName   string
		stageBase   string
		copyFrom    string
		isViolation bool
		name        string
	}{
		{stageName: "foo",    copyFrom: "bar",        isViolation: false, name: "Stage name=foo, copy --from=bar."},
		{stageName: "foo",    copyFrom: "foo",        isViolation:  true, name: "Stage name=foo, copy --from=foo."},
		{stageName: "",       copyFrom: "foo",        isViolation: false, name: "Stage name=\"\", copy --from=bar."},
		{stageName: "fooBar", copyFrom: "foo",        isViolation: false, name: "Stage name=fooBar, copy --from=foo."},
		{stageName: "foo",    copyFrom: "fooBar",     isViolation: false, name: "Stage name=foo, copy --from=fooBar."},
		{stageName: "foo",    copyFrom: "foo:1.2",    isViolation: false, name: "Stage name=foo, copy --from=foo:1.2."},
		{
			stageName: "builder",
			stageBase: "foo",
			copyFrom: "foo:latest",
			isViolation:  true,
			name: "Stage name=foo, copy --from=foo:latest."},
		{
			stageName: "builder",
			stageBase: "foo:latest",
			copyFrom: "foo",
			isViolation:  true,
			name: "Stage name=foo, copy --from=foo:latest."},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// nolint:exhaustivestruct
			stage := instructions.Stage{
				Name:     testCase.stageName,
				BaseName: testCase.stageBase,
				Commands: []instructions.Command{
					&instructions.CopyCommand{
						From: testCase.copyFrom,
					},
				},
			}

			assert.Equal(t, testCase.isViolation, RuleSet.ValidateCpy006(stage).IsViolated())
		})
	}
}
