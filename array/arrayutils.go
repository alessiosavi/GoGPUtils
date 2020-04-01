package arrayutils

import (
	"strconv"
	"strings"
)

//RemoveElementsFromStringByIndex delete the element of the indexes contained in j of the data in input
func RemoveElementsFromStringByIndex(data []string, j []int) []string {
	var (
		newArray []string
		toAdd    bool = true
	)

	if len(j) == 0 {
		return data
	}
	//sort.Ints(j)
	for i := 0; i < len(data); i++ {
		for _, k := range j {
			// if k < i || k > i {
			// 	break
			// } else
			if i == k {
				toAdd = false
				break
			}
		}

		if toAdd {
			newArray = append(newArray, data[i])
		}

		toAdd = true
	}

	return newArray
}

//RemoveElement is delegated to delete the element related to index i
func RemoveElement(s []string, i int) []string {
	if i < len(s) {
		s[len(s)-1], s[i] = s[i], s[len(s)-1]
		return s[:len(s)-1]
	}
	return s
}

// JoinStrings use a strings.Builder for concatenate the input string array.
// It concatenate the strings among the delimiter in input
func JoinStrings(strs []string, delimiter string) string {
	if len(strs) == 0 {
		return ""
	}

	var sb strings.Builder

	for i := range strs {
		sb.WriteString(strs[i])
		sb.WriteString(delimiter)
	}

	return strings.TrimSuffix(sb.String(), delimiter)
}

// JoinInts use a strings.Builder for concatenate the input string array.
// It concatenate the strings among the delimiter in input
func JoinInts(ints []int, delimiter string) string {
	if len(ints) == 0 {
		return ""
	}
	var sb strings.Builder

	for i := range ints {
		sb.WriteString(strconv.Itoa(ints[i]))
		sb.WriteString(delimiter)
	}
	return strings.TrimSuffix(sb.String(), delimiter)
}

// ReverseArrayInt is delegated to return the inverse rappresentation of the array
func ReverseArrayInt(n1 []int) []int {
	var result []int = make([]int, len(n1))
	for i := len(n1) - 1; i >= 0; i-- {
		v := n1[i]
		j := len(n1) - 1 - i
		result[j] = v
	}
	return result
}

// ReverseArrayString is delegated to return the inverse rappresentation of the array
func ReverseArrayString(n1 []string) []string {
	var result []string = make([]string, len(n1))
	for i := len(n1) - 1; i >= 0; i-- {
		v := n1[i]
		j := len(n1) - 1 - i
		result[j] = v
	}
	return result
}

// RemoveIntByIndex is delegated to remove the element of index s
func RemoveIntByIndex(slice []int, s int) []int {
	if s < 0 || s >= len(slice) {
		return slice
	}
	return append(slice[:s], slice[s+1:]...)
}

// RemoveIntByValue is delegated to remove the element that contains the given value
func RemoveIntByValue(slice []int, value int) []int {
	for i := 0; i < len(slice); i++ {
		if slice[i] == value {
			slice = append(slice[:i], slice[i+1:]...)
			i--
		}
	}
	return slice
}

// In is delegated to verify if the given value is present in the target slice
func InInt(slice []int, target int) bool {
	for _, b := range slice {
		if b == target {
			return true
		}
	}
	return false
}

// In is delegated to verify if the given value is present in the target slice
func InRune(slice []rune, target rune) bool {
	for _, b := range slice {
		if b == target {
			return true
		}
	}
	return false
}

// RemoveStringByIndex the item in position s from the input array
func RemoveStringByIndex(slice []string, s int) []string {
	if s < 0 || s >= len(slice) {
		return slice
	}
	return append(slice[:s], slice[s+1:]...)
}

// RemoveStrings is delegated to remove the input 'toRemove' value from the given slice
func RemoveStrings(slice []string, toRemove []string) []string {
	for i := 0; i < len(slice); i++ {
		for j := 0; j < len(toRemove); j++ {
			if slice[i] == toRemove[j] {
				slice = RemoveStringByIndex(slice, i)
				// reset the index
				i--
				break
			}
		}
	}
	return slice
}
