package ruleset

import (
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	log "github.com/sirupsen/logrus"
	set "github.com/zoumo/goset"
)

// STL -> Stage List.
var _ = NewRule("STL001", "Stage name alias must be unique.", "", ValError, ValidateStl001)

func ValidateStl001(stageList []instructions.Stage) RuleValidationResult {
	result := RuleValidationResult{isViolated: false, LocationRange: LocationRange{}}
	stageNameSet := set.NewSet()

	for _, stage := range stageList {
		if stageNameSet.Contains(stage.Name) && stage.Name != "" { // found a non-unique build stage alias
			result.SetViolated()
			result.LocationRange = ParseLocationFromRawParser(stage.Name, stage.Location)
		}

		err := stageNameSet.Add(stage.Name)
		if err != nil {
			log.Error(err)
		}
	}

	return result
}
