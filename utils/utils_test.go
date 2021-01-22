// nolint:funlen
package utils_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	Utils "github.com/cremindes/whalelint/utils"
)

// TestSplitBashChainLex tests the EqualsEither ability to find correctly whether str is part of strSlice.
//
// Scenario: EqualsEither is called with inStr and inStrSlice as inputs and it returns a boolean according to docs.
//
// G | inStr | and |      inStrSlice       | W | EqualsEither | T | the expected return value is
// I |  Foo  |     | ["Foo"]               | H |  is called   | H |  true
// V |  Foo  |     | ["FooWithSuffix"]     | E |              | E | false
// E |  Foo  |     | ["prefixedFoo"]       | N |              | N | false
// N |  Foo  |     | ["foo"]               |   |              |   | false
//   |  Foo  |     | ["bar", "Foo", "Bar"] |   |              |   |  true
//
func TestEqualsEither(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		inStr      string
		inStrSlice []string
		expected   bool
	}{
		{
			name:       "String equals with first and only element in string slice.",
			inStr:      "Foo",
			inStrSlice: []string{"Foo"},
			expected:   true,
		},
		{
			name:       "String does not equal with string+suffix as single element in string slice.",
			inStr:      "Foo",
			inStrSlice: []string{"FooWithSuffix"},
			expected:   false,
		},
		{
			name:       "String does not equal with prefix+string as single element in string slice.",
			inStr:      "Foo",
			inStrSlice: []string{"prefixedFoo"},
			expected:   false,
		},
		{
			name:       "String does not equal with lowercase string as single element in string.",
			inStr:      "Foo",
			inStrSlice: []string{"foo"},
			expected:   false,
		},
		{
			name:       "String equals with second out of free elements in string slice.",
			inStr:      "Foo",
			inStrSlice: []string{"bar", "Foo", "Bar"},
			expected:   true,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, Utils.EqualsEither(testCase.inStr, testCase.inStrSlice), testCase.expected)
		})
	}
}

// G |           inStr           | W | RemoveExtraSpaces | trim  | T | the expected return value is
// I |  "Test with 3 spaces."    | H |     is called     | false | H | "Test with 3 spaces" - no change
// V |  ""                       | E |       with        | false | E | "" - no change
// E |  "   "                    | N |                   | false | N | " "
// N |  "   "                    | N |                   |  true |   | ""
//   |  "  a  "                  |   |                   | false |   | " a "
//   |  "  a  "                  |   |                   |  true |   | "a"
//   |  " Foo  bar foo    bar  " |   |                   | false |   | " Foo bar foo bar "
//   |  " Foo  bar foo    bar  " |   |                   |  true |   | "Foo bar foo bar"
func TestRemoveExtraSpaces(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		inStr    string
		inTrim   bool
		expected string
	}{
		{
			name:     "No change in no extra spaces situation.",
			inStr:    "Test with 3 spaces.",
			inTrim:   false,
			expected: "Test with 3 spaces.",
		},
		{
			name:     "No change in no extra spaces situation.",
			inStr:    "",
			inTrim:   false,
			expected: "",
		},
		{
			name:     "Reduce 3 string of 3 spaces to 1.",
			inStr:    "   ",
			inTrim:   false,
			expected: " ",
		},
		{
			name:     "Reduce 3 string of 3 spaces to 0 when param trim is true.",
			inStr:    "   ",
			inTrim:   true,
			expected: "",
		},
		{
			name:     "Remove extra spaces around single character, but 1-1 before and after it.",
			inStr:    "  a  ",
			inTrim:   false,
			expected: " a ",
		},
		{
			name:     "Remove extra spaces and trim around single character.",
			inStr:    "  a  ",
			inTrim:   true,
			expected: "a",
		},
		{
			name:     "Remove extra spaces complex test without trimming.",
			inStr:    " Foo  bar foo    bar  ",
			inTrim:   false,
			expected: " Foo bar foo bar ",
		},
		{
			name:     "Remove extra spaces complex test with trimming.",
			inStr:    " Foo  bar foo    bar  ",
			inTrim:   true,
			expected: "Foo bar foo bar",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, Utils.RemoveExtraSpaces(testCase.inStr, testCase.inTrim), testCase.expected)
		})
	}
}

func TestSplitMulti(t *testing.T) {
	t.Parallel()

	t.SkipNow() // TODO: enable, once SplitMulti is implemented

	testCases := []struct {
		name           string
		inStr          string
		inPatternSlice []string
		expected       []string
	}{
		{
			name:           "Do not split string with no delimiter rune '.' present in inStr.",
			inStr:          "foo bar",
			inPatternSlice: []string{"."},
			expected:       []string{"foo bar"},
		},
		{
			name:           "Split string into two around delimiter rune '.'.",
			inStr:          "foo.bar",
			inPatternSlice: []string{"."},
			expected:       []string{"foo", "bar"},
		},
		{
			name:           "Split string into three around delimiter rune '.'.",
			inStr:          ".foo.bar",
			inPatternSlice: []string{"."},
			expected:       []string{"", "foo", "bar"},
		},
		{
			name:           "Split string into three around delimiter runes '.' and '|'.",
			inStr:          "foo|bar.foo2",
			inPatternSlice: []string{".", "|"},
			expected:       []string{"foo", "bar", "foo2"},
		},
		{
			name:           "Split string only into two around delimiter rune '|' but not '.'.",
			inStr:          "foo|bar.foo2",
			inPatternSlice: []string{"&", "|"},
			expected:       []string{"foo", "bar.foo2"},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			t.Log(Utils.SplitMulti(testCase.inStr, testCase.inPatternSlice), testCase.expected)

			assert.ElementsMatch(t, Utils.SplitMulti(testCase.inStr, testCase.inPatternSlice), testCase.expected)
		})
	}
}

