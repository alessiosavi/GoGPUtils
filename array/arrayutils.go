package arrayutils

import (
	"fmt"
	"github.com/alessiosavi/GoGPUtils/datastructure/types"
	"sort"
	"strings"
)

func PadSlice[T any](data *[]T, n int, v T) {
	if len(*data) >= n {
		*data = (*data)[:n]
		return
	}
	res := make([]T, n-len(*data))
	for i := 0; i < cap(res); i++ {
		res[i] = v
	}
	*data = append(*data, res...)
}
func RemoveElementsFromMatrixByIndex(data [][]string, j []int) [][]string {
	var (
		newArray [][]string
		toAdd    = true
	)

	if len(j) == 0 {
		return data
	}
	for i := 0; i < len(data); i++ {
		for _, k := range j {
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

// RemoveElementsFromStringByIndex delete the element of the indexes contained in j of the data in input
func RemoveElementsFromStringByIndex(data []string, j []int) []string {
	var (
		newArray []string
		toAdd    = true
	)

	if len(j) == 0 {
		return data
	}
	// sort.Ints(j)
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

// RemoveElement is delegated to delete the element related to index i
func RemoveElement(s []string, i int) []string {
	if i < len(s) {
		s[len(s)-1], s[i] = s[i], s[len(s)-1]
		return s[:len(s)-1]
	}
	return s
}

// JoinStrings use a strings.Builder for concatenate the input string array.
// It concatenates the strings among the delimiter in input
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

// JoinNumber use a strings.Builder for concatenate the input string array.
// It concatenates the strings among the delimiter in input
func JoinNumber[T types.Number](n []T, delimiter string) string {
	if len(n) == 0 {
		return ""
	}
	var sb strings.Builder

	for i := range n {
		sb.WriteString(fmt.Sprintf("%v", n[i]))
		sb.WriteString(delimiter)
	}
	return strings.TrimSuffix(sb.String(), delimiter)
}

// ReverseArray is delegated to return the inverse rappresentation of the array
// FIXME: Use the same array instead of allocate a new array
func ReverseArray[T types.Number](n1 []T) []T {
	var result = make([]T, len(n1))
	for i := len(n1) - 1; i >= 0; i-- {
		v := n1[i]
		j := len(n1) - 1 - i
		result[j] = v
	}
	return result
}

// ReverseArrayString is delegated to return the inverse rappresentation of the array
func ReverseArrayString(n1 []string) []string {
	var result = make([]string, len(n1))
	for i := len(n1) - 1; i >= 0; i-- {
		v := n1[i]
		j := len(n1) - 1 - i
		result[j] = v
	}
	return result
}

// RemoveByIndex is delegated to remove the element of index s
func RemoveByIndex[T types.Number](slice []T, s int) []T {
	if s < 0 || s >= len(slice) {
		return slice
	}
	return append(slice[:s], slice[s+1:]...)
}

// RemoveByValue is delegated to remove the element that contains the given value
func RemoveByValue[T types.Number](slice []T, value T) []T {
	for i := 0; i < len(slice); i++ {
		if slice[i] == value {
			slice = append(slice[:i], slice[i+1:]...)
			i--
		}
	}
	return slice
}

// InInt is delegated to verify if the given value is present in the target slice
func InInt(slice []int, target int) bool {
	for _, b := range slice {
		if b == target {
			return true
		}
	}
	return false
}

// InRune is delegated to verify if the given value is present in the target slice
func InRune(slice []rune, target rune) bool {
	for i := range slice {
		if slice[i] == target {
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
func RemoveStrings(slice, toRemove []string) []string {
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

// InStrings is delegated to verify if the given string arrays contains the target
func InStrings(slice []string, target string) bool {
	for _, element := range slice {
		if element == target {
			return true
		}
	}
	return false
}

// ContainStrings is delegated to verify if the given string arrays contains the target
func ContainStrings(slice []string, target string) bool {
	for _, element := range slice {
		if strings.Contains(target, element) {
			return true
		}
	}
	return false
}

func UniqueString(slice []string) []string {
	var m = make(map[string]struct{})
	for _, x := range slice {
		m[x] = struct{}{}
	}
	slice = []string{}
	for x := range m {
		slice = append(slice, x)
	}

	sort.Strings(slice)
	return slice
}

func ToByte(slice []string, separator string) []byte {
	var sb strings.Builder
	for i := range slice {
		sb.WriteString(slice[i])
		sb.WriteString(separator)
	}
	return []byte(strings.TrimSuffix(sb.String(), separator))
}

// SplitEqual is delegated to split the given data into slice of equal length
func SplitEqual[T any](data []T, n int) ([][]T, []T) {
	var ret = make([][]T, 0, len(data)/n+1)
	var i int
	for i = 0; i < len(data)-n; i += n {
		ret = append(ret, data[i:i+n])
	}
	return ret, data[i:]
}

func Filter[T any](slice []T, f func(T) bool) []T {
	var n = make([]T, 0, len(slice))
	for _, e := range slice {
		if f(e) {
			n = append(n, e)
		}
	}
	return n
}

func Apply[T any](v *[]T, fn func(int, T) T, inplace bool) []T {
	var res []T
	if inplace {
		res = *v
	} else {
		res = make([]T, len(*v))
		copy(res, *v)
	}
	for i := range res {
		res[i] = fn(i, res[i])
	}
	return res
}
