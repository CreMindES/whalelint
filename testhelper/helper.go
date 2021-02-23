package testhelper

import (
	"errors"
	"reflect"
	"testing"
)

func CheckForErrorRecursively(t *testing.T, err error, target error) bool {
	t.Helper()

	// first check for unwrapped errors, especially where both are nil
	if err == nil && target == nil {
		return true
	}

	if errors.Is(err, target) {
		return true
	}

	unWrappedErr := err
	targetType := reflect.TypeOf(target)

	for {
		e := errors.Unwrap(unWrappedErr)
		if e != nil {
			unWrappedErr = e

			if reflect.TypeOf(e) == targetType {
				return true
			}
		} else {
			return false
		}
	}
}