// G |         inSlice               | inPatterns     | W | FindIndexOfSliceElement | the expected return value is
// I |  ["foo"]                      | ["foo"]        | H |        is called        |  0
// V |  ["foo"]                      | ["bar"]        | E |                         | -1
// E |  ["bar", "foo", "bar", "foo"] | ["foo"]        | N |                         |  1
// N |  ["foo1", "foo", "bar"]       | ["bar", "foo"] |   |                         |  1
//   |  ["foo1", "foo", "bar"]       | []             |   |                         | -1
//   |  ["foo1", "foo", "bar"]       | [""]           |   |                         | -1
func TestFindIndexOfSliceElement(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		inSlice   []string
		inPattern interface{}
		expected  int
	}{
		{
			name:      "Index of matching [str] with str.",
			inSlice:   []string{"foo"},
			inPattern: "foo",
			expected:  0,
		},
		{
			name:      "Index of not matching [str] with other str.",
			inSlice:   []string{"foo"},
			inPattern: "bar",
			expected:  -1,
		},
		{
			name:      "Index of first match.",
			inSlice:   []string{"bar", "foo", "bar", "foo"},
			inPattern: "foo",
			expected:  1,
		},
		{
			name:      "Index of first match in inSlice with any of inPattern",
			inSlice:   []string{"foo1", "foo", "bar"},
			inPattern: []string{"bar", "foo"},
			expected:  1,
		},
		{
			name:      "Nil search pattern.",
			inSlice:   []string{"foo1", "foo", "bar"},
			inPattern: []string{},
			expected:  -1,
		},
		{
			name:      "Empty string search pattern.",
			inSlice:   []string{"foo1", "foo", "bar"},
			inPattern: []string{""},
			expected:  -1,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, Utils.FindIndexOfSliceElement(testCase.inSlice, testCase.inPattern), testCase.expected)
		})
	}
}

func TestInsertIntoSlice(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		inSlice   []string
		inElement string
		inIndex   int
		expectedV []string
		expectedE error
	}{
		{
			name:      "Insert element in front of slice.",
			inSlice:   []string{"foo", "bar"},
			inElement: "item",
			inIndex:   0,
			expectedV: []string{"item", "foo", "bar"},
			expectedE: nil,
		},
		{
			name:      "Insert element into slice with index 1.",
			inSlice:   []string{"foo", "bar"},
			inElement: "item",
			inIndex:   1,
			expectedV: []string{"foo", "item", "bar"},
			expectedE: nil,
		},
		{
			name:      "Insert element into slice at the end",
			inSlice:   []string{"foo", "bar"},
			inElement: "item",
			inIndex:   2,
			expectedV: []string{"foo", "bar", "item"},
			expectedE: nil,
		},
		{
			name:      "Insert element into slice at invalid position.",
			inSlice:   []string{"foo", "bar"},
			inElement: "item",
			inIndex:   3,
			expectedV: []string{"foo", "bar"},
			expectedE: Utils.ErrOutOfBounds,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			newSlice, err := Utils.InsertIntoSlice(testCase.inSlice, testCase.inElement, testCase.inIndex)

			assert.Equal(t, newSlice, testCase.expectedV)
			assert.Equal(t,      err, testCase.expectedE)  // nolint:gofmt,gofumpt,goimports
		})
	}
}

func TestSplitKeyValue(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		inStr         string
		inRune        rune
		expectedKey   string
		expectedValue string
	}{
		{
			name:          "Split simple key=value into key and value.",
			inStr:         "key=value",
			inRune:        '=',
			expectedKey:   "key",
			expectedValue: "value",
		},
		{
			name:          "Tyy to split key=value into key and value by '|' rune.",
			inStr:         "key=value",
			inRune:        '|',
			expectedKey:   "key=value",
			expectedValue: "",
		},
		{
			name:          "Tyy to split key=value into key and value by '=' rune.",
			inStr:         "key=value=1.2",
			inRune:        '=',
			expectedKey:   "key",
			expectedValue: "value=1.2",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			key, value := Utils.SplitKeyValue(testCase.inStr, testCase.inRune)

			assert.Equal(t,   key, testCase.expectedKey  ) // nolint:gofmt,gofumpt,goimports
			assert.Equal(t, value, testCase.expectedValue)
		})
	}
}

func TestParseKeyValueMap(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		inStrSlice     []string
		inRune         rune
		inFinishOnMiss bool
		expected       map[string]string
	}{
		{
			name:           "Split simple, unique key=value x 2 list.",
			inStrSlice:     []string{"key1=value1", "key2=value2"},
			inRune:         '=',
			inFinishOnMiss: false,
			expected:       map[string]string{"key1": "value1", "key2": "value2"},
		},
		{
			name:           "Split duplicated key=value list.",
			inStrSlice:     []string{"key1=value1", "key1=value2"},
			inRune:         '=',
			inFinishOnMiss: false,
			expected:       map[string]string{"key1": "value2"},
		},
		{
			name:           "Split key=value list with irrelevant end.",
			inStrSlice:     []string{"key1=value1", "key2=value2", "arg1", "key2=value3"},
			inRune:         '=',
			inFinishOnMiss: true,
			expected:       map[string]string{"key1": "value1", "key2": "value2"},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			assert.EqualValues(t,
				Utils.ParseKeyValueMap(testCase.inStrSlice, testCase.inRune, testCase.inFinishOnMiss),
				testCase.expected)
		})
	}
}
