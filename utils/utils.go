// Package utils is collections of utility functions for the Whalelint project.
// It mostly involves slice versions of functions found in the strings package and some other helper ones.
package utils

import (
	"errors"
	"regexp"
	"strings"
)

/* String helper functions. */

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
func ParseKeyValueMap(strList []string, separator rune) map[string]string {
	argMap := make(map[string]string, len(strList))

	for _, item := range strList {
		var key, value string

		equalSignIndex := strings.IndexRune(item, separator)
		if equalSignIndex == -1 {
			key = item
		} else {
			key = item[0:equalSignIndex]
			value = item[equalSignIndex+1:]
		}

		argMap[key] = value
	}

	return argMap
}
