// Package utils is collections of utility functions for the Whalelint project.
// It mostly involves slice versions of functions found in the strings package and some other helper ones.
package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
	log "github.com/sirupsen/logrus"
)

/* Errors. */

var ErrUnSupportedType = errors.New("unsupported type")

/* String helper functions. */

// EqualsEither returns true, if str string is a match to any element of the targetList, false otherwise.
func EqualsEither(str string, targetList []string) bool {
	for _, target := range targetList {
		if str == target {
			return true
		}
	}

	return false
}

// SplitMulti splits str string on any pattern match from patternSlice.
// Note: not-implemented!
func SplitMulti(str string, patternSlice []string) []string {
	return []string{}
}

// RemoveExtraSpaces removes all extra consecutive spaces from a string.
func RemoveExtraSpaces(str string, trim bool) string {
	space := regexp.MustCompile(`\s+`)

	result := space.ReplaceAllString(str, " ")
	if trim {
		result = strings.Trim(result, " ")
	}

	return result
}

func SliceContains(arr []string, patternInterface interface{}) bool {
	return FindIndexOfSliceElement(arr, patternInterface) != -1
}

// FindIndexOfSliceElement returns the fist index of the slice element that matched with the pattern.
// Pattern can be a string or a string slice. In the latter case the index of the match of the first element with any
// matching pattern will be returned.
func FindIndexOfSliceElement(arr []string, patternInterface interface{}) int {
	switch pattern := patternInterface.(type) {
	case string:
		for i, item := range arr {
			if item == pattern {
				return i
			}
		}
	case []string:
		for i, item := range arr {
			for _, patternItem := range pattern {
				if item == patternItem {
					return i
				}
			}
		}
	}

	return -1
}

// ErrOutOfBounds is an error when indexing out of a slice.
var ErrOutOfBounds = errors.New("out of bounds")

/* String slice functions. */

// InsertIntoSlice inserts element at index into originalSlice.
func InsertIntoSlice(originalSlice []string, element string, index int) ([]string, error) {
	if index == len(originalSlice) {
		return append(originalSlice, element), nil
	}

	if index > len(originalSlice) {
		return originalSlice, ErrOutOfBounds
	}

	result := append(originalSlice[:index+1], originalSlice[index:]...)
	result[index] = element

	return result, nil
}

/* String map functions. */

// ParseKeyValueMap parses a string list into map[string]string based on the separator rune.
// If finishOnMiss is true, the parser terminates on the first element, where there is no separator rune present
// and returns with the map built so far without the last, unsplittable element.
func ParseKeyValueMap(strList []string, separator rune, finishOnMiss bool) map[string]string {
	resultMap := make(map[string]string)

	for _, item := range strList {
		key, value := SplitKeyValue(item, separator)

		if finishOnMiss && len(value) == 0 {
			break
		}

		resultMap[key] = value
	}

	return resultMap
}

// SplitKeyValue splits s string into two strings on the r rune.
// It's ideal for building maps, as there is always two element, a key and a value returned.
// Naturally, if there is no r rune present in s string, the second element will be an empty string.
func SplitKeyValue(s string, r rune) (string, string) {
	idx := strings.IndexRune(s, r)

	if idx == -1 {
		return s, ""
	}

	return s[0:idx], s[idx+1:]
}

func FilterMapByKey(kv map[string]string, f func(string) bool) []string {
	return filterMap(kv, f, true, false)
}

func FilterMapByValue(kv map[string]string, f func(string) bool) []string {
	return filterMap(kv, f, false, true)
}

func FilterMapKeys(kv map[string]string, f func(string) bool) []string {
	return filterMap(kv, f, true, true)
}

func FilterMapValues(kv map[string]string, f func(string) bool) []string {
	return filterMap(kv, f, false, false)
}

func filterMap(kv map[string]string, f func(string) bool, byKeyOrValue bool, keyOrValue bool) []string {
	result := make([]string, 0)

	filterOut := ""
	sortBy := ""

	for k, v := range kv {
		if byKeyOrValue {
			sortBy = k
		} else {
			sortBy = v
		}

		if keyOrValue {
			filterOut = k
		} else {
			filterOut = v
		}

		if f(sortBy) {
			result = append(result, filterOut)
		}
	}

	return result
}

/* Unix stuff. */

