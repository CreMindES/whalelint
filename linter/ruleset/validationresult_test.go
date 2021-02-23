package ruleset_test

import (
	"encoding/json"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func newMockRule() *RuleSet.Rule {
	mockFunc := func(command *instructions.Command) /* RuleSet.RuleValidationResult */ {}
	mockRule := RuleSet.NewRule("FakeID", "MockDef", "MockDesc", RuleSet.ValUnknown, mockFunc)

	return mockRule
}

func newMockLocation() RuleSet.LocationRange {
	return RuleSet.NewLocationRange(1, 2, 3, 4)
}

func TestRuleValidationResult_MarshalJSON(t *testing.T) {
	t.Parallel()

	mockRule := newMockRule()
	mockLoc := newMockLocation()

	referenceRuleValidationResult := RuleSet.NewRuleValidationResult(mockRule, false, "", mockLoc)
	var duplicateRuleValidationResult RuleSet.RuleValidationResult // nolint:wsl

	// serialize
	ruleValidationResultJSON, err := json.Marshal(referenceRuleValidationResult)
	if err != nil {
		t.Error(err)
	}
	// deserialize
	err = json.Unmarshal(ruleValidationResultJSON, &duplicateRuleValidationResult)
	if err != nil {
		t.Error(err)
	}

	// First check the rule equality by fields.
	// Note: there is no point in checking the validationFunc, as unmarshaling a unknown function cannot
	//       provide the actual function body.
	// nolint: gofmt, gofumpt, goimports
	assert.Equal(t, referenceRuleValidationResult.RuleID(),      duplicateRuleValidationResult.RuleID())
	assert.Equal(t, referenceRuleValidationResult.Severity(),    duplicateRuleValidationResult.Severity())
	assert.Equal(t, referenceRuleValidationResult.Description(), duplicateRuleValidationResult.Description())

	// Putting MockRule to both validation results, so all the other fields besides Rule can be checked in one go.
	duplicateRuleValidationResult.SetRule(mockRule)

	assert.Equal(t, *referenceRuleValidationResult, duplicateRuleValidationResult)
}

func TestRuleValidationResult_SetViolated(t *testing.T) { // nolint:funlen
	t.Parallel()

	mockRule := newMockRule()
	mockLoc  := newMockLocation() // nolint: gofmt, gofumpt, goimports

	testCases := []struct {
		ExampleName      string
		validationResult *RuleSet.RuleValidationResult
		params           []bool
		IsViolation      bool
	}{
		{
			ExampleName:      "notViolated.SetViolated()",
			validationResult: RuleSet.NewRuleValidationResult(mockRule, false, "", mockLoc),
			params:           []bool{},
			IsViolation:      true,
		},
		{
			ExampleName:      "violated.SetViolated()",
			validationResult: RuleSet.NewRuleValidationResult(mockRule, true, "", mockLoc),
			params:           []bool{},
			IsViolation:      true,
		},
		{
			ExampleName:      "notViolated.SetViolated(false)",
			validationResult: RuleSet.NewRuleValidationResult(mockRule, false, "", mockLoc),
			params:           []bool{false},
			IsViolation:      false,
		},
		{
			ExampleName:      "notViolated.SetViolated(true)",
			validationResult: RuleSet.NewRuleValidationResult(mockRule, false, "", mockLoc),
			params:           []bool{true},
			IsViolation:      true,
		},
		{
			ExampleName:      "notViolated.SetViolated(true)",
			validationResult: RuleSet.NewRuleValidationResult(mockRule, true, "", mockLoc),
			params:           []bool{false},
			IsViolation:      true,
		},
		{
			ExampleName:      "notViolated.SetViolated(true)",
			validationResult: RuleSet.NewRuleValidationResult(mockRule, true, "", mockLoc),
			params:           []bool{true},
			IsViolation:      true,
		},
		{
			ExampleName:      "notViolated.SetViolated(true, FORCE)",
			validationResult: RuleSet.NewRuleValidationResult(mockRule, false, "", mockLoc),
			params:           []bool{false, RuleSet.FORCE},
			IsViolation:      false,
		},
		{
			ExampleName:      "notViolated.SetViolated(true, FORCE)",
			validationResult: RuleSet.NewRuleValidationResult(mockRule, false, "", mockLoc),
			params:           []bool{true, RuleSet.FORCE},
			IsViolation:      true,
		},
		{
			ExampleName:      "violated.SetViolated(true, FORCE)",
			validationResult: RuleSet.NewRuleValidationResult(mockRule, true, "", mockLoc),
			params:           []bool{false, RuleSet.FORCE},
			IsViolation:      false,
		},
		{
			ExampleName:      "violated.SetViolated(true, FORCE)",
			validationResult: RuleSet.NewRuleValidationResult(mockRule, true, "", mockLoc),
			params:           []bool{true, RuleSet.FORCE},
			IsViolation:      true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.ExampleName, func(t *testing.T) {
			t.Parallel()

			testCase.validationResult.SetViolated(testCase.params...)

			assert.Equal(t, testCase.IsViolation, testCase.validationResult.IsViolated())
		})
	}
	// var buffer bytes.Buffer
	// // var fields Log.Fields
	//
	// logger := Log.New()
	// logger.Out = &buffer
	//
	// logger.Println("test")

	// 	_, hook := test.NewNullLogger()
	//
	// 	valResult := RuleSet.NewRuleValidationResult(MockRule, false, "", mockLoc)
	// 	valResult.SetViolated(true, true, true)
	// // 	assert.Contains(t, "Invalid params to RuleValidationResult::SetViolated" ,buffer.String())
	// 	assert.Equal(t, Log.ErrorLevel, hook.LastEntry().Level)

	return // nolint:gosimple
}

func TestRuleValidationResult_SetLocation(t *testing.T) {
	t.Parallel()

	mockRule := newMockRule()
	mockLoc := newMockLocation()
	ruleValidationResult := RuleSet.NewRuleValidationResult(mockRule, false, "", mockLoc)

	startLine, startChar, endLine, endChar :=
		ruleValidationResult.Location().Start().LineNumber()+2,
		ruleValidationResult.Location().Start().CharNumber()+2,
		ruleValidationResult.Location().End().LineNumber()+2,
		ruleValidationResult.Location().End().CharNumber()+2

	ruleValidationResult.SetLocation(startLine, startChar, endLine, endChar)

	assert.Equal(t, startLine, ruleValidationResult.Location().Start().LineNumber())
	assert.Equal(t, startChar, ruleValidationResult.Location().Start().CharNumber())
	assert.Equal(t,   endLine, ruleValidationResult.Location().End().LineNumber()) // nolint:gofmt,gofumpt,goimports
	assert.Equal(t,   endChar, ruleValidationResult.Location().End().CharNumber())
}

func TestRuleValidationResult_SetLocationRangeFrom(t *testing.T) {
	t.Parallel()

	mockRule := newMockRule()
	mockLoc1 := RuleSet.NewLocationRange(1, 2, 3, 4)
	mockLoc2 := RuleSet.NewLocationRange(5, 6, 7, 8)
	ruleValidationResult := RuleSet.NewRuleValidationResult(mockRule, false, "", mockLoc1)

	assert.Equal(t, mockLoc1, *ruleValidationResult.Location())

	ruleValidationResult.SetLocationRangeFrom(mockLoc2)

	assert.Equal(t, mockLoc2, *ruleValidationResult.Location())
}

// nolint:gofmt,gofumpt,goimports
func TestRuleValidationResult_Message(t *testing.T) {
	t.Parallel()

	mockRule := newMockRule()
	mockLoc  := newMockLocation()
	mockMessage := "Fake Message"
	ruleValidationResultWithoutMessage := RuleSet.NewRuleValidationResult(mockRule, false, "", mockLoc)
	ruleValidationResultWithMessage    := RuleSet.NewRuleValidationResult(mockRule, false, mockMessage, mockLoc)

	assert.Equal(t, mockRule.Definition(), ruleValidationResultWithoutMessage.Message())
	assert.Equal(t, mockMessage, ruleValidationResultWithMessage.Message())
}
