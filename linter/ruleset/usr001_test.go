package ruleset_test

import (
	"testing"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	RuleSet "github.com/cremindes/whalelint/linter/ruleset"
)

func TestValidateUsr001(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		userList    []string
		isViolation bool
		name        string
	}{
		{userList: []string{"foo"},                isViolation: false, name: "One user foo."},
		{userList: []string{"root"},               isViolation:  true, name: "One user root."},
		{userList: []string{"foo", "bar"},         isViolation: false, name: "Two users, foo then bar."},
		{userList: []string{"root", "bar"},        isViolation: false, name: "Two users, root then bar."},
		{userList: []string{"foo", "root"},        isViolation:  true, name: "Two users, foo then root."},
		{userList: []string{"foo", "root", "bar"}, isViolation: false, name: "Three users, foo, root and then bar."},
		{userList: []string{""},                   isViolation: false, name: "No user."},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			stageList := make([]instructions.Stage, 0, len(testCase.userList))
			for _, userName := range testCase.userList {
				// nolint:exhaustivestruct
				stageList = append(stageList, instructions.Stage{
					Commands: []instructions.Command{
						&instructions.UserCommand{
							User: userName,
						},
					},
				})
			}

			assert.Equal(t, testCase.isViolation, RuleSet.ValidateUsr001(stageList).IsViolated())
		})
	}
}