func IsUnixPortValid(portParam interface{}) bool {
	portMin := 0
	portMax := 65535

	var (
		err  error
		port int
	)

	switch portAssert := portParam.(type) {
	case int:
		port = portAssert
	case string:
		// convert our port string to integer
		portAssert = strings.TrimSpace(portAssert)
		port, err = strconv.Atoi(portAssert)

		if err != nil {
			log.Debug("Cannot convert port string to int!")
			return false // nolint:nlreturn
		}
	default:
		log.Error("Unsupported portParam type.")
		return false // nolint:nlreturn
	}

	if (portMin <= port) && (port <= portMax) {
		return true
	}

	return false
}

/* File */

func ReadFileContents(filePath string) (string, error) {
	// TODO: migrate to 1.16 os.ReadFile
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return string(data), nil
}

/* Docker related functions. */

// MatchDockerImageNames compares Docker image name strings with the addition of including the fact, that when no tag
// is specified, ":latest" is assumed.
func MatchDockerImageNames(str1, str2 string) bool {
	return strings.TrimSuffix(str1, ":latest") == strings.TrimSuffix(str2, ":latest")
}

func GetDockerfileAst(filePathString string) ([]instructions.Stage, []instructions.ArgCommand, error) {
	filePath := filepath.Clean(filePathString)

	fileHandle, err := os.Open(filePath)
	if err != nil {
		// log.Error("Cannot open Dockerfile \"", filePath, "\".", err)
		return nil, nil, fmt.Errorf("%w", err)
	}
	defer fileHandle.Close()

	dockerfile, err := parser.Parse(fileHandle)
	if err != nil {
		// log.Error("Cannot parse Dockerfile \"", filePath, "\"", err)
		return nil, nil, fmt.Errorf("dockerfile parse | %w", err)
	}

	stageList, metaArgs, err := instructions.Parse(dockerfile.AST)
	if err != nil {
		log.Debug("Cannot create Dockerfile AST from \"", filePath, "\".", err)
		stageList, metaArgs = ParseDockerfileInstructionsSafely(dockerfile, fileHandle) // nolint:wsl
	}

	return stageList, metaArgs, nil
}

// ParseDockerfileInstructionsSafely parses Dockerfile Instructions representation by iteratively trying to correct
// the AST representation - if needed.
//
// If there is an error, it tries to
// - find the offending line from the error message of buildkit::instructions::Parse
// - match that line to a dockerfile.AST.children element
// - remove that child
// - try again until there is
//   - either a valid AST tree that can be parsed further
//   - there is no more child, in which case it returns an empty stage.
// nolint:funlen
func ParseDockerfileInstructionsSafely(dockerfile *parser.Result, fileHandle io.ReadSeeker) ([]instructions.Stage,
	[]instructions.ArgCommand) {
	if dockerfile == nil {
		return []instructions.Stage{}, []instructions.ArgCommand{}
	}

	var (
		stageList []instructions.Stage
		metaArgs  []instructions.ArgCommand
		err       error
	)

	for stageList == nil {
		stageList, metaArgs, err = instructions.Parse(dockerfile.AST)
		// Try to correct, if there is any error
		if err != nil {
			log.Trace("Cannot create Dockerfile AST", err)
			// parse offending node number
			regexpOffendingLine := regexp.MustCompile(" parse error line ([1-9]+[0-9]*):")
			strSlice := regexpOffendingLine.FindStringSubmatch(err.Error())
			offendingLineIndex, _ := strconv.Atoi(strSlice[1])

			_, errSeek := fileHandle.Seek(0, 0) // go back to the beginning of the file
			if errSeek != nil {
				log.Error(errSeek)
				os.Exit(1)
			}

			var offendingLineStr string

			reader := bufio.NewReader(fileHandle)

			for i := 1; i <= offendingLineIndex; i++ {
				offendingLineStr, err = reader.ReadString('\n')
				if err != nil && !errors.Is(err, io.EOF) {
					break
				}
			}
			log.Trace("Found offending line index: ", offendingLineIndex)

			offendingLineStr = strings.TrimSuffix(offendingLineStr, "\n")
			offendingAstIdx := -1

			for i, child := range dockerfile.AST.Children {
				if child.Original == offendingLineStr {
					log.Trace("Matched offending line ", offendingLineStr, " to child ", i, ".")
					offendingAstIdx = i
				}
			}

			if offendingAstIdx > 0 {
				dockerfile.AST.Children = append(
					dockerfile.AST.Children[:offendingAstIdx],
					dockerfile.AST.Children[offendingAstIdx+1:]...)
			}
		}

		if len(dockerfile.AST.Children) == 0 {
			log.Error("Fallback to Dockerfile AST correction was unsuccessful.", err)

			return []instructions.Stage{}, []instructions.ArgCommand{}
		}
	}

	return stageList, metaArgs
}
