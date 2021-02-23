package ruleset

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/Masterminds/sprig"
	Log "github.com/sirupsen/logrus"
)

// RuleDocMap is a collection of RuleID -> RuleDoc(s).
type RuleDocMap map[string]RuleDoc

// RuleDoc represent docs for a single lint Rule.
// It is primarily intended to be source for automatic rule docs generation.
// It's aggregated from the rule itself and it's test cases.
type RuleDoc struct {
	Rule     *Rule
	DocsRef  DocsReference
	TestDocs []TestCaseDocs
}

// TestCaseDocs holds common test case parts for lint rules, mainly used in RuleDoc.
type TestCaseDocs struct {
	ExampleName string
	DocsContext string
	IsViolation bool
}

var ruleDocMap = RuleDocMap{} // nolint:gochecknoglobals

// ExtractDocFieldsFromTestCase parses the common fields of an arbitrary test case struct into a TestCaseDocs.
// It return an error is any of the fields of a TestCaseDocs is missing from a lint Rule test case.
func ExtractDocFieldsFromTestCase(testDocsReflect reflect.Value, parent reflect.Value, i int) (TestCaseDocs, error) {
	// Parse common TestCaseDocs <-> Rule lint TestCase struct fields
	docsContextReflectValue := testDocsReflect.FieldByName("DocsContext")
	exampleNameReclectValue := testDocsReflect.FieldByName("ExampleName")
	isViolationReflectValue := testDocsReflect.FieldByName("IsViolation")

	// Check success/validity of fields
	if !docsContextReflectValue.IsValid() ||
		!isViolationReflectValue.IsValid() ||
		!exampleNameReclectValue.IsValid() {
		// TODO: error
		err := &reflect.ValueError{}

		return TestCaseDocs{}, err
	}

	// Convert to actual types from reflect.Value
	docsContext := docsContextReflectValue.String()
	exampleName := exampleNameReclectValue.String()
	isViolation := isViolationReflectValue.Bool()

	// Run template engine
	docsContext, err := ApplyTemplate(docsContext, parent.Index(i).Interface())
	if err != nil {
		return TestCaseDocs{}, err
	}

	// Remove extra tabs and spaces
	space := regexp.MustCompile(`[^\S\n]{2,}`)
	docsContext = space.ReplaceAllString(docsContext, "")

	// Assemble result
	testCaseDocs := TestCaseDocs{
		ExampleName: exampleName,
		DocsContext: docsContext,
		IsViolation: isViolation,
	}

	return testCaseDocs, nil
}

// ApplyTemplate generates text from templateStr and sourceStruct.
func ApplyTemplate(templateStr string, sourceStruct interface{}) (string, error) {
	templateEngine, _ := template.New("ruleDocs").Parse(templateStr)
	strBuilder := strings.Builder{}

	if err := templateEngine.Execute(&strBuilder, sourceStruct); err != nil {
		return "", fmt.Errorf("failed to execute template engine: %w", err)
	}

	return strBuilder.String(), nil
}

// RegisterTestCaseDocs registers a test case docs for Rule DocsContext generation.
func RegisterTestCaseDocs(ruleID string, testDocsIntrfc interface{}) {
	rule := Get().GetRuleByName(ruleID, nil)
	testDocsSlice := make([]TestCaseDocs, 0)

	// get DocsContext and IsViolation values
	switch reflect.TypeOf(testDocsIntrfc).Kind() { // nolint:exhaustive
	case reflect.Slice:
		s := reflect.ValueOf(testDocsIntrfc)

		for i := 0; i < s.Len(); i++ {
			testDocsReflect := s.Index(i)

			testCaseDocs, err := ExtractDocFieldsFromTestCase(testDocsReflect, s, i)
			if err != nil {
				Log.Error(err) // TODO
				return         // nolint:nlreturn
			}

			testDocsSlice = append(testDocsSlice, testCaseDocs)
		}
	default:
		// TODO: handle error gracefully
		return
	}

	ruleDocMap[ruleID] = RuleDoc{
		Rule:     &rule,
		TestDocs: testDocsSlice,
		DocsRef:  rule.DocsReference(),
	}

	GenerateRuleDocs()
}

func GenerateRuleDocs() {
	folder := "./" // ./linter/ruleset/

	// docsTemplate, err := template.ParseFiles(folder + "docs.gotemplate")
	// if err != nil {
	// 	// TODO: err
	// 	Log.Error(err)
	// }

	docsTemplate := template.Must(
		template.New("docs.gotemplate").Funcs(sprig.FuncMap()).ParseGlob("docs.gotemplate"),
	)

	for ruleID, ruleDoc := range ruleDocMap {
		if ruleID == "STL001" || ruleID == "CPY001" || ruleID == "CPY006" {
			f, err := os.Create(folder + strings.ToLower(ruleID) + ".md")
			if err != nil {
				return
			}

			w := bufio.NewWriter(f)

			// Run the template engine
			err = docsTemplate.Execute(w, ruleDoc)
			if err != nil {
				Log.Error(err)
			}

			// Flush writer
			if errFlush := w.Flush(); errFlush != nil {
				Log.Error(errFlush)
			}

			// Close file
			if errClose := f.Close(); errClose != nil {
				Log.Error(errClose)
			}

			break
		}
	}
}
