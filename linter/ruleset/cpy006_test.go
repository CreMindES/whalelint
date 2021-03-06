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
			IsViolation: false, ExampleName: "1st stage name is `foo`, copy from `bar`.",
			StageName: "foo", CopyFrom: "bar",
			DocsContext: `FROM golang:1.15 as {{ .CopyFrom }}
                          RUN go build app
                          FROM ubuntu:20.14 as {{ .StageName }}
                          COPY --from {{ .CopyFrom }}`,
		},
		{
			IsViolation:  true, ExampleName: "2nd stage name is `foo`, copy from `foo`.",
			StageName: "foo", CopyFrom: "foo",
			DocsContext: `FROM golang:1.15 as bar
                          RUN go build app
                          FROM ubuntu:20.14 as {{ .StageName }}
                          COPY --from {{ .CopyFrom }}`,
		},
		{
			IsViolation: false, ExampleName: "No stage name, but copy from `bar`",
			StageName: "", CopyFrom: "foo",
			DocsContext: `FROM golang:1.15
                          RUN go build app
                          FROM ubuntu:20.14
                          COPY --from {{ .CopyFrom }}`,
		},
		{
			IsViolation: false, ExampleName: "1st stage name is `fooBar`, copy from `foo`.",
			StageName: "fooBar", CopyFrom: "foo",
			DocsContext: `FROM golang:1.15 as {{ .StageName }}
                          RUN go build app
                          FROM ubuntu:20.14
                          COPY --from {{ .CopyFrom }}`,
		},
		{
			IsViolation: false, ExampleName: "1st stage name is `foo`, copy from `fooBar`.",
			StageName: "foo", CopyFrom: "fooBar",
			DocsContext: `FROM golang:1.15 as {{ .StageName }}
                          RUN go build app
                          FROM ubuntu:20.14
                          COPY --from {{ .CopyFrom }}`,

		},
		{
			IsViolation: false, ExampleName: "1st stage name is foo, copy from `foo:1.2`.",
			StageName: "foo", CopyFrom: "foo:1.2",
			DocsContext: `FROM golang:1.15 as {{ .StageName }}
                          RUN go build app
                          FROM ubuntu:20.14
                          COPY --from {{ .CopyFrom }}`,
		},
		{
			IsViolation: true,
			ExampleName: "1st stage alias is `builder` and 2nd base image is `foo`, copy from `foo:latest`.",
			StageName: "builder", StageBase: "foo", CopyFrom: "foo:latest",
			DocsContext: `FROM golang:1.15 as {{ .StageName }}
                          RUN go build app
                          FROM {{ .StageBase }}
                          COPY --from {{ .CopyFrom }}`,

		},
		{
			IsViolation: true,
			ExampleName: "1st stage alias is `builder` and 2nd base image is `foo:latest`, copy from `foo`.",
			StageName: "builder", StageBase: "foo:latest", CopyFrom: "foo",
			DocsContext: `FROM golang:1.15 as {{ .StageName }}
                          RUN go build app
                          FROM {{ .StageBase }}
                          COPY --from {{ .CopyFrom }}`,
		},
		{
			IsViolation: false,
			ExampleName: "Simple COPY src dst",
			StageName: "", StageBase: "foo:latest", CopyFrom: "",
			DocsContext: `FROM golang:1.15 as {{ .StageName }}
                          RUN go build app
                          FROM {{ .StageBase }}
                          COPY src dst`,
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
