// Package utils is collections of utility functions for the Whalelint project.
// It mostly involves slice versions of functions found in the strings package and some other helper ones.
package utils

import (
	"errors"
	"regexp"
	"strings"
)

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

// RemoveExtraSpaces removes all extra consecutive spaces from a string.
func RemoveExtraSpaces(str string) string {
	space := regexp.MustCompile(`\s+`)

	return space.ReplaceAllString(str, " ")
}

// SplitMulti splits a string on any delimiter string found in the delimiterList.
func SplitMulti(str string, delimiterList []string) []string {
	indiceList := FindIndexOfSubString(str, delimiterList)
	if len(indiceList) == 0 {
		return []string{str}
	}

	result := make([]string, len(indiceList))

	prev := 0

	for i, indice := range indiceList {
		result[i] = str[prev:indice]
		prev = i + 1
	}

	return result
}

// FindIndexOfSubString finds the index of any of the patterns in patternList inside the str string.
func FindIndexOfSubString(str string, patternList []string) []int {
	var indexList []int

	for _, pattern := range patternList {
		if idx := strings.Index(str, pattern); idx != -1 {
			indexList = append(indexList, idx)
		}
	}

	return indexList
}

// FindIndexOfSliceElement returns the index of the slice element that matched with the pattern.
// Pattern can be a string or a string slice. In the latter case the first matched index will be returned.
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
