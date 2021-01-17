package ruleset_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

// Test that file names, validation function names and rule names match.
func TestRulesetNamingScheme(t *testing.T) {
	t.Parallel()

	packagePath := "./"

	validationFnNameRegexp := regexp.MustCompile("Validate([A-Z][a-z]{2}[0-9]{3})")
	filenameRegexp := regexp.MustCompile("[a-z]{3}[0-9]{3}.go")
	// ruleNameRegexp := regexp.MustCompile("[A-Z]{3}[0-9]{3}.go")

	// Parse package list on packagePath
	packageList := parsePackageList(t, packagePath)

	// Get all the validation functions and their file names
	validationFnMap := parseValidationFnMap(t, packageList)

	// Validate rule function and filename naming convention:
	// - Rule name: [A-Z]{3}[0-9]{3}
	// - File name: ruleName.toLower() + ".go", i.e. [a-z]{3}[0-9]{3}.go
	// - ValidationFn name: "Validate" + rule name as [A-Z][A-Z]{2}[0-9]{3}

	for key, ruleNameCheck := range validationFnMap {
		// Check ValidationFunction name.
		if !validationFnNameRegexp.MatchString(ruleNameCheck.validationFuncName) {
			t.Error("Validation function name \"", ruleNameCheck.validationFuncName,
				"\" does not conform to naming convention.")
			t.FailNow()
		}

		// Check filename.
		if !filenameRegexp.MatchString(ruleNameCheck.filename) {
			t.Error("Validation rule filename \"", ruleNameCheck.filename,
				"\" does not conform to naming convention.")
			t.FailNow()
		}

		// Check that ValidationFunction is declared in a file with the right name.
		validateFuncStr := validationFnNameRegexp.FindStringSubmatch(ruleNameCheck.validationFuncName)
		validateFuncStrRuleNamePart := strings.ToLower(validateFuncStr[1])
		filenameCandidate := strings.ToLower(validateFuncStrRuleNamePart) + ".go"

		if filenameCandidate != ruleNameCheck.filename {
			t.Error("Validation function name does not match with its filename.",
				ruleNameCheck.validationFuncName, "shouldn't be in", ruleNameCheck.filename)
			t.FailNow()
		}

		// Check for validation function that is defined but not used
		if ruleNameCheck.newRuleCall.ruleID == "" {
			t.Log("Rule declared but not added to RuleMap!")
			t.SkipNow()
		}

		// Check that NewRule was call with matching Validation function
		if (ruleNameCheck.validationFuncName != ruleNameCheck.newRuleCall.validationFuncName) ||
			(key != strings.ToLower(ruleNameCheck.newRuleCall.ruleID)) {
			t.Error("Naming mismatch for rule", key, "and", ruleNameCheck.newRuleCall.ruleID)
		}
	}
}

func parsePackageList(t *testing.T, path string) map[string]*ast.Package {
	t.Helper()

	packList, err := parser.ParseDir(token.NewFileSet(), path, nil, 0)
	if err != nil {
		t.Error("Failed to parse package.", err)
	}

	return packList
}

type RuleNameCheck struct {
	filename           string
	validationFuncName string
	// the actual id given to the rule when called NewRule
	newRuleCall struct {
		ruleID             string
		validationFuncName string
	}
}

type RuleNameCheckMap map[string]*RuleNameCheck

func parseValidationFnMap(t *testing.T, packageList map[string]*ast.Package) RuleNameCheckMap {
	t.Helper()

	// map of ruleID.toLower() - ruleNameCheck
	fnMap := RuleNameCheckMap{}

	for _, pack := range packageList {
		for filePath, fileAst := range pack.Files {
			// TODO: think about making this mess nicer
			filename := filepath.Base(filePath)

			for _, d := range fileAst.Decls {
				if fn, isFn := d.(*ast.FuncDecl); isFn { //nolint:nestif
					functionName := fn.Name.Name

					if filterFuncByReturnType(fn, RuleSet.RuleValidationResult{}) {
						// Skip the special, main, func Validate(param interface{}) ValidationResult.
						if functionName == "Validate" {
							continue
						}

						key := strings.ToLower(strings.TrimPrefix(functionName, "Validate"))

						fnMap.update(key, filename, functionName, "", "")
					}
				} else if newRuleCallAst := filterAstNewRuleCall(d); newRuleCallAst != nil {
					if basicLit, isBacislit := newRuleCallAst.Args[0].(*ast.BasicLit); isBacislit {
						ruleID := RemoveQuotes(basicLit.Value)

						validationFuncName := newRuleCallAst.Args[3].(*ast.Ident).String()
						key := strings.ToLower(ruleID)

						fnMap.update(key, filename, "", ruleID, validationFuncName)
					}
				}
			}
		}
	}

	return fnMap
}

func RemoveQuotes(str string) string {
	return str[1 : len(str)-1]
}

// nolint
func (rncm *RuleNameCheckMap) update(key, filename,
	validationFuncName,
	ruleID,
	validationFuncName2 string) {

	if _, ok := (*rncm)[key]; !ok {
		(*rncm)[key] = &RuleNameCheck{}
	}

	if            filename != "" {(*rncm)[key].filename = filename}
	if  validationFuncName != "" {(*rncm)[key].validationFuncName = validationFuncName}
	if              ruleID != "" {(*rncm)[key].newRuleCall.ruleID = ruleID}
	if validationFuncName2 != "" {(*rncm)[key].newRuleCall.validationFuncName = validationFuncName2}
}

func filterAstNewRuleCall(decl ast.Decl) *ast.CallExpr {
	if genDecl, isGenDecl := decl.(*ast.GenDecl); isGenDecl { //nolint:nestif
		if valueSpec, isValueSpec := genDecl.Specs[0].(*ast.ValueSpec); isValueSpec {
			if callExp, isCallExp := valueSpec.Values[0].(*ast.CallExpr); isCallExp {
				if ident, isIdent := callExp.Fun.(*ast.Ident); isIdent {
					if ident.Name == "NewRule" {
						return callExp
					}
				}
			}
		}
	}

	return nil
}

func filterFuncByReturnType(fn *ast.FuncDecl, t interface{}) bool {
	returnTypeName := "notWhatWeAreLookingFor"

	targetReturnTypeName := reflect.TypeOf(t).Name()

	// check for no return type
	if fn.Type.Results == nil {
		return targetReturnTypeName == "nil"
	}

	funcReturnResultType := fn.Type.Results.List[0].Type
	if ident, ok := funcReturnResultType.(*ast.Ident); ok {
		returnTypeName = ident.Name
	}

	return returnTypeName == reflect.TypeOf(t).Name()
}
