package ruleset

import (
	"encoding/json"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
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
	location := LocationRange{
		start: Location{
			lineNumber: 1,
			charNumber: 0,
		},
		end: Location{
			lineNumber: 1,
			charNumber: 0,
		},
	}

	commandLocation := command.Location()
	if commandLocation != nil {
		location.start.lineNumber = commandLocation[0].Start.Line
		location.start.charNumber = commandLocation[0].Start.Character
		location.end.lineNumber = commandLocation[0].End.Line
		location.end.charNumber = commandLocation[0].End.Character
	}

	return location
}
