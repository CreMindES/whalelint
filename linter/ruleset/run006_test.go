package ruleset_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

// nolint:funlen
func TestValidateRun006(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		ExampleName string
		CommandStr  string
		IsViolation bool
		DocsContext string
	}{
		{
			IsViolation: false, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "apt-get update && apt-get install -y vim=1.2.3 && apt-get clean",
		},
		{
			IsViolation: false, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "apt-get update && apt-get install -y vim=1.2.3 && rm -rf /var/lib/apt/lists",
		},
		{
			IsViolation: false, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "apt-get update && apt-get install -y vim=1.2.3 && apt-get clean && rm -rf /var/lib/apt/lists",
		},
		{
			IsViolation: true, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "apt update && apt install vim",
		},
		{
			IsViolation: true, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "DEBIAN_FRONTEND=noninteractive apt-get update",
		},
		{
			IsViolation: true, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "DEBIAN_FRONTEND=noninteractive apt update",
		},
		{
			IsViolation: false, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "yum update -y && yum install -y git && yum clean all && date",
		},
		{
			IsViolation: true, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "yum update -y && yum install -y git && date",
		},
		{
			IsViolation: false, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "dnf update -y && dnf install -y git && dnf clean all",
		},
		{
			IsViolation: true, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "dnf update -y && dnf install -y git && date",
		},
		{
			IsViolation: false, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "zypper refresh -y && zypper install -y git && zypper clean all",
		},
		{
			IsViolation: true, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "zypper refresh -y && zypper install -y git && date",
		},
		{
			IsViolation: false, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "date",
		},
		{
			IsViolation: false, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "pip install --no-cache-dir pytorch",
		},
		{
			IsViolation: true, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "pip install --update pytorch",
		},
		{
			IsViolation: false, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "apk add --update --no-cache git",
		},
		{
			IsViolation: true, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
			CommandStr: "apk update && apk add git",
		},
		// { TODO
		// 	IsViolation: false, ExampleName: "", DocsContext: "`RUN` {{ .CommandStr }}",
		// 	CommandStr: "apk --no-cache add git",
		// },
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.CommandStr, func(t *testing.T) {
			t.Parallel()

			runCommand := RuleSet.NewRunCommand(testCase.CommandStr, RuleSet.NewLocationRange(
				1, 0, 1, len(testCase.CommandStr)))

			// test validation rule
			assert.Equal(t, testCase.IsViolation, RuleSet.ValidateRun006(runCommand).IsViolated())
		})
	}
}
