// ruleset provides a set of rules and their corresponding validator functions
// for linting Dockerfile AST elements.
package ruleset

import (
	"encoding/json"
	"reflect"

	log "github.com/sirupsen/logrus"
)

// Severity type represents a severity, with an int level and a String function.
type Severity int

const (
	ValError Severity = iota
	ValDeprecation
	ValInfo
	ValWarning
)

// Severity.String() converts the raw Severity into a string.
func (severity Severity) String() string {
	switch severity {
	case ValDeprecation:
		return "Deprecation"
	case ValError:
		return "Error"
	case ValInfo:
		return "Info"
	case ValWarning:
		return "Warning"
	default:
		return "Unknown"
	}
}

// DocsReference returns an official reference link connected to the rule itself, most likely directly linking to a
// Docker documentation webpage.
func (rule *Rule) DocsReference() DocsReference {
	docsReference, ok := docsReferenceMap[rule.id[:3]]
	if !ok {
		return ToDoReference
	}

	return docsReference
}

// Rule represents a Dockerfile lint validation rule.
// It has the basic id, description, severity attributes and a validation function as an
// empty interface. For further details on validateFunc, please see Validate how it is
// utilized.
type Rule struct {
	id             string
	description    string
	severity       Severity
	validationFunc interface{}
}

// Validation calls the the rule's validationFunc validation function
// in the correct form, after converting from interface{} to the concrete type.
//
// example: func(runCommand *instructions.RunCommand) RuleValidationResult where runCommand is
// asserted param as *instructions.RunCommand.
func (rule Rule) Validate(param interface{}) RuleValidationResult {
	// Assemble validationFunc reflect type, based on param type, as they are always
	// func(param *paramActualType) RuleValidationResult
	paramType := reflect.TypeOf(param)
	returnType := reflect.TypeOf(RuleValidationResult{})
	funcType := reflect.FuncOf([]reflect.Type{paramType}, []reflect.Type{returnType}, false)
	funcReflect := reflect.ValueOf(rule.validationFunc).Convert(funcType)
	log.Trace("RuleSet | ValidationReflect> funcType:", funcType)

	// Type assert param into the actual type
	paramCasted := reflect.ValueOf(param).Convert(paramType)

	// Call the reflection function representation
	funcReflectResult := funcReflect.Call([]reflect.Value{paramCasted})

	// Get back actual result and assign rule to rule validation result
	result, ok := funcReflectResult[0].Interface().(RuleValidationResult)
	if ok {
		result.rule = rule
	} else {
		log.Error("Cannot retrieve RuleValidationResult from reflect call.")
	}

	return result
}

// NewRule creates a new Rule by joining it's id, description, severity and validation function.
// It automatically gets assigned into a slice/set of rules corresponding to a specific
// Dockerfile AST element, inside the ruleMap's corresponding bin, based on the Dockerfile AST
// element's type. See below, how reflect.TypeOf().String() is used to achieve this.
func NewRule(id string, description string, severity Severity, param interface{}) interface{} {
	rule := Rule{
		id:             id,
		description:    description,
		severity:       severity,
		validationFunc: param,
	}

	targetBin := reflect.TypeOf(param).In(0).String()

	if val, ok := ruleMap[targetBin]; ok {
		ruleMap[targetBin] = append(val, rule)
	} else {
		ruleMap[targetBin] = []Rule{rule}
	}

	return true
}

// ID returns the rule's id string.
func (rule *Rule) ID() string {
	return rule.id
}

// Severity returns the rule's severity, e.g. error, warning, info, deprecation.
func (rule *Rule) Severity() Severity {
	return rule.severity
}

// Description returns the rule's description, i.e. the rule itself as a statement/guidance.
func (rule *Rule) Description() string {
	return rule.description
}

func (rule *Rule) ValidationFunc() interface{} {
	return rule.validationFunc
}

// MarshalJSON converts a Rule instance to JSON.
func (rule Rule) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		ID          string
		Description string
		Severity    string
	}{
		ID:          rule.id,
		Description: rule.description,
		Severity:    rule.severity.String(),
	})
}

// RuleMapType represents a set of rules for each Dockerfile AST element
// identified by its type's string value (through reflect).
type RuleMapType map[string][]Rule

// ruleMap stores a ruleset for each Dockerfile AST element that they need to be validated against.
// It's a map[reflect.TypeOf(Dockerfile AST element).String()][]Rule under the hood.
var ruleMap RuleMapType = map[string][]Rule{} // nolint:gochecknoglobals

// Count gives back the total number of rules in the ruleset.
// Note: each AST element has a set of corresponding rules in the rule map.
func (rm RuleMapType) Count() int {
	sum := 0
	for _, astElementRuleList := range rm {
		sum += len(astElementRuleList)
	}

	return sum
}

// Get returns ruleset's ruleMap.
func Get() RuleMapType {
	return ruleMap
}

// GetRulesForAstElement returns a Rule slice with all the rules that
// the given Dockerfile AST element needs to be validated against.
func GetRulesForAstElement(astElementInterface interface{}) []Rule {
	return ruleMap[reflect.TypeOf(astElementInterface).String()]
}
