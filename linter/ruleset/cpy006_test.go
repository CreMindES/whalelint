package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

// nolint:funlen
func TestValidateCpy006(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		StageName   string
		StageBase   string
		CopyFrom    string
		IsViolation bool
		ExampleName string
		DocsContext string
	}{
		{
			IsViolation: false, ExampleName: "Stage name=foo, copy --from=bar.",
			StageName: "foo", CopyFrom: "bar",
			DocsContext: `FROM golang:1.15 as {{ .StageName }}
                          RUN go build app
                          FROM ubuntu:20.14
                          COPY --from {{ .CopyFrom }}`,
		},
		{
			IsViolation:  true, ExampleName: "Stage name=foo, copy --from=foo.",
			StageName: "foo", CopyFrom: "foo",
		},
		{
			IsViolation: false, ExampleName: "Stage name=\"\", copy --from=bar.",
			StageName: "", CopyFrom: "foo",
		},
		{
			IsViolation: false, ExampleName: "Stage name=fooBar, copy --from=foo.",
			StageName: "fooBar", CopyFrom: "foo",
		},
		{
			IsViolation: false, ExampleName: "Stage name=foo, copy --from=fooBar.",
			StageName: "foo", CopyFrom: "fooBar",
		},
		{
			IsViolation: false, ExampleName: "Stage name=foo, copy --from=foo:1.2.",
			StageName: "foo", CopyFrom: "foo:1.2",
		},
		{
			IsViolation: true, ExampleName: "Stage name=foo, copy --from=foo:latest.",
			StageName: "builder", StageBase: "foo", CopyFrom: "foo:latest",

		},
		{
			IsViolation: true, ExampleName: "Stage name=foo, copy --from=foo:latest.",
			StageName: "builder", StageBase: "foo:latest", CopyFrom: "foo",
		},
	}

	RuleSet.RegisterTestCaseDocs("CPY006", testCases)

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.ExampleName, func(t *testing.T) {
			t.Parallel()

			// nolint:exhaustivestruct
			stage := instructions.Stage{
				Name:     testCase.StageName,
				BaseName: testCase.StageBase,
				Commands: []instructions.Command{
					&instructions.CopyCommand{
						From: testCase.CopyFrom,
					},
				},
			}

			assert.Equal(t, testCase.IsViolation, RuleSet.ValidateCpy006(stage).IsViolated())
		})
	}
}
