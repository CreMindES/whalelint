package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateSts003(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		IsViolation   bool
		StagePlatform string
		ExampleName   string
		DocsContext   string
	}{
		{
			IsViolation:   false,
			ExampleName:   "One stage FROM ubuntu:20.04.",
			StagePlatform: "",
			DocsContext:   "`FROM` {{ .StageBaseName }}",
		},
		{
			IsViolation:   true,
			ExampleName:   "One stage FROM --platform=armv7 ubuntu:20.14.",
			StagePlatform: "armv7",
			DocsContext:   "`FROM` {{ .StageBaseName }}",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.ExampleName, func(t *testing.T) {
			t.Parallel()

			stageSourceCode := "FROM "
			if len(testCase.StagePlatform) > 0 {
				stageSourceCode += "--platform=" + testCase.StagePlatform + " "
			}
			stageSourceCode += "ubuntu:20.14"

			// nolint:exhaustivestruct
			stage := instructions.Stage{
				BaseName:   "ubuntu:20.14",
				Platform:   testCase.StagePlatform,
				SourceCode: stageSourceCode,
			}

			assert.Equal(t, testCase.IsViolation, RuleSet.ValidateSts003(stage).IsViolated())
		})
	}
}
