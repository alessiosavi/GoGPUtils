package arrayutils

import (
	"fmt"
	"github.com/alessiosavi/GoGPUtils/helper"
	"golang.org/x/exp/slices"
	"reflect"
	"strconv"
	"testing"
)

func TestRemoveElementsFromMatrixByIndex(t *testing.T) {

	var data = [][]string{
		{"0", "0", "0", "0", "0", "0", "0", "0", "0"},
		{"1", "1", "1", "1", "1", "1", "1", "1", "1"},
		{"2", "2", "2", "2", "2", "2", "2", "2", "2"},
		{"3", "3", "3", "3", "3", "3", "3", "3", "3"},
		{"4", "4", "4", "4", "4", "4", "4", "4", "4"},
		{"5", "5", "5", "5", "5", "5", "5", "5", "5"},
		{"6", "6", "6", "6", "6", "6", "6", "6", "6"},
	}
	var target = [][]string{
		{"0", "0", "0", "0", "0", "0", "0", "0", "0"},
		{"1", "1", "1", "1", "1", "1", "1", "1", "1"},
		{"4", "4", "4", "4", "4", "4", "4", "4", "4"},
		{"5", "5", "5", "5", "5", "5", "5", "5", "5"},
		{"6", "6", "6", "6", "6", "6", "6", "6", "6"},
	}

	data = RemoveElementsFromMatrixByIndex(data, []int{2, 3})

	if !reflect.DeepEqual(data, target) {
		t.Errorf("Got: %v\nExpected:%+v\n", data, target)
	}

}

func TestRemoveElementsFromString1(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	data = RemoveElementsFromStringByIndex(data, []int{0, 1, 2})
	for i := range data {
		if data[i] != strconv.Itoa(i+4) {
			t.Fail()
		}
	}
}

func TestRemoveElementsFromString2(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	data = RemoveElementsFromStringByIndex(data, []int{5, 4, 3})
	for i := range data {
		if data[i] != strconv.Itoa(i+1) {
			t.Fail()
		}
	}
}

func TestRemoveElementsFromString3(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	data = RemoveElementsFromStringByIndex(data, []int{1, 2, 3, 5, 4, 3, 0})
	for i := range data {
		if data[i] != strconv.Itoa(i+1) {
			t.Fail()
		}
	}
}

func TestRemoveElement(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	lenBefore := len(data)
	data = RemoveElement(data, 0)
	if lenBefore-1 != len(data) {
		t.Fail()
	}
}

func TestJoinStrings(t *testing.T) {
	data := []string{"a", "b", "c", "d", "e", "f", "g"}
	joined := JoinStrings(data, "")
	if joined != "abcdefg" {
		t.Fail()
	}
	joined = JoinStrings(data, " ")
	if joined != "a b c d e f g" {
		t.Fail()
	}
}
func TestJoinNumber(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	joined := JoinNumber[int](data, "")
	if joined != "0123456" {
		t.Fail()
	}
	joined = JoinNumber[int](data, " ")
	if joined != "0 1 2 3 4 5 6" {
		t.Fail()
	}
}
func TestReverseArrayInt(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	reversed := ReverseArray[int](data)
	for i, j := len(reversed)-1, 0; i > 0; i, j = i-1, j+1 {
		if reversed[i] != data[j] {
			t.Fail()
		}
	}
}
func TestReverseArrayString(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	reversed := ReverseArrayString(data)
	for i, j := len(reversed)-1, 0; i > 0; i, j = i-1, j+1 {
		if reversed[i] != data[j] {
			t.Fail()
		}
	}
}
func TestRemoveIntByIndex(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	test := []int{1, 2, 3, 4, 5, 6}
	deleted := RemoveByIndex[int](data, 0)
	if !reflect.DeepEqual(deleted, test) {
		t.Fail()
	}
}

func TestRemoveIntByIndex2(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	test := []int{0, 1, 2, 3, 4, 5}
	deleted := RemoveByIndex[int](data, 6)
	if !reflect.DeepEqual(deleted, test) {
		t.Fail()
	}
}

func TestRemoveIntByIndex3(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	deleted := RemoveByIndex[int](data, 7)
	if !reflect.DeepEqual(data, deleted) {
		t.Fail()
	}
}

