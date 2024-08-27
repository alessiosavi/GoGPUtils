package arrayutils

import (
	"fmt"
	"github.com/alessiosavi/GoGPUtils/datastructure/types"
	"github.com/alessiosavi/GoGPUtils/helper"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
	"reflect"
	"testing"
)

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
			name: "ok_1",
			args: args[string]{
				data:   []string{"1", "2", "3"},
				n:      5,
				expect: []string{"1", "2", "3", "0", "0"},
			},
		},
		{
			name: "ok_2",
			args: args[string]{
				data:   []string{"1", "2", "3"},
				n:      3,
				expect: []string{"1", "2", "3"},
			},
		},
		{
			name: "ok_3",
			args: args[string]{
				data:   []string{"1", "2", "3"},
				n:      2,
				expect: []string{"1", "2"},
			},
		},
		{
			name: "ok_4",
			args: args[string]{
				data:   []string{"1", "2", "3"},
				n:      0,
				expect: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Pad(&tt.args.data, tt.args.n, "0")
		})
		if !slices.Equal(tt.args.data, tt.args.expect) {
			fmt.Println(tt.args)
		}

	}
}

func TestTrimSlice(t *testing.T) {
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
			name: "ok_1",
			args: args[string]{
				data:   []string{"1", "2", "3", "0", "0"},
				n:      3,
				expect: []string{"1", "2", "3"},
			},
		},
		{
			name: "ok_2",
			args: args[string]{
				data:   []string{"1", "2", "3"},
				n:      3,
				expect: []string{"1", "2", "3"},
			},
		},
		{
			name: "ok_3",
			args: args[string]{
				data:   []string{"1", "2", "3"},
				n:      0,
				expect: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Trim(&tt.args.data, tt.args.n)
		})
		if !slices.Equal(tt.args.data, tt.args.expect) {
			fmt.Println(tt.args)
		}

	}
}

