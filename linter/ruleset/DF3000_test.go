package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"

	RuleSet "../ruleset"
)

func TestValidateDl3000(t *testing.T) {
	t.Parallel()

	absWorkdirCommand := &instructions.WorkdirCommand{Path: "/go"}

	if RuleSet.ValidateDl3000(absWorkdirCommand).IsViolated() != false {
		t.Errorf("validateDf3000, a.k.a validate WORKDIR is absolute path, should pass for \"/go\"!")
	}

	nonAbsWorkdirCommand1 := &instructions.WorkdirCommand{Path: "./go"}
	if RuleSet.ValidateDl3000(nonAbsWorkdirCommand1).IsViolated() != true {
		t.Errorf("validateDf3000, a.k.a validate WORKDIR is absolute path, should not pass for \"./go\"!")
	}

	nonAbsWorkdirCommand2 := &instructions.WorkdirCommand{Path: "go/src"}
	if RuleSet.ValidateDl3000(nonAbsWorkdirCommand2).IsViolated() != true {
		t.Errorf("validateDf3000, a.k.a validate WORKDIR is absolute path, should not pass for \"go/src\"!")
	}
}