func TestRemoveIntByIndex4(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6}
	deleted := RemoveByIndex[int](data, -1)
	if !reflect.DeepEqual(data, deleted) {
		t.Fail()
	}
}

func TestRemoveStringByIndex(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	test := []string{"2", "3", "4", "5", "6"}
	deleted := RemoveStringByIndex(data, 0)
	if !reflect.DeepEqual(deleted, test) {
		t.Fail()
	}
}

func TestRemoveStringByIndex2(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	test := []string{"1", "2", "3", "4", "5"}
	deleted := RemoveStringByIndex(data, 5)
	if !reflect.DeepEqual(deleted, test) {
		t.Log(deleted)
		t.Log(test)
	}
}

func TestRemoveStringByIndex3(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	deleted := RemoveStringByIndex(data, 7)
	if !reflect.DeepEqual(data, deleted) {
		t.Fail()
	}
}

func TestRemoveStringByIndex4(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	deleted := RemoveStringByIndex(data, -1)
	if !reflect.DeepEqual(data, deleted) {
		t.Fail()
	}
}
func TestRemoveStrings(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	test := []string{"1", "2", "3", "4", "5", "6"}
	deleted := RemoveStrings(data, test)
	if len(deleted) != 0 {
		t.Error(deleted)
	}
}

func TestRemoveStrings1(t *testing.T) {
	data := []string{"1", "2", "3", "4", "5", "6"}
	test := []string{"7", "8", "9", "10", "11", "12"}
	deleted := RemoveStrings(data, test)
	if !reflect.DeepEqual(data, deleted) {
		t.Log(data)
		t.Log(deleted)
	}
}

func TestSplitEqual(t *testing.T) {
	type args[T any] struct {
		data []T
		n    int
	}
	type testCase[T any] struct {
		name  string
		args  args[T]
		want  [][]T
		want1 []T
	}
	tests := []testCase[int]{
		{
			name: "OK",
			args: args[int]{
				data: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
				n:    5,
			},
			want:  [][]int{{0, 1, 2, 3, 4}, {5, 6, 7, 8, 9}},
			want1: []int{10, 11, 12, 13, 14},
		},

		{
			name: "OK",
			args: args[int]{
				data: []int{0, 1, 2, 3, 4},
				n:    2,
			},
			want:  [][]int{{0, 1}, {2, 3}},
			want1: []int{4},
		},
		{
			name: "ok",
			args: args[int]{
				data: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
				n:    6,
			},
			want:  [][]int{{0, 1, 2, 3, 4, 5}, {6, 7, 8, 9, 10, 11}},
			want1: []int{12, 13, 14},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SplitEqual(tt.args.data, tt.args.n)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitEqual() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SplitEqual() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func BenchmarkApplyInplace(b *testing.B) {
	data := helper.GenerateSequentialArray[byte](50000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Apply(&data, func(i int, v byte) byte {
			if v%2 == 0 {
				return v + 1
			} else {
				return v - 1
			}
		}, true)
	}
}

func BenchmarkApply(b *testing.B) {
	data := helper.GenerateSequentialArray[byte](50000)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Apply(&data, func(i int, v byte) byte {
			if v%2 == 0 {
				return v + 1
			} else {
				return v - 1
			}
		}, false)
	}
}

func TestApply(t *testing.T) {
	type args[T any] struct {
		v       *[]T
		fn      func(int, T) T
		inplace bool
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "1",
			args: args[int]{
				v: &[]int{1, 2, 3, 4, 5},
				fn: func(i int, v int) int {
					if i%2 == 0 {
						v++
					}
					return v
				},
				inplace: false,
			},
			want: []int{2, 2, 4, 4, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Apply(tt.args.v, tt.args.fn, tt.args.inplace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Apply() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPadSlice(t *testing.T) {
	type args[T any] struct {
		data   []T
		n      int
		expect []T
	}
	type testCase[T any] struct {
		name string
		args args[T]
	}
	tests := []testCase[string]{
		{
			name: "ok1",
			args: args[string]{
				data:   []string{"1", "2", "3"},
				n:      5,
				expect: []string{"1", "2", "3", "0", "0"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PadSlice(&tt.args.data, tt.args.n, "0")
		})
		if !slices.Equal(tt.args.data, tt.args.expect) {
			fmt.Println(tt.args)
		}

	}
}
