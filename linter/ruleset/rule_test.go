package ruleset_test

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestSeverity_String(t *testing.T) {
	t.Parallel()

	severitySlice := []RuleSet.Severity{
		RuleSet.ValDeprecation, RuleSet.ValError, RuleSet.ValInfo, RuleSet.ValUnknown, RuleSet.ValWarning,
	}
	severityStringSlice := []string{"Deprecation", "Error", "Info", "Unknown", "Warning"}
	invalidSeverity := RuleSet.Severity(math.MaxUint32)

	assert.Equal(t, len(severitySlice), len(severityStringSlice))

	for i, severity := range severitySlice {
		assert.Equal(t, severityStringSlice[i], severity.String())
	}

	assert.Equal(t, "Unknown", invalidSeverity.String())
}

func TestSeverity_MarshalJSON(t *testing.T) {
	t.Parallel()

	severitySlice := []RuleSet.Severity{
		RuleSet.ValDeprecation, RuleSet.ValError, RuleSet.ValInfo, RuleSet.ValUnknown, RuleSet.ValWarning,
	}

	for _, severity := range severitySlice {
		var severityUnmarshalInto RuleSet.Severity

		severityJSON, errMarshal := json.Marshal(severity)
		if errMarshal != nil {
			t.Error(errMarshal)
		}

		errUnMarshal := json.Unmarshal(severityJSON, &severityUnmarshalInto)
		if errUnMarshal != nil {
			t.Error(errUnMarshal)
		}

		assert.Equal(t, severity, severityUnmarshalInto)
	}

	var severityUnmarshalInto RuleSet.Severity
	err := json.Unmarshal([]byte("\"RandomSeverity\""), &severityUnmarshalInto)
	assert.NotEqual(t, nil, err)
}

func TestRule_Validate(t *testing.T) {
	t.Parallel()

	type MockRule struct {
		called int
	}

	mockFunc := func(rule *MockRule) RuleSet.RuleValidationResult {
		rule.called++
		return RuleSet.RuleValidationResult{} // nolint: nlreturn
	}

	a := &MockRule{called: 0}
	rule := RuleSet.NewRule("MockID", "Mock", "MockDesc", RuleSet.ValInfo, mockFunc)

	rule.Validate(a)
	assert.Equal(t, 1, a.called)
}

func TestRule_ValidationFunc(t *testing.T) {
	t.Parallel()

	mockFunc := func(int) {}
	mockRule := RuleSet.NewRule("MockID", "Mock", "MockDesc", RuleSet.ValInfo, mockFunc)

	assert.ObjectsAreEqual(mockFunc, mockRule.ValidationFunc())
}

func TestRuleMapType_Count(t *testing.T) {
	t.Parallel()

	ruleMap := RuleSet.RuleMapType{}
	mockRule := newMockRule()
	mockRuleSet := []RuleSet.Rule{*mockRule, *mockRule}
	ruleMap["mock"] = mockRuleSet

	assert.Equal(t, len(mockRuleSet), ruleMap.Count())
}

func TestRuleMapType_GetRuleByName(t *testing.T) {
	t.Parallel()

	ruleMap := RuleSet.RuleMapType{}
	mockFunc := func(int) /* RuleSet.RuleValidationResult */ {}
	targetName := "FakeID2"
	mockRule1 := RuleSet.NewRule("FakeID1", "Fake definition 1", "MockDesc 1", RuleSet.ValInfo, mockFunc)
	mockRule2 := RuleSet.NewRule(targetName, "Fake definition 2", "MockDesc 2", RuleSet.ValUnknown, mockFunc)

	mockRuleSet := []RuleSet.Rule{*mockRule1, *mockRule2, *mockRule1}
	ruleMap["mock"] = mockRuleSet

	foundRule := ruleMap.GetRuleByName(targetName, int(1))
	assert.Equal(t, mockRule2.Description(), foundRule.Description())

	assert.Equal(t, RuleSet.Rule{}, ruleMap.GetRuleByName(targetName, float32(1)))
}

func TestRule_DocsReference(t *testing.T) {
	t.Parallel()

	mockFunc := func(int) {}
	ruleCopy := RuleSet.NewRule("CPY000", "", "", RuleSet.ValUnknown, mockFunc)
	ruleNone := RuleSet.NewRule("XXX000", "", "", RuleSet.ValUnknown, mockFunc)

	assert.Equal(t, RuleSet.DocsReferenceMap["CPY"], ruleCopy.DocsReference())
	assert.Equal(t, RuleSet.ToDoReference, ruleNone.DocsReference())
}
