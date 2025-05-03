package common

import (
	"crypto/sha256"
	"errors"
	"math/rand"
	"strconv"
	"time"
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

// Hash returns the SHA-256 hash of a string
func Hash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return string(h.Sum(nil))
}

var r *rand.Rand

func InitSeed() {
	// Seed random number generator
	source := rand.NewSource(time.Now().UnixNano())
	r = rand.New(source)
}

func Shuffle(n int, swap func(i, j int)) {
	r.Shuffle(n, swap)
}
