package ruleset

import (
	"encoding/json"
	"fmt"
	"sort"

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
	return BKRangeSliceToLocationRange(command.Location())
}

func LocationRangeToBKRange(locationRange LocationRange) parser.Range {
	return parser.Range{
		Start: parser.Position{
			Line:      locationRange.Start().LineNumber(),
			Character: locationRange.Start().CharNumber(),
		},
		End:   parser.Position{
			Line:      locationRange.End().LineNumber(),
			Character: locationRange.End().CharNumber(),
		},
	}
}

func BKRangeSliceToLocationRange(parserRange []parser.Range) LocationRange {
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
		return BKRangeSliceToLocationRange(window)
	}

	location := NewLocationFrom4Int(
		Parser.RawParser.StringLocation(str, window),
	)

	if location.Start().LineNumber() == -1 {
		return BKRangeSliceToLocationRange(window)
	}

	return location
}

func ParseLocationSliceFromRawParser(strSlice []string, window []parser.Range) []LocationRange {
	if !Parser.RawParser.IsInitialized() {
		return []LocationRange{BKRangeSliceToLocationRange(window)}
	}

	location := NewLocationFrom4IntSlice(
		Parser.RawParser.StringSliceLocation(strSlice, window),
	)

	if location[0].Start().LineNumber() == -1 {
		return []LocationRange{BKRangeSliceToLocationRange(window)}
	}

	return location
}

func NewLocationFrom4Int(locationRange [4]int) LocationRange {
	return LocationRange{
		start: &Location{locationRange[0], locationRange[1]},
		end:   &Location{locationRange[2], locationRange[3]},
	}
}

func NewLocationFrom4IntSlice(locationRangeSlice [][4]int) []LocationRange {
	result := make([]LocationRange, len(locationRangeSlice))

	for i, locationRange := range locationRangeSlice {
		result[i] = LocationRange{
			start: &Location{locationRange[0], locationRange[1]},
			end:   &Location{locationRange[2], locationRange[3]},
		}
	}

	return result
}

func UnionOfLocationRanges(locationRangeSlice []LocationRange) LocationRange {
	// TODO deep copy
	SortLocationRanges(locationRangeSlice)

	return LocationRange{
		start: &Location{
			lineNumber: locationRangeSlice[0].start.lineNumber,
			charNumber: locationRangeSlice[0].start.charNumber,
		},
		end: &Location{
			lineNumber: locationRangeSlice[len(locationRangeSlice)-1].end.lineNumber,
			charNumber: locationRangeSlice[len(locationRangeSlice)-1].end.charNumber,
		},
	}
}

func SortLocationRanges(locationRangeSlice []LocationRange) {
	sortFunc := func(i, j int) bool {
		// compare by 1st variable
		if locationRangeSlice[i].start.lineNumber != locationRangeSlice[j].start.lineNumber {
			return locationRangeSlice[i].start.lineNumber < locationRangeSlice[j].start.lineNumber
		}
		// compare by 2nd variable
		if locationRangeSlice[i].start.charNumber != locationRangeSlice[j].start.charNumber {
			return locationRangeSlice[i].start.charNumber < locationRangeSlice[j].start.charNumber
		}
		// compare by 3rd variable
		if locationRangeSlice[i].end.lineNumber != locationRangeSlice[j].end.lineNumber {
			return locationRangeSlice[i].end.lineNumber < locationRangeSlice[j].end.lineNumber
		}
		// compare by 4th variable
		if locationRangeSlice[i].end.charNumber != locationRangeSlice[j].end.charNumber {
			return locationRangeSlice[i].end.charNumber < locationRangeSlice[j].end.charNumber
		}

		// all equal
		return false
	}

	sort.Slice(locationRangeSlice, sortFunc)
}
