package arrayutils

import (
	"fmt"
	"github.com/alessiosavi/GoGPUtils/datastructure/types"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
	"strings"
)

// Pad is delegated to pad/trim the given `data` to `n`, using the value `v`
func Pad[T any](data *[]T, n int, v T) []T {
	if len(*data) > n {
		return Trim(data, n)
	}
	res := make([]T, n-len(*data))
	for i := 0; i < cap(res); i++ {
		res[i] = v
	}
	*data = append(*data, res...)
	return *data
}

// Trim is delegated to trim the given `data` to `n`
func Trim[T any](data *[]T, n int) []T {
	if len(*data) > n {
		*data = (*data)[:n]
	}
	return *data
}

// RemoveElement is delegated to delete the element related to index i
func RemoveElement(s []string, i int) []string {
	if i < len(s) {
		s[len(s)-1], s[i] = s[i], s[len(s)-1]
		return s[:len(s)-1]
	}
	return s
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

// RemoveByIndex is delegated to remove the element of index s
func RemoveByIndex[T any](slice []T, s int) []T {
	if s < 0 || s >= len(slice) {
		return slice
	}
	slice = slices.Delete(slice, s, s+1)
	return slice
}

// RemoveByValue is delegated to remove the element that contains the given value
func RemoveByValue[T comparable](slice []T, v T) []T {
	for i := 0; i < len(slice); i++ {
		if slice[i] == v {
			slice = append(slice[:i], slice[i+1:]...)
			i--
		}
	}
	return slice
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

func Filter[T any](slice []T, f func(int, T) bool) []T {
	var n = make([]T, 0, len(slice))
	for i, e := range slice {
		if f(i, e) {
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

func Unique[T constraints.Ordered](slice []T) []T {
	var m = make(map[T]struct{})
	var i int
	for _, x := range slice {
		if _, ok := m[x]; !ok {
			m[x] = struct{}{}
			slice[i] = x
			i++
		}
	}

	return slices.Clip(slice[:i])
}

func Count[T comparable](data []T, target T) int {
	var c = 0
	for i := range data {
		if data[i] == target {
			c++
		}
	}
	return c
}

func Partition[T any](slice []T, f func(int, T) bool) ([]T, []T) {
	var ok []T
	var ko []T
	for i, e := range slice {
		if f(i, e) {
			ok = append(ok, e)
		} else {
			ko = append(ko, e)
		}
	}
	return ok, ko
}
func EachSlice(slice []int, size int) <-chan []int {
	ch := make(chan []int)
	go func() {
		defer close(ch)
		for i := 0; i < len(slice); i += size {
			end := i + size
			if end > len(slice) {
				end = len(slice)
			}
			ch <- slice[i:end]
		}
	}()
	return ch
}
