package ruleset

import (
	"encoding/json"
	"fmt"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"

	Parser "github.com/cremindes/whalelint/parser"
)

type LocationRange struct {
	start *Location
	end   *Location
}

type Location struct {
	lineNumber int
	charNumber int
}

func (location *Location) LineNumber() int {
	return location.lineNumber
}

func (location *Location) CharNumber() int {
	return location.charNumber
}

func (location Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		LineNumber int
		CharNumber int
	}{
		LineNumber: location.lineNumber,
		CharNumber: location.charNumber,
	})
}

func (location *Location) UnmarshalJSON(data []byte) error {
	loc := struct {
		LineNumber int
		CharNumber int
	}{}

	err := json.Unmarshal(data, &loc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal Location: %w", err)
	}

	location.lineNumber = loc.LineNumber
	location.charNumber = loc.CharNumber

	return nil
}

func (locationRange *LocationRange) Start() *Location {
	return locationRange.start
}

func (locationRange *LocationRange) End() *Location {
	return locationRange.end
}

func (locationRange LocationRange) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Start Location
		End   Location
	}{
		Start: *locationRange.start,
		End:   *locationRange.end,
	})
}

func (locationRange *LocationRange) UnmarshalJSON(data []byte) error {
	lr := struct {
		Start Location
		End   Location
	}{}

	err := json.Unmarshal(data, &lr)
	if err != nil {
		return fmt.Errorf("failed to unmarshal LocationRange: %w", err)
	}

	locationRange.start = &lr.Start
	locationRange.end   = &lr.End // nolint:gofmt,gofumpt,goimports

	return nil
}

func LocationRangeFromCommand(command instructions.Command) LocationRange {
	return CopyLocationRange(command.Location())
}

func CopyLocationRange(parserRange []parser.Range) LocationRange {
	if parserRange == nil {
		return LocationRange{
			start: &Location{
				lineNumber: 1,
				charNumber: 0,
			},
			end: &Location{
				lineNumber: 1,
				charNumber: 0,
			},
		}
	}

	location := LocationRange{
		start: &Location{
			lineNumber: parserRange[0].Start.Line,
			charNumber: parserRange[0].Start.Character,
		},
		end: &Location{
			lineNumber: parserRange[len(parserRange)-1].End.Line,
			charNumber: parserRange[len(parserRange)-1].End.Character,
		},
	}

	return location
}

func NewLocationRange(startLine, startChar, endLine, endChar int) LocationRange {
	return LocationRange{
		start: &Location{
			lineNumber: startLine,
			charNumber: startChar,
		},
		end: &Location{
			lineNumber: endLine,
			charNumber: endChar,
		},
	}
}

func ParseLocationFromRawParser(str string, window []parser.Range) LocationRange {
	if !Parser.RawParser.IsInitialized() {
		return CopyLocationRange(window)
	}

	return NewLocationFrom4Int(
		Parser.RawParser.StringLocation(str, window),
	)
}

func NewLocationFrom4Int(locationRange [4]int) LocationRange {
	return LocationRange{
		start: &Location{locationRange[0], locationRange[1]},
		end:   &Location{locationRange[2], locationRange[3]},
	}
}