func TestJoinNumber1(t *testing.T) {
	type args[T types.Number] struct {
		n         []T
		delimiter string
	}
	type testCase[T types.Number] struct {
		name string
		args args[T]
		want string
	}
	testFloat := []testCase[float32]{
		{
			name: "ok_1",
			args: args[float32]{
				n:         []float32{0, 1, 2, 3, 4, 5, 6},
				delimiter: "",
			},
			want: "0123456",
		},
		{
			name: "ok_2",
			args: args[float32]{
				n:         []float32{0.0, 1.1, 2.2, 3.3, 4.4, 5.5, 6.6},
				delimiter: " ",
			},
			want: "0 1.1 2.2 3.3 4.4 5.5 6.6",
		},
	}
	testInt := []testCase[int]{
		{
			name: "ok_1",
			args: args[int]{
				n:         []int{0, 1, 2, 3, 4, 5, 6},
				delimiter: "",
			},
			want: "0123456",
		},
		{
			name: "ok_1",
			args: args[int]{
				n:         []int{0, 1, 2, 3, 4, 5, 6},
				delimiter: " ",
			},
			want: "0 1 2 3 4 5 6",
		},
	}
	for _, tt := range testFloat {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinNumber(tt.args.n, tt.args.delimiter); got != tt.want {
				t.Errorf("JoinNumber() = %v, want %v", got, tt.want)
			}
		})
	}

	for _, tt := range testInt {
		t.Run(tt.name, func(t *testing.T) {
			if got := JoinNumber(tt.args.n, tt.args.delimiter); got != tt.want {
				t.Errorf("JoinNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveByIndex(t *testing.T) {
	type args[T any] struct {
		slice []T
		s     int
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want []T
	}
	testsInt := []testCase[int]{
		{
			name: "INT_1",
			args: args[int]{
				slice: []int{0, 1, 2, 3, 4, 5, 6},
				s:     0,
			},
			want: []int{1, 2, 3, 4, 5, 6},
		},

		{
			name: "INT_2",
			args: args[int]{
				slice: []int{0, 1, 2, 3, 4, 5, 6},
				s:     6,
			},
			want: []int{0, 1, 2, 3, 4, 5},
		},
		{
			name: "INT_3",
			args: args[int]{
				slice: []int{0, 1, 2, 3, 4, 5, 6},
				s:     7,
			},
			want: []int{0, 1, 2, 3, 4, 5, 6},
		},
		{
			name: "INT_4",
			args: args[int]{
				slice: []int{0, 1, 2, 3, 4, 5, 6},
				s:     -1,
			},
			want: []int{0, 1, 2, 3, 4, 5, 6},
		},
	}
	testsString := []testCase[string]{
		{
			name: "STRING_1",
			args: args[string]{
				slice: []string{"1", "2", "3", "4", "5", "6"},
				s:     0,
			},
			want: []string{"2", "3", "4", "5", "6"},
		},

		{
			name: "STRING_2",
			args: args[string]{
				slice: []string{"1", "2", "3", "4", "5", "6"},
				s:     5,
			},
			want: []string{"1", "2", "3", "4", "5"},
		},
		{
			name: "STRING_3",
			args: args[string]{
				slice: []string{"1", "2", "3", "4", "5", "6"},
				s:     7,
			},
			want: []string{"1", "2", "3", "4", "5", "6"},
		},
		{
			name: "STRING_4",
			args: args[string]{
				slice: []string{"1", "2", "3", "4", "5", "6"},
				s:     -1,
			},
			want: []string{"1", "2", "3", "4", "5", "6"},
		},
	}
	for _, tt := range testsInt {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveByIndex(tt.args.slice, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveByIndex() = %v, want %v", got, tt.want)
			}
		})
	}

	for _, tt := range testsString {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveByIndex(tt.args.slice, tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveByIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveByValue(t *testing.T) {
	type args[T types.Number] struct {
		slice []T
		v     T
	}
	type testCase[T types.Number] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "ok_1",
			args: args[int]{
				slice: []int{6, 5, 4, 3, 2, 1, 0},
				v:     0,
			},
			want: []int{6, 5, 4, 3, 2, 1},
		},

		{
			name: "ok_2",
			args: args[int]{
				slice: []int{6, 5, 3, 3, 2, 3, 0},
				v:     3,
			},
			want: []int{6, 5, 2, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveByValue(tt.args.slice, tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveByValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	type args[T constraints.Ordered] struct {
		slice []T
	}
	type testCase[T constraints.Ordered] struct {
		name string
		args args[T]
		want []T
	}
	tests := []testCase[int]{
		{
			name: "1",
			args: args[int]{
				slice: []int{0, 1, 2, 3, 4, 5, 5, 6, 7},
			},
			want: []int{0, 1, 2, 3, 4, 5, 6, 7},
		},
		{
			name: "2",
			args: args[int]{
				slice: []int{0, 0, 0, 0, 0, 0},
			},
			want: []int{0},
		},
		{
			name: "3",
			args: args[int]{
				slice: []int{5, 4, 3, 2, 1, 1},
			},
			want: []int{5, 4, 3, 2, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Unique(tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_eachSlice(t *testing.T) {
	type args struct {
		slice []int
		size  int
	}

	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{
			name: "ok1",
			args: args{
				slice: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				size:  3,
			},
			want: [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
		},
		{
			name: "ok2",
			args: args{
				slice: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				size:  3,
			},
			want: [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10}},
		},
		{
			name: "ok3",
			args: args{
				slice: []int{1, 2},
				size:  3,
			},
			want: [][]int{{1, 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := EachSlice(tt.args.slice, tt.args.size)
			var got [][]int
			for v := range res {
				got = append(got, v)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("eachSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPartition(t *testing.T) {
	type args[T any] struct {
		arr       []T
		condition func(int, T) bool
	}
	type testCase[T any] struct {
		name               string
		args               args[T]
		wantSatisfies      []T
		wantDoesNotSatisfy []T
	}
	tests := []testCase[int]{
		{
			name: "even_odd",
			args: args[int]{
				arr: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
				condition: func(i, v int) bool {
					return v%2 == 0
				},
			},
			wantSatisfies:      Filter([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, func(i, v int) bool { return v%2 == 0 }),
			wantDoesNotSatisfy: Filter([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, func(i, v int) bool { return v%2 != 0 }),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSatisfies, gotDoesNotSatisfy := Partition(tt.args.arr, tt.args.condition)
			if !reflect.DeepEqual(gotSatisfies, tt.wantSatisfies) {
				t.Errorf("Partition() gotSatisfies = %v, want %v", gotSatisfies, tt.wantSatisfies)
			}
			if !reflect.DeepEqual(gotDoesNotSatisfy, tt.wantDoesNotSatisfy) {
				t.Errorf("Partition() gotDoesNotSatisfy = %v, want %v", gotDoesNotSatisfy, tt.wantDoesNotSatisfy)
			}
		})
	}
}
