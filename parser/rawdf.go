package parser

import (
	"fmt"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"

	Utils "github.com/cremindes/whalelint/utils"
)

var RawParser RawDockerfileParser = RawDockerfileParser{rawStr: ""} // nolint:gochecknoglobals

type RawDockerfileParser struct {
	rawStr   string
	rawLines []string
}

func (r *RawDockerfileParser) IsInitialized() bool {
	return len(r.rawStr) > 0 && len(r.rawLines) > 0
}

func (r *RawDockerfileParser) ParseRawLineRange(p []parser.Range) []string {
	if !r.IsInitialized() {
		return nil
	}

	return r.rawLines[p[0].Start.Line-1 : p[len(p)-1].End.Line]
}

func (r *RawDockerfileParser) UpdateRawStr(str string) {
	r.rawStr = str
	r.rawLines = strings.Split(r.rawStr, "\n")
}

func (r *RawDockerfileParser) ParseDockerfile(filePath string) error {
	str, err := Utils.ReadFileContents(filePath)

	if err == nil {
		r.UpdateRawStr(str)

		return nil
	}

	return fmt.Errorf("%w", err)
}

func (r *RawDockerfileParser) StringLocation(str string, window []parser.Range) [4]int {
	windowStart, windowEnd := 0, 0
	windowStartChar := 0

	switch len(window) {
	case 0:
		windowStart, windowEnd = 0, len(r.rawLines)-1
	case 1:
		windowStart, windowEnd = window[0].Start.Line-1, window[0].Start.Line
		windowStartChar = window[0].Start.Character
	default:
		windowStart, windowEnd = window[0].Start.Line-1, window[len(window)-1].End.Line
		windowStartChar = window[0].Start.Character
	}

	searchWindow := r.rawLines[windowStart:windowEnd]

	for i, line := range searchWindow {
		// tmp workaround till bash parser can work together with raw parser
		if i == 0 {
			line = line[windowStartChar:]
		}

		if strings.Contains(line, str) {
			startLine := i
			startChar := strings.Index(line, str) + windowStartChar
			endLine := i
			endChar := strings.Index(line, str) + len(str) + windowStartChar

			return [4]int{
				windowStart + startLine + 1,
				startChar,
				windowStart + endLine + 1,
				endChar,
			}
		}

		if i == 0 {
			windowStartChar = 0
		}
	}

	return [4]int{-1, -1, -1, -1}
}

func (r *RawDockerfileParser) StringSliceLocation(strSlice []string, window []parser.Range) [][4]int {
	if len(strSlice) == 0 {
		w := window[0]

		return [][4]int{{w.Start.Line, w.Start.Character, w.End.Line, w.End.Character}}
	}

	result := make([][4]int, len(strSlice))

	for i, str := range strSlice {
		result[i] = r.StringLocation(str, window)
	}

	return result
}
