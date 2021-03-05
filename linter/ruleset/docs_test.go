// +build ruledocs

package ruleset_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestExtractDocFieldsFromTestCase(t *testing.T) {
	t.Parallel()

	testCaseMockOK := RuleSet.TestCaseDocs{
		ExampleName: "exampleName",
		DocsContext: "docsContext",
		IsViolation: false,
	}

	testCaseMockNotOk := struct {
		RandomField bool
	}{
		RandomField: false,
	}

	testCaseMockOkWrongTemplate := struct {
		ExampleName string
		DocsContext string
		IsViolation bool
		RandomField bool
	}{
		ExampleName: "exampleName",
		DocsContext: "{{ .NotRandomField }}",
		IsViolation: false,
		RandomField: false,
	}

	testCaseMockSlice := []interface{}{testCaseMockOK, testCaseMockNotOk, testCaseMockOkWrongTemplate}

	// testCaseMockOK
	testCaseDocs, err := RuleSet.ExtractDocFieldsFromTestCase(reflect.ValueOf(testCaseMockOK),
		reflect.ValueOf(testCaseMockSlice), 0)

	assert.Equal(t, nil, err)
	assert.Equal(t, testCaseMockOK, testCaseDocs)

	// testCaseMockNotOk
	_, err = RuleSet.ExtractDocFieldsFromTestCase(reflect.ValueOf(testCaseMockNotOk),
		reflect.ValueOf(testCaseMockSlice), 1)

	assert.NotEqual(t, nil, err) // TODO

	// testCaseMockOkWrongTemplate
	_, err = RuleSet.ExtractDocFieldsFromTestCase(reflect.ValueOf(testCaseMockOkWrongTemplate),
		reflect.ValueOf(testCaseMockSlice), 2)

	assert.NotEqual(t, nil, err) // TODO
}

func TestRegisterTestCaseDocs(t *testing.T) {
	t.Parallel()
}

func TestGenerateRuleDocs(t *testing.T) {
	t.Parallel()
}
