package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
	"time"
)

// AssertEqual checks if values are equal.
func AssertEqual(t *testing.T, a interface{}, expected interface{}) {
	if a == expected {
		return
	}

	// Get the filename + line of where the assertion failed.
	_, filename, line, _ := runtime.Caller(1)
	fmt.Printf("%s:%d expected %v (type %v), got %v (type %v)\n", filepath.Base(filename), line, expected, reflect.TypeOf(expected), a, reflect.TypeOf(a))
	t.FailNow()
}

// IsValidDate checks if a given date string is a valid date.
func IsValidDate(dateString string) bool {
	if dateString == "" {
		return false
	}

	_, err := time.Parse("2006-01-02", dateString)
	return err == nil
}

// MustParseDate accepts a date string and returns a time.Time value.
func MustParseDate(dateString string) time.Time {
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		log.Fatalf("Could not convert %v to a date.\n", dateString)
	}
	return date
}
