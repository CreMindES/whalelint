package parser_test

import (
	"io/ioutil"
	"syscall"
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
	"github.com/stretchr/testify/assert"

	Parser "github.com/cremindes/whalelint/parser"
	TestHelper "github.com/cremindes/whalelint/testhelper"
)

func TestRawDockerfileParser_IsInitialized(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		FileContent string
		Expected    bool
	}{
		{
			Name:        "File with 3 lines.",
			FileContent: "first line\nsecond line\nthird line",
			Expected:    true,
		},
		{
			Name:        "Empty file.",
			FileContent: "",
			Expected:    false,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			rdfp := Parser.RawDockerfileParser{}
			rdfp.UpdateRawStr(testCase.FileContent)

			assert.Equal(t, testCase.Expected, rdfp.IsInitialized())
		})
	}
}

func TestRawDockerfileParser_ParseRawLineRange(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		FileContent string
		ParserRange []parser.Range
		Expected    []string
	}{
		{
			Name:        "Get 2nd and 3rd line.",
			FileContent: "first line\nsecond line\nthird line\nforth line\n",
			ParserRange: []parser.Range{
				{
					Start: parser.Position{Line: 2, Character: 0},
					End:   parser.Position{Line: 2, Character: 0},
				},
				{
					Start: parser.Position{Line: 3, Character: 0},
					End:   parser.Position{Line: 3, Character: 0},
				},
			},
			Expected: []string{"second line", "third line"},
		},
		{
			Name:        "Empty file.",
			FileContent: "",
			ParserRange: []parser.Range{{
				Start: parser.Position{Line: 2, Character: 0},
				End:   parser.Position{Line: 2, Character: 0},
			}},
			Expected: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			rdfp := Parser.RawDockerfileParser{}
			rdfp.UpdateRawStr(testCase.FileContent)

			assert.Equal(t, testCase.Expected, rdfp.ParseRawLineRange(testCase.ParserRange))
		})
	}
}

// nolint:funlen
func TestRawDockerfileParser_StringLocation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		FileContent string
		SearchStr   string
		ParserRange []parser.Range
		Expected    [4]int
	}{
		{
			Name:        "Get location of \"second\" from file of 4 lines with window of first 3 lines.",
			FileContent: "first line\nsecond line\nthird line\nforth line\n",
			SearchStr:   "second",
			ParserRange: []parser.Range{
				{
					Start: parser.Position{Line: 1, Character: 0},
					End:   parser.Position{Line: 1, Character: 0},
				},
				{
					Start: parser.Position{Line: 3, Character: 0},
					End:   parser.Position{Line: 3, Character: 0},
				},
			},
			Expected: [4]int{2, 0, 2, 6},
		},
		{
			Name:        "Get location of \"line\" from file of 4 lines within window of 2nd line.",
			FileContent: "first line\nsecond line\nthird line\nforth line\n",
			SearchStr:   "line",
			ParserRange: []parser.Range{
				{
					Start: parser.Position{Line: 2, Character: 0},
					End:   parser.Position{Line: 2, Character: 0},
				},
			},
			Expected: [4]int{2, 7, 2, 11},
		},
		{
			Name:        "Get location of \"forth\" from file of 4 lines with a window of the first 2 lines.",
			FileContent: "first line\nsecond line\nthird line\nforth line\n",
			SearchStr:   "forth",
			ParserRange: []parser.Range{
				{
					Start: parser.Position{Line: 1, Character: 0},
					End:   parser.Position{Line: 1, Character: 0},
				},
				{
					Start: parser.Position{Line: 2, Character: 0},
					End:   parser.Position{Line: 2, Character: 0},
				},
			},
			Expected: [4]int{-1, -1, -1, -1},
		},
		{
			Name:        "Get location of \"forth\" from file of 4 lines without any window.",
			FileContent: "first line\nsecond line\nthird line\nforth line\n",
			SearchStr:   "forth",
			ParserRange: nil,
			Expected:    [4]int{4, 0, 4, 5},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			rdfp := Parser.RawDockerfileParser{}
			rdfp.UpdateRawStr(testCase.FileContent)

			assert.Equal(t, testCase.Expected, rdfp.StringLocation(testCase.SearchStr, testCase.ParserRange))
		})
	}
}

func TestRawDockerfileParser_ParseDockerfile(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name        string
		FileContent string
		FileExists  bool
		Expected    error
	}{
		{
			Name:        "Parse simple Dockerfile",
			FileContent: "FROM golang:1.16",
			FileExists:  true,
			Expected:    nil,
		},
		{
			Name:        "Parse empty Dockerfile",
			FileContent: "",
			FileExists:  true,
			Expected:    nil,
		},
		{
			Name:        "Parse non-existing file",
			FileContent: "",
			FileExists:  false,
			Expected:    syscall.ENOENT,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.Name, func(t *testing.T) {
			t.Parallel()

			rdfp := Parser.RawDockerfileParser{}
			filePath := ""

			if testCase.FileExists {
				tmpFile, errCreate := ioutil.TempFile("", "test_rawdf.*")
				assert.Nil(t, errCreate)
				_, errWrite := tmpFile.WriteString(testCase.FileContent)
				assert.Nil(t, errWrite)
				errSync := tmpFile.Sync()
				assert.Nil(t, errSync)

				filePath = tmpFile.Name()
			} else {
				filePath = "bogus/File/Path"
			}

			err := rdfp.ParseDockerfile(filePath)

			result := TestHelper.CheckForErrorRecursively(t, err, testCase.Expected)

			assert.Equal(t, true, result)
		})
	}
}
