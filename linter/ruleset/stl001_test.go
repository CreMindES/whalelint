package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateStl001(t *testing.T) { // nolint:funlen
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		StageNameList []string
		IsViolation   bool
		ExampleName   string
		DocsContext   string
	}{
		{
			StageNameList: []string{"builder"},
			IsViolation:   false,
			ExampleName:   "One stage with alias.",
			DocsContext:   `FROM golang:1.15 as {{ index .StageNameList 0 }}
			                RUN go --version`,
		},
		{
			StageNameList: []string{"builder_foo", "builder_bar"},
			IsViolation:   false,
			ExampleName:   "Two stages with aliases.",
			DocsContext:   `FROM golang:1.15 as {{ index .StageNameList 0 }}
						    RUN go build app
						    FROM ubuntu:20.04 as {{ index .StageNameList 1 }}
						    COPY --from {{ index .StageNameList 0 }} /app ./app`,
		},
		{
			StageNameList: []string{"builder_foo", "builder_foo" },
			IsViolation:   true,
			ExampleName:   "Two stages with the same aliases.",
			DocsContext:   `FROM golang:1.15 as {{ index .StageNameList 0 }}
                            RUN go build app
                            FROM ubuntu:20.04 as {{ index .StageNameList 1 }}
							COPY --from {{ index .StageNameList 0 }} /app ./app`,
		},
		{
			StageNameList: []string{"builder_foo", "", ""},
			IsViolation:   false,
			ExampleName:   "Three stages, but only one has an alias.",
			DocsContext:   `FROM golang:1.15 as {{ index .StageNameList 0 }}
                            RUN go build app
                            FROM golang:1.16
                            RUN go build app
                            FROM scratch
							COPY --from {{ index .StageNameList 0 }} /app ./app`,
		},
	}

	RuleSet.RegisterTestCaseDocs("STL001", testCases)

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.ExampleName, func(t *testing.T) {
			t.Parallel()

			stageList := make([]instructions.Stage, 0, len(testCase.StageNameList))
			for _, stageName := range testCase.StageNameList {
				stageList = append(stageList, instructions.Stage{Name: stageName}) // nolint:exhaustivestruct
			}

			assert.Equal(t, testCase.IsViolation, RuleSet.ValidateStl001(stageList).IsViolated())
		})
	}
}
