package ruleset_test

import (
	"encoding/json"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestLocationRangeGetters(t *testing.T) {
	t.Parallel()

	locationRange := RuleSet.NewLocationRange(1, 2, 3, 4)

	assert.Equal(t, 1, locationRange.Start().LineNumber())
	assert.Equal(t, 2, locationRange.Start().CharNumber())
	assert.Equal(t, 3, locationRange.End().LineNumber())
	assert.Equal(t, 4, locationRange.End().CharNumber())
}

func TestLocationRange_MarshalJSON(t *testing.T) {
	t.Parallel()

	locationRange := RuleSet.NewLocationRange(1, 2, 3, 4)
	locationRangeDuplicate := RuleSet.LocationRange{}

	// serialize
	locationRangeJSON, err := json.Marshal(locationRange)
	if err != nil {
		t.Error(err)
	}
	// deserialize
	err = json.Unmarshal(locationRangeJSON, &locationRangeDuplicate)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, locationRange, locationRangeDuplicate)
}

type MockCommand struct {
	location []parser.Range
}

func (mockCommand *MockCommand) Location() []parser.Range {
	return mockCommand.location
}

func (mockCommand *MockCommand) Name() string {
	return ""
}

func TestLocationRangeFromCommand(t *testing.T) {
	t.Parallel()

	referenceLocation := RuleSet.NewLocationRange(1, 2, 3, 4)

	mockCommand := &MockCommand{location: []parser.Range{{
		Start: parser.Position{
			Line:      1,
			Character: 2,
		},
		End: parser.Position{
			Line:      3,
			Character: 4,
		},
	}}}

	location := RuleSet.LocationRangeFromCommand(mockCommand)

	assert.Equal(t, referenceLocation, location)
}
