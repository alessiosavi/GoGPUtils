package arrayutils

import "strings"

//RemoveElement delete the element of the indexes contained in j of the data in input
func RemoveElement(data []string, j []int) []string {
	for i := 0; i < len(j); i++ {
		data[j[i]] = data[len(data)-1] // Copy last element to index i.
		data[len(data)-1] = ""         // Erase last element (write zero value).
		data = data[:len(data)-1]      // Truncate slice.
	}
	return data
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
	return RemoveElement(data, toDelete)
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
