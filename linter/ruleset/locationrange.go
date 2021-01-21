package ruleset

import (
	"encoding/json"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type LocationRange struct {
	start Location
	end   Location
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

func (locationRange *LocationRange) Start() Location {
	return locationRange.start
}

func (locationRange *LocationRange) End() Location {
	return locationRange.end
}

func (locationRange LocationRange) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Start Location
		End   Location
	}{
		Start: locationRange.start,
		End:   locationRange.end,
	})
}

func LocationRangeFromCommand(command instructions.Command) LocationRange {
	return CopyLocationRange(command.Location())
}

func CopyLocationRange(parserRange []parser.Range) LocationRange {
	if parserRange == nil {
		return LocationRange{
			start: Location{
				lineNumber: 1,
				charNumber: 0,
			},
			end: Location{
				lineNumber: 1,
				charNumber: 0,
			},
		}
	}

	location := LocationRange{
		start: Location{
			lineNumber: parserRange[0].Start.Line,
			charNumber: parserRange[0].Start.Character,
		},
		end: Location{
			lineNumber: parserRange[len(parserRange)-1].End.Line,
			charNumber: parserRange[len(parserRange)-1].End.Character,
		},
	}

	return location
}
