// Package parser package help parse various strings into meaningful
// - bash commands or chains of bash commands.
package parser

import (
	"strings"

	"github.com/google/shlex"
	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	log "github.com/sirupsen/logrus"

	Utils "github.com/cremindes/whalelint/utils"
)

// BashCommandChain represents a chain of bash commands.
// e.g. "date; mkdir /test && ls test || echo "test"
// BashCommandList holds the individual bash commands, while OperatorList the operators - "&", "&&", "|", "||", ";".
type BashCommandChain struct {
	BashCommandList []BashCommand
	OperatorList    []string
}

// BashCommand represents a single bash command in a semantically mostly parsed form.
// It includes
// - any environment variables defined before, the
// - binary,
// - sub-command optionally - only for selected binaries for which a rule exits, like apt-get,
// - options, e.g. --yes,
// - rest of the argument list,
// - raw string of the bash command,
// - sudo modifier.
type BashCommand struct {
	envVars    map[string]string
	bin        string
	subCommand string
	optionMap  map[string]string
	argMap     map[string]string
	hasSudo    bool
	rawString  string
}

// EnvVars returns the environment variables defined before the binary on the same line.
func (bashCommand *BashCommand) EnvVars() map[string]string {
	return bashCommand.envVars
}

// Bin returns the binary of the command.
func (bashCommand *BashCommand) Bin() string {
	return bashCommand.bin
}

// SubCommand returns the optional subcommand for the selected set of pre-known binaries.
func (bashCommand *BashCommand) SubCommand() string {
	return bashCommand.subCommand
}

// ArgMap returns the argument list of the bash command.
func (bashCommand *BashCommand) ArgMap() map[string]string {
	return bashCommand.argMap
}

// OptionList returns the options passed to the binary and/or the subcommand, like --yes.
func (bashCommand *BashCommand) OptionList() map[string]string {
	return bashCommand.optionMap
}

// HasSudo tells whether the command has a sudo modifier in front of it.
func (bashCommand *BashCommand) HasSudo() bool {
	return bashCommand.hasSudo
}

// String returns the raw string of the bash command, that served as the basis of the parsing.
func (bashCommand *BashCommand) String() string {
	return bashCommand.rawString
}

// ParseBashCommandList parses a list of bash commands, either from a raw string or buildkit::*instructions.RunCommand.
func ParseBashCommandList(command interface{}) []BashCommand {
	return ParseBashCommandChain(command).BashCommandList
}

// ParseBashCommandChain parses a chain of bash commands separated by bash operators [&, &&, |, ||, ;, etc.].
// Currently it can digest either a raw string or a buildkit::*instructions.RunCommand.
func ParseBashCommandChain(command interface{}) BashCommandChain {
	var (
		err error
		lex []string
	)

	switch c := command.(type) {
	case string:
		lex, err = shlex.Split(c)
	case *instructions.RunCommand:
		lex, err = shlex.Split(c.String()[4:])
	}

	if err != nil && len(lex) == 0 {
		log.Error("Cannot lex bash command.", err)

		return BashCommandChain{}
	}

	bashCommandChain := BashCommandChain{}

	lex = convertSemicolonsToLexItems(lex)

	bashCommandChainLex, delimiterLex := splitBashChainLex(lex)
	bashCommandChain.OperatorList = delimiterLex

	for _, bashCommandLex := range bashCommandChainLex {
		bashCommandChain.BashCommandList = append(bashCommandChain.BashCommandList, ParseBashCommand(bashCommandLex))
	}

	if bashCommandChain.BashCommandList == nil {
		log.Println("fuck")
	}

	return bashCommandChain
}

// ParseBashCommand parses a bash command from a []string format.
// The latter is currently obtained by github.com/google/shlex::Split.
func ParseBashCommand(bashCommandLex []string) BashCommand {
	bashCommand := BashCommand{
		envVars:   make(map[string]string),
		argMap:    make(map[string]string),
		optionMap: make(map[string]string),
	}

	// Not intended as a full list, just mostly what the rules are using, so keep this in mind while debugging!
	subCommandMap := map[string][]string{
		"apt-get": {"clean", "install", "remove", "update", "upgrade", "dist-upgrade"},
		"pip":     {"install", "freeze", "list", "download"},
	}

	if len(bashCommandLex) == 0 {
		return bashCommand
	}

	// TODO: sudo

	// env vars
	for i, lexItem := range bashCommandLex {
		if strings.Contains(lexItem, "=") {
			keyValue := strings.Split(lexItem, "=")
			bashCommand.envVars[keyValue[0]] = keyValue[1]
		} else {
			// only get the first n key-value pairs, as they are the env vars
			bashCommandLex = bashCommandLex[i:]
			break //nolint:nlreturn
		}
	}

	// binary
	bashCommand.bin, bashCommandLex = bashCommandLex[0], bashCommandLex[1:]

	// optional subcommand
	for _, subCommand := range subCommandMap[bashCommand.bin] {
		if bashCommandLex[0] == subCommand {
			bashCommand.subCommand = subCommand
			bashCommandLex = bashCommandLex[1:]
			break // nolint:wsl,nlreturn
		}
	}

	// options, everything that starts with a - or --
	// TODO: option values
	lastOptionIndex := 0
	for i, lexItem := range bashCommandLex { //nolint:wsl
		if strings.HasPrefix(lexItem, "-") {
			bashCommand.optionMap[lexItem] = ""
			lastOptionIndex = i
		}
	}

	if len(bashCommandLex)-1 > lastOptionIndex {
		bashCommandLex = bashCommandLex[lastOptionIndex+1:]
	}

	// args
	bashCommand.argMap = Utils.ParseKeyValueMap(bashCommandLex, '=')

	return bashCommand
}

// splitBashChainLex splits a bash command lex chain on a set of delimiters.
// It returns the list of bash commands lexes in the chain and the delimiters between them.
func splitBashChainLex(strList []string) ([][]string, []string) {
	var (
		bashCommandList [][]string
		delimiterList   []string
	)

	delimiterSet := []string{";", "|", "||", "&", "&&", ">", "<"}

	for len(strList) > 0 {
		delimiterIndex := Utils.FindIndexOfSliceElement(strList, delimiterSet)
		if delimiterIndex == -1 {
			bashCommandList = append(bashCommandList, strList)
			break // nolint:nlreturn
		}

		bashCommandList = append(bashCommandList, strList[0:delimiterIndex])
		delimiterList = append(delimiterList, strList[delimiterIndex])
		strList = strList[delimiterIndex+1:]
	}

	return bashCommandList, delimiterList
}

// convertSemicolonsToLexItems is a helper function to handle the special case, when some of the bash commands are not
// chained together by bash operators, but by a simple semicolon (";") at the end of one of the lex items.
// This function converts those semicolons into their own lex item. So the special case can be handled like the rest.
func convertSemicolonsToLexItems(strList []string) []string {
	result := make([]string, len(strList), cap(strList))
	copy(result, strList)

	var indexList []int

	for i, str := range strList {
		if str[len(str)-1] == ';' {
			indexList = append(indexList, i+1)
			result[i] = result[i][0 : len(result[i])-1]
		}
	}

	var err error
	for i, index := range indexList {
		result, err = Utils.InsertIntoSlice(result, ";", index+i)
		if err != nil {
			log.Error("Cannot copy into slice.")

			return nil
		}
	}

	return result
}
