package parser_test

import (
	"testing"

	"github.com/google/shlex"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/stretchr/testify/assert"

	Parser "github.com/cremindes/whalelint/parser"
)

func TestBashCommand_OptionKeyList(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		optionMap      map[string]string
		optionKeySlice []string
		name           string
	}{
		{
			optionMap:      map[string]string{"foo": "bar"},
			optionKeySlice: []string{"foo"},
			name:           "Key-value option --foo=bar.",
		},
		{
			optionMap:      map[string]string{"foo": ""},
			optionKeySlice: []string{"foo"},
			name:           "Key-value option --foo.",
		},
		{
			optionMap:      map[string]string{"foo": "bar", "x": ""},
			optionKeySlice: []string{"foo", "x"},
			name:           "Key-value option --foo=bar -x.",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			bashCommand := Parser.NewBashCommand(map[string]string{}, "", "", testCase.optionMap,
				map[string]string{}, false, "")

			assert.Equal(t, testCase.optionKeySlice, bashCommand.OptionKeyList())
		})
	}
}

// TestSplitBashChainLex tests the SplitBashChainLex function's ability to properly split the chained raw lex string
// into a set of bashCommands and a set of delimiters.
//
// Scenario: Raw bash string needs to be parsed into a bashCommandList and a delimiterList.
//
// GIVEN |            raw string            | W | SplitBashChainLex  | T | the expectation is
//       | "echo \"ok\""                    | H |     is called      | H | 1 bashCommand, nil delimiterSet
//       | "echo \"ok\" && echo --version"  | E |                    | E | 2 bashCommand,  1  delimiter in set
//       | "echo \"ok\" ; || date && date"  | N |                    | N | 4 bashCommand,  3  delimiter in set
//
func TestSplitBashChainLex(t *testing.T) {
	t.Parallel()

	// nolint:gofmt,gofumpt,goimports
	testCases := []struct {
		name          string
		rawStr        string
		length        int
		delimiterList []string
	}{
		{"simple command" , "echo \"ok\""                  , 1, nil},
		{"two commands"   , "echo \"ok\" && echo --version", 2, []string{"&&"}},
		{"complex command", "echo \"ok\" ; || date && date", 4, []string{";", "||", "&&"}},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			lex, err := shlex.Split(testCase.rawStr)
			if err != nil {
				t.Error("Failed to run lexer on raw string.")
			}
			bashCommandList, delimiterList := Parser.SplitBashChainLex(lex)

			assert.Equal(t, len(bashCommandList), testCase.length)
			assert.EqualValues(t, delimiterList, testCase.delimiterList)
		})
	}
}

func TestParseBashCommand(t *testing.T) { // nolint:funlen
	t.Parallel()

	testCases := []struct {
		name           string
		bashCommandLex []string
		bashCommand    Parser.BashCommand
	}{
		{
			"simple command",
			[]string{"echo", "ok"},
			Parser.NewBashCommand(
				map[string]string{},
				"echo",
				"",
				map[string]string{},
				map[string]string{"ok": ""},
				false,
				"echo ok", // shlex removes quotes - TODO: reconsider this part, with locationParsing in mind!
			),
		},
		{
			"complex command",
			[]string{"envVar=envVar", "sudo", "apt-get", "install", "--yes", "vim=1.2.3", "test"},
			Parser.NewBashCommand(
				map[string]string{"envVar": "envVar"},
				"apt-get",
				"install",
				map[string]string{"--yes": ""},
				map[string]string{"vim": "1.2.3", "test": ""},
				true,
				"envVar=envVar sudo apt-get install --yes vim=1.2.3 test",
			),
		},
		{
			"empty command",
			[]string{},
			Parser.NewBashCommand(
				map[string]string{},
				"",
				"",
				map[string]string{},
				map[string]string{},
				false,
				"",
			),
		},
		{
			"pip install pytorch command",
			[]string{"pip", "install", "pytorch"},
			Parser.NewBashCommand(
				map[string]string{},
				"pip",
				"install",
				map[string]string{},
				map[string]string{"pytorch": ""},
				false,
				"pip install pytorch",
			),
		},
		{
			"pip install pytorch command",
			[]string{"pip", "install", "pytorch", "pytorch-lightning"},
			Parser.NewBashCommand(
				map[string]string{},
				"pip",
				"install",
				map[string]string{},
				map[string]string{"pytorch": "", "pytorch-lightning": ""},
				false,
				"pip install pytorch pytorch-lightning",
			),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			bashCommand := Parser.ParseBashCommand(testCase.bashCommandLex)
			assert.EqualValues(t, bashCommand, testCase.bashCommand)

			// Test getters also
			assert.Equal(t, bashCommand.EnvVars(), testCase.bashCommand.EnvVars())
			assert.Equal(t, bashCommand.Bin(), testCase.bashCommand.Bin())
			assert.Equal(t, bashCommand.SubCommand(), testCase.bashCommand.SubCommand())
			assert.Equal(t, bashCommand.ArgMap(), testCase.bashCommand.ArgMap())
			assert.Equal(t, bashCommand.HasSudo(), testCase.bashCommand.HasSudo())
			assert.Equal(t, bashCommand.OptionList(), testCase.bashCommand.OptionList())
			assert.Equal(t, bashCommand.String(), testCase.bashCommand.String())
		})
	}
}

func TestParseBashCommandChain(t *testing.T) { // nolint:funlen
	t.Parallel()

	testCases := []struct {
		name             string                  // test name
		input            interface{}             // input
		bashCommandChain Parser.BashCommandChain // expected output
	}{
		{
			"Parse basic string into a bash command chain.",
			"echo 1",
			Parser.BashCommandChain{
				BashCommandList: []Parser.BashCommand{Parser.NewBashCommand(
					map[string]string{},
					"echo",
					"",
					map[string]string{},
					map[string]string{"1": ""},
					false,
					"echo 1",
				)},
				OperatorList: nil,
			},
		},
		{
			"Parse *instrunctions.RunCommand into a bash command chain.",
			&instructions.RunCommand{
				ShellDependantCmdLine: instructions.ShellDependantCmdLine{
					CmdLine:      []string{"echo 1"},
					PrependShell: true,
				},
			},
			Parser.BashCommandChain{
				BashCommandList: []Parser.BashCommand{Parser.NewBashCommand(
					map[string]string{},
					"echo",
					"",
					map[string]string{},
					map[string]string{"1": ""},
					false,
					"echo 1",
				)},
				OperatorList: nil,
			},
		},
		{
			"Parse empty *instrunctions.RunCommand into a bash command chain.",
			&instructions.RunCommand{
				ShellDependantCmdLine: instructions.ShellDependantCmdLine{
					CmdLine:      []string{},
					PrependShell: true,
				},
			},
			Parser.BashCommandChain{
				BashCommandList: []Parser.BashCommand{Parser.NewBashCommand(
					map[string]string{},
					"",
					"",
					map[string]string{},
					map[string]string{},
					false,
					"",
				)},
				OperatorList: nil,
			},
		},
		{
			"Parse empty string array into a bash command chain.",
			[]string{},
			Parser.BashCommandChain{
				BashCommandList: []Parser.BashCommand{Parser.NewBashCommand(
					map[string]string{},
					"",
					"",
					map[string]string{},
					map[string]string{},
					false,
					"",
				)},
				OperatorList: nil,
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			bashCommandChain := Parser.ParseBashCommandChain(testCase.input)
			assert.EqualValues(t, bashCommandChain, testCase.bashCommandChain)
		})
	}
}
