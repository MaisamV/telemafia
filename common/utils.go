package common

import (
	"errors"
	"strconv"
)

// Contains checks if a string is in a slice
func Contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// StringToInt64 converts a digit string to an int64
func StringToInt64(s string) (int64, error) {
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i, nil
	} else {
		return 0, errors.New("invalid digit string")
	}
}
