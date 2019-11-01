package arrayutils

import (
	"strings"
)

//RemoveElement delete the element of the indexes contained in j of the data in input
func RemoveElements(data []string, j []int) []string {
	for i := 0; i < len(j); i++ {
		data[j[i]] = data[len(data)-1] // Copy last element to index i.
		data[len(data)-1] = ""         // Erase last element (write zero value).
		data = data[:len(data)-1]      // Truncate slice.
	}
	return data
}

func RemoveElement(s []string, i int) []string {
	if i < len(s) {
		s[len(s)-1], s[i] = s[i], s[len(s)-1]
		return s[:len(s)-1]
	}
	return s
}

func RemoveElementsTuned(s []string, index []int) []string {
	if len(index) <= len(s) {
		for i := range index {
			s[len(s)-1], s[i] = s[i], s[len(s)-1]
			s = s[:len(s)-1]
			if i+1 < len(index) {
				index[i+1] = index[i+1] - 1
			}
		}
	}
	return s

}

// RemoveWhiteSpaceArray is delegated to iterate every array item and remove the whitespace from the given string
func RemoveWhiteSpaceArray(data []string) []string {
	var toDelete []int
	// Iterate the string in the list
	for i := 0; i < len(data); i++ {
		// Iterate the char in the string
		for _, c := range data[i] {
			if c == 32 { // if whitespace
				toDelete = append(toDelete, i)
			}
		}
	}
	return RemoveElements(data, toDelete)
}

// JoinStrings concatenate every data in the array and return the string content
func JoinStrings(data []string) string {
	var sb strings.Builder
	for i := 0; i < len(data); i++ {
		sb.WriteString(data[i] + " ")
	}
	return sb.String()
}

// RemoveFromByte Remove a given element from a string
// NOTE: Panic in case of index out of bound
func RemoveFromByte(s []byte, i int) []byte {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

// ReverseArray is delegated to return the inverse rappresentation of the array
func ReverseArray(n1 []int) []int {
	var result []int = make([]int, len(n1))
	for i := len(n1) - 1; i >= 0; i-- {
		v := n1[i]
		j := len(n1) - 1 - i
		result[j] = v
	}
	return result
}
