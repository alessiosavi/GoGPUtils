package searchutils

import (
	"github.com/alessiosavi/ahocorasick"
)

// LinearSearch is a simple for delegated to find the target value
func LinearSearch[T comparable](data []T, target T) int {
	var i int
	for i = range data {
		if target == data[i] {
			return i
		}
	}
	return -1
}

// ContainsString is delegated to verify if the given string is present in the data
func ContainsString(data, target string) bool {
	matcher := ahocorasick.NewStringMatcher([]string{target})
	return matcher.Contains([]byte(data))
}

// ContainsStringByte is delegated to verify if the given string is present in the data
func ContainsStringByte(data []byte, target string) bool {
	matcher := ahocorasick.NewStringMatcher([]string{target})
	return matcher.Contains(data)
}

// ContainsStrings is delegated to verify if the given array of string are present in the data
func ContainsStrings(data string, targets []string) bool {
	matcher := ahocorasick.NewStringMatcher(targets)
	return matcher.Contains([]byte(data))
}

// ContainsStringsByte is delegated to verify if the given array of string are present in the data
func ContainsStringsByte(data []byte, targets []string) bool {
	matcher := ahocorasick.NewStringMatcher(targets)
	return matcher.Contains(data)
}

// ContainsWhichStrings is delegated to verify which strings are present in the data
func ContainsWhichStrings(data string, target []string) []string {
	matcher := ahocorasick.NewStringMatcher(target)
	hits := matcher.Match([]byte(data))
	found := make([]string, len(hits))
	for i := range hits {
		found[i] = target[hits[i]]
	}
	return found
}

// ContainsWhichStringsByte is delegated to verify which strings are present in the data
func ContainsWhichStringsByte(data []byte, target []string) []string {
	matcher := ahocorasick.NewStringMatcher(target)
	hits := matcher.Match(data)
	found := make([]string, len(hits))
	for i := range hits {
		found[i] = target[hits[i]]
	}
	return found
}
