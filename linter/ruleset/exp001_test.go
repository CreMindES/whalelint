package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateExp001(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		PortValue   []string
		IsViolation bool
		ExampleName string
		DocsContext string
	}{
		{
			PortValue:   []string{"4242"},
			IsViolation: false,
			ExampleName: "EXPOSE 4242",
			DocsContext: "FROM golang:1.15\nEXPOSE {{ .PortValue }}",
		},
		{
			PortValue:   []string{"4242/tcp"}, IsViolation: false,
			ExampleName: "EXPOSE 4242/tcp", DocsContext: "FROM golang:1.15\nEXPOSE {{ .PortValue }}",
		},
		{
			PortValue:   []string{"4242/udp"}, IsViolation: false,
			ExampleName: "EXPOSE 4242/udp", DocsContext: "FROM golang:1.15\nEXPOSE {{ .PortValue }}",
		},
		{
			PortValue:   []string{"4242/yyy"}, IsViolation: true,
			ExampleName: "EXPOSE 4242/yyy", DocsContext: "FROM golang:1.15\nEXPOSE {{ .PortValue }}",
		},
		{
			PortValue:   []string{"4242:tcp"}, IsViolation: true,
			ExampleName: "EXPOSE 4242:tcp", DocsContext: "FROM golang:1.15\nEXPOSE {{ .PortValue }}",
		},
		{
			PortValue:   []string{"4242", "4242/tcp", "4242/udp"}, IsViolation: false,
			ExampleName: "EXPOSE 4242, 4242/tcp, 4242/udp", DocsContext: "FROM golang:1.15\nEXPOSE {{ .PortValue }}",
		},
		{
			PortValue:   []string{"67999"}, IsViolation: true,
			ExampleName: "EXPOSE 67999", DocsContext: "FROM golang:1.15\nEXPOSE {{ .PortValue }}",
		},
		{
			PortValue:   []string{"4242", "67999", "4242/udp" }, IsViolation: true,
			ExampleName: "EXPOSE 4242, 67999, 4242/udp", DocsContext: "FROM golang:1.15\nEXPOSE {{ .PortValue }}",
		},
	}

	RuleSet.RegisterTestCaseDocs("EXP001", testCases)

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.ExampleName, func(t *testing.T) {
			t.Parallel()

			command := &instructions.ExposeCommand{
				Ports: testCase.PortValue,
			}

			assert.Equal(t, testCase.IsViolation, RuleSet.ValidateExp001(command).IsViolated())
		})
	}
}
