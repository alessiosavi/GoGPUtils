package sliceutil

import (
	"cmp"
	"slices"
	"testing"
)

// ============================================================================
// Filter Tests
// ============================================================================

func TestFilter(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		want      []int
	}{
		{
			name:      "filter evens",
			input:     []int{1, 2, 3, 4, 5, 6},
			predicate: func(n int) bool { return n%2 == 0 },
			want:      []int{2, 4, 6},
		},
		{
			name:      "filter none match",
			input:     []int{1, 3, 5},
			predicate: func(n int) bool { return n%2 == 0 },
			want:      []int{},
		},
		{
			name:      "filter all match",
			input:     []int{2, 4, 6},
			predicate: func(n int) bool { return n%2 == 0 },
			want:      []int{2, 4, 6},
		},
		{
			name:      "empty slice",
			input:     []int{},
			predicate: func(n int) bool { return true },
			want:      []int{},
		},
		{
			name:      "nil slice",
			input:     nil,
			predicate: func(n int) bool { return true },
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Filter(tt.input, tt.predicate)
			if !Equal(got, tt.want) {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterInPlace(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}
	got := FilterInPlace(input, func(n int) bool { return n%2 == 0 })
	want := []int{2, 4, 6}

	if !Equal(got, want) {
		t.Errorf("FilterInPlace() = %v, want %v", got, want)
	}

	// Verify it's the same underlying array
	if &got[0] != &input[0] {
		t.Error("FilterInPlace should use same underlying array")
	}
}

// ============================================================================
// Map Tests
// ============================================================================

func TestMap(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		transform func(int) int
		want      []int
	}{
		{
			name:      "double values",
			input:     []int{1, 2, 3},
			transform: func(n int) int { return n * 2 },
			want:      []int{2, 4, 6},
		},
		{
			name:      "empty slice",
			input:     []int{},
			transform: func(n int) int { return n },
			want:      []int{},
		},
		{
			name:      "nil slice",
			input:     nil,
			transform: func(n int) int { return n },
			want:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Map(tt.input, tt.transform)
			if !Equal(got, tt.want) {
				t.Errorf("Map() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMapTypeConversion(t *testing.T) {
	input := []int{1, 2, 3}
	got := Map(input, func(n int) string {
		return string(rune('a' + n - 1))
	})
	want := []string{"a", "b", "c"}

	if !slices.Equal(got, want) {
		t.Errorf("Map type conversion = %v, want %v", got, want)
	}
}

func TestMapWithIndex(t *testing.T) {
	input := []string{"a", "b", "c"}
	got := MapWithIndex(input, func(i int, s string) string {
		return s + string(rune('0'+i))
	})
	want := []string{"a0", "b1", "c2"}

	if !slices.Equal(got, want) {
		t.Errorf("MapWithIndex() = %v, want %v", got, want)
	}
}

// ============================================================================
// Reduce Tests
// ============================================================================

func TestReduce(t *testing.T) {
	t.Run("sum", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		got := Reduce(input, 0, func(acc, n int) int { return acc + n })
		want := 15
		if got != want {
			t.Errorf("Reduce sum = %v, want %v", got, want)
		}
	})

	t.Run("product", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		got := Reduce(input, 1, func(acc, n int) int { return acc * n })
		want := 24
		if got != want {
			t.Errorf("Reduce product = %v, want %v", got, want)
		}
	})

	t.Run("empty slice", func(t *testing.T) {
		input := []int{}
		got := Reduce(input, 42, func(acc, n int) int { return acc + n })
		want := 42
		if got != want {
			t.Errorf("Reduce empty = %v, want %v", got, want)
		}
	})
}

// ============================================================================
// Contains Tests
// ============================================================================

func TestContains(t *testing.T) {
	tests := []struct {
		name   string
		slice  []int
		target int
		want   bool
	}{
		{"found first", []int{1, 2, 3}, 1, true},
		{"found middle", []int{1, 2, 3}, 2, true},
		{"found last", []int{1, 2, 3}, 3, true},
		{"not found", []int{1, 2, 3}, 4, false},
		{"empty slice", []int{}, 1, false},
		{"nil slice", nil, 1, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Contains(tt.slice, tt.target); got != tt.want {
				t.Errorf("Contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContainsFunc(t *testing.T) {
	type person struct {
		name string
		age  int
	}
	people := []person{{"Alice", 30}, {"Bob", 25}, {"Charlie", 35}}

	if !ContainsFunc(people, func(p person) bool { return p.age > 30 }) {
		t.Error("ContainsFunc should find person over 30")
	}

	if ContainsFunc(people, func(p person) bool { return p.age > 40 }) {
		t.Error("ContainsFunc should not find person over 40")
	}
}

// ============================================================================
// IndexOf Tests
// ============================================================================

func TestIndexOf(t *testing.T) {
	tests := []struct {
		name   string
		slice  []string
		target string
		want   int
	}{
		{"found first", []string{"a", "b", "c"}, "a", 0},
		{"found middle", []string{"a", "b", "c"}, "b", 1},
		{"found last", []string{"a", "b", "c"}, "c", 2},
		{"not found", []string{"a", "b", "c"}, "d", -1},
		{"empty slice", []string{}, "a", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IndexOf(tt.slice, tt.target); got != tt.want {
				t.Errorf("IndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLastIndexOf(t *testing.T) {
	tests := []struct {
		name   string
		slice  []int
		target int
		want   int
	}{
		{"single occurrence", []int{1, 2, 3}, 2, 1},
		{"multiple occurrences", []int{1, 2, 3, 2, 4}, 2, 3},
		{"not found", []int{1, 2, 3}, 5, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LastIndexOf(tt.slice, tt.target); got != tt.want {
				t.Errorf("LastIndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ============================================================================
// Unique Tests
// ============================================================================

func TestUnique(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"with duplicates", []int{1, 2, 2, 3, 1, 4}, []int{1, 2, 3, 4}},
		{"no duplicates", []int{1, 2, 3}, []int{1, 2, 3}},
		{"all same", []int{1, 1, 1}, []int{1}},
		{"empty", []int{}, []int{}},
		{"nil", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Unique(tt.input)
			if !Equal(got, tt.want) {
				t.Errorf("Unique() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUniqueFunc(t *testing.T) {
	type item struct {
		id   int
		name string
	}
	input := []item{{1, "a"}, {2, "b"}, {1, "c"}, {3, "d"}}
	got := UniqueFunc(input, func(i item) int { return i.id })
	want := []item{{1, "a"}, {2, "b"}, {3, "d"}}

	if len(got) != len(want) {
		t.Errorf("UniqueFunc() length = %v, want %v", len(got), len(want))
	}
	for i := range got {
		if got[i].id != want[i].id {
			t.Errorf("UniqueFunc()[%d].id = %v, want %v", i, got[i].id, want[i].id)
		}
	}
}

// ============================================================================
// Chunk Tests
// ============================================================================

func TestChunk(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		size  int
		want  [][]int
	}{
		{
			name:  "even split",
			input: []int{1, 2, 3, 4, 5, 6},
			size:  2,
			want:  [][]int{{1, 2}, {3, 4}, {5, 6}},
		},
		{
			name:  "uneven split",
			input: []int{1, 2, 3, 4, 5},
			size:  2,
			want:  [][]int{{1, 2}, {3, 4}, {5}},
		},
		{
			name:  "size larger than slice",
			input: []int{1, 2},
			size:  5,
			want:  [][]int{{1, 2}},
		},
		{
			name:  "size of 1",
			input: []int{1, 2, 3},
			size:  1,
			want:  [][]int{{1}, {2}, {3}},
		},
		{
			name:  "empty slice",
			input: []int{},
			size:  2,
			want:  [][]int{},
		},
		{
			name:  "nil slice",
			input: nil,
			size:  2,
			want:  nil,
		},
		{
			name:  "zero size",
			input: []int{1, 2, 3},
			size:  0,
			want:  nil,
		},
		{
			name:  "negative size",
			input: []int{1, 2, 3},
			size:  -1,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Chunk(tt.input, tt.size)
			if len(got) != len(tt.want) {
				t.Errorf("Chunk() length = %v, want %v", len(got), len(tt.want))
				return
			}
			for i := range got {
				if !Equal(got[i], tt.want[i]) {
					t.Errorf("Chunk()[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

// ============================================================================
// Flatten Tests
// ============================================================================

func TestFlatten(t *testing.T) {
	tests := []struct {
		name  string
		input [][]int
		want  []int
	}{
		{
			name:  "normal",
			input: [][]int{{1, 2}, {3}, {4, 5, 6}},
			want:  []int{1, 2, 3, 4, 5, 6},
		},
		{
			name:  "with empty inner",
			input: [][]int{{1, 2}, {}, {3}},
			want:  []int{1, 2, 3},
		},
		{
			name:  "all empty",
			input: [][]int{{}, {}, {}},
			want:  []int{},
		},
		{
			name:  "nil",
			input: nil,
			want:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Flatten(tt.input)
			if !Equal(got, tt.want) {
				t.Errorf("Flatten() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ============================================================================
// Reverse Tests
// ============================================================================

func TestReverse(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"normal", []int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
		{"single", []int{1}, []int{1}},
		{"two", []int{1, 2}, []int{2, 1}},
		{"empty", []int{}, []int{}},
		{"nil", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Reverse(tt.input)
			if !Equal(got, tt.want) {
				t.Errorf("Reverse() = %v, want %v", got, tt.want)
			}
			// Verify original unchanged
			if tt.input != nil && len(tt.input) > 0 {
				original := []int{1, 2, 3, 4}[:len(tt.input)]
				for i, v := range tt.input {
					if v != original[i] {
						t.Error("Reverse modified original slice")
						break
					}
				}
			}
		})
	}
}

func TestReverseInPlace(t *testing.T) {
	input := []int{1, 2, 3, 4}
	ReverseInPlace(input)
	want := []int{4, 3, 2, 1}

	if !Equal(input, want) {
		t.Errorf("ReverseInPlace() = %v, want %v", input, want)
	}
}

// ============================================================================
// Set Operations Tests
// ============================================================================

func TestIntersect(t *testing.T) {
	tests := []struct {
		name string
		a, b []int
		want []int
	}{
		{"normal", []int{1, 2, 3, 4}, []int{2, 4, 6}, []int{2, 4}},
		{"no overlap", []int{1, 2, 3}, []int{4, 5, 6}, []int{}},
		{"same", []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}},
		{"a nil", nil, []int{1, 2}, nil},
		{"b nil", []int{1, 2}, nil, nil},
		{"with duplicates", []int{1, 2, 2, 3}, []int{2, 2, 4}, []int{2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Intersect(tt.a, tt.b)
			if !Equal(got, tt.want) {
				t.Errorf("Intersect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDifference(t *testing.T) {
	tests := []struct {
		name string
		a, b []int
		want []int
	}{
		{"normal", []int{1, 2, 3, 4}, []int{2, 4}, []int{1, 3}},
		{"no overlap", []int{1, 2, 3}, []int{4, 5}, []int{1, 2, 3}},
		{"all overlap", []int{1, 2}, []int{1, 2, 3}, []int{}},
		{"a nil", nil, []int{1}, nil},
		{"b nil", []int{1, 2}, nil, []int{1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Difference(tt.a, tt.b)
			if !Equal(got, tt.want) {
				t.Errorf("Difference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUnion(t *testing.T) {
	tests := []struct {
		name string
		a, b []int
		want []int
	}{
		{"normal", []int{1, 2}, []int{2, 3}, []int{1, 2, 3}},
		{"no overlap", []int{1, 2}, []int{3, 4}, []int{1, 2, 3, 4}},
		{"same", []int{1, 2}, []int{1, 2}, []int{1, 2}},
		{"both nil", nil, nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Union(tt.a, tt.b)
			if !Equal(got, tt.want) {
				t.Errorf("Union() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ============================================================================
// GroupBy Tests
// ============================================================================

func TestGroupBy(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}
	got := GroupBy(input, func(n int) string {
		if n%2 == 0 {
			return "even"
		}
		return "odd"
	})

	if !Equal(got["even"], []int{2, 4, 6}) {
		t.Errorf("GroupBy even = %v, want [2, 4, 6]", got["even"])
	}
	if !Equal(got["odd"], []int{1, 3, 5}) {
		t.Errorf("GroupBy odd = %v, want [1, 3, 5]", got["odd"])
	}

	// Test nil
	if GroupBy[int, string](nil, func(n int) string { return "" }) != nil {
		t.Error("GroupBy(nil) should return nil")
	}
}

// ============================================================================
// Partition Tests
// ============================================================================

func TestPartition(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	evens, odds := Partition(input, func(n int) bool { return n%2 == 0 })

	if !Equal(evens, []int{2, 4}) {
		t.Errorf("Partition evens = %v, want [2, 4]", evens)
	}
	if !Equal(odds, []int{1, 3, 5}) {
		t.Errorf("Partition odds = %v, want [1, 3, 5]", odds)
	}

	// Test nil
	nilA, nilB := Partition[int](nil, func(n int) bool { return true })
	if nilA != nil || nilB != nil {
		t.Error("Partition(nil) should return nil, nil")
	}
}

// ============================================================================
// Take/Drop Tests
// ============================================================================

func TestTake(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		n     int
		want  []int
	}{
		{"normal", []int{1, 2, 3, 4, 5}, 3, []int{1, 2, 3}},
		{"more than length", []int{1, 2}, 5, []int{1, 2}},
		{"zero", []int{1, 2, 3}, 0, nil},
		{"negative", []int{1, 2, 3}, -1, nil},
		{"nil slice", nil, 2, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Take(tt.input, tt.n)
			if !Equal(got, tt.want) {
				t.Errorf("Take() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTakeLast(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	got := TakeLast(input, 3)
	want := []int{3, 4, 5}

	if !Equal(got, want) {
		t.Errorf("TakeLast() = %v, want %v", got, want)
	}
}

func TestDrop(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		n     int
		want  []int
	}{
		{"normal", []int{1, 2, 3, 4, 5}, 2, []int{3, 4, 5}},
		{"drop all", []int{1, 2, 3}, 5, []int{}},
		{"drop none", []int{1, 2, 3}, 0, []int{1, 2, 3}},
		{"negative", []int{1, 2, 3}, -1, []int{1, 2, 3}},
		{"nil", nil, 2, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Drop(tt.input, tt.n)
			if !Equal(got, tt.want) {
				t.Errorf("Drop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTakeWhile(t *testing.T) {
	input := []int{1, 2, 3, 4, 1, 2}
	got := TakeWhile(input, func(n int) bool { return n < 4 })
	want := []int{1, 2, 3}

	if !Equal(got, want) {
		t.Errorf("TakeWhile() = %v, want %v", got, want)
	}
}

func TestDropWhile(t *testing.T) {
	input := []int{1, 2, 3, 4, 1, 2}
	got := DropWhile(input, func(n int) bool { return n < 4 })
	want := []int{4, 1, 2}

	if !Equal(got, want) {
		t.Errorf("DropWhile() = %v, want %v", got, want)
	}
}

// ============================================================================
// All/Any/None Tests
// ============================================================================

func TestAll(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		want      bool
	}{
		{"all match", []int{2, 4, 6}, func(n int) bool { return n%2 == 0 }, true},
		{"some match", []int{2, 3, 6}, func(n int) bool { return n%2 == 0 }, false},
		{"none match", []int{1, 3, 5}, func(n int) bool { return n%2 == 0 }, false},
		{"empty", []int{}, func(n int) bool { return false }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := All(tt.input, tt.predicate); got != tt.want {
				t.Errorf("All() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAny(t *testing.T) {
	tests := []struct {
		name      string
		input     []int
		predicate func(int) bool
		want      bool
	}{
		{"all match", []int{2, 4, 6}, func(n int) bool { return n%2 == 0 }, true},
		{"some match", []int{1, 2, 3}, func(n int) bool { return n%2 == 0 }, true},
		{"none match", []int{1, 3, 5}, func(n int) bool { return n%2 == 0 }, false},
		{"empty", []int{}, func(n int) bool { return true }, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Any(tt.input, tt.predicate); got != tt.want {
				t.Errorf("Any() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNone(t *testing.T) {
	if !None([]int{1, 3, 5}, func(n int) bool { return n%2 == 0 }) {
		t.Error("None should return true when no elements match")
	}
	if None([]int{1, 2, 3}, func(n int) bool { return n%2 == 0 }) {
		t.Error("None should return false when some elements match")
	}
}

func TestCount(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}
	got := Count(input, func(n int) bool { return n%2 == 0 })
	want := 3

	if got != want {
		t.Errorf("Count() = %v, want %v", got, want)
	}
}

// ============================================================================
// Find Tests
// ============================================================================

func TestFind(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}

	// Find existing
	val, ok := Find(input, func(n int) bool { return n > 3 })
	if !ok || val != 4 {
		t.Errorf("Find() = %v, %v; want 4, true", val, ok)
	}

	// Not found
	val, ok = Find(input, func(n int) bool { return n > 10 })
	if ok {
		t.Errorf("Find() should return false for no match")
	}
}

func TestFindLast(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}

	val, ok := FindLast(input, func(n int) bool { return n%2 == 0 })
	if !ok || val != 4 {
		t.Errorf("FindLast() = %v, %v; want 4, true", val, ok)
	}
}

// ============================================================================
// Min/Max Tests
// ============================================================================

func TestMin(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  int
		ok    bool
	}{
		{"normal", []int{3, 1, 4, 1, 5}, 1, true},
		{"single", []int{42}, 42, true},
		{"empty", []int{}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Min(tt.input)
			if ok != tt.ok {
				t.Errorf("Min() ok = %v, want %v", ok, tt.ok)
			}
			if ok && got != tt.want {
				t.Errorf("Min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  int
		ok    bool
	}{
		{"normal", []int{3, 1, 4, 1, 5}, 5, true},
		{"single", []int{42}, 42, true},
		{"empty", []int{}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Max(tt.input)
			if ok != tt.ok {
				t.Errorf("Max() ok = %v, want %v", ok, tt.ok)
			}
			if ok && got != tt.want {
				t.Errorf("Max() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMinMaxFunc(t *testing.T) {
	type item struct {
		val int
	}
	items := []item{{3}, {1}, {4}, {1}, {5}}
	cmpFunc := func(a, b item) int { return cmp.Compare(a.val, b.val) }

	min, _ := MinFunc(items, cmpFunc)
	if min.val != 1 {
		t.Errorf("MinFunc() = %v, want 1", min.val)
	}

	max, _ := MaxFunc(items, cmpFunc)
	if max.val != 5 {
		t.Errorf("MaxFunc() = %v, want 5", max.val)
	}
}

// ============================================================================
// Equal Tests
// ============================================================================

func TestEqual(t *testing.T) {
	tests := []struct {
		name string
		a, b []int
		want bool
	}{
		{"equal", []int{1, 2, 3}, []int{1, 2, 3}, true},
		{"different length", []int{1, 2}, []int{1, 2, 3}, false},
		{"different values", []int{1, 2, 3}, []int{1, 2, 4}, false},
		{"both empty", []int{}, []int{}, true},
		{"both nil", nil, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Equal(tt.a, tt.b); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ============================================================================
// Pad Tests
// ============================================================================

func TestPad(t *testing.T) {
	tests := []struct {
		name   string
		input  []int
		length int
		fill   int
		want   []int
	}{
		{"pad needed", []int{1, 2}, 5, 0, []int{1, 2, 0, 0, 0}},
		{"no pad needed", []int{1, 2, 3}, 2, 0, []int{1, 2, 3}},
		{"exact length", []int{1, 2, 3}, 3, 0, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Pad(tt.input, tt.length, tt.fill)
			if !Equal(got, tt.want) {
				t.Errorf("Pad() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPadLeft(t *testing.T) {
	got := PadLeft([]int{1, 2}, 5, 0)
	want := []int{0, 0, 0, 1, 2}

	if !Equal(got, want) {
		t.Errorf("PadLeft() = %v, want %v", got, want)
	}
}

// ============================================================================
// Remove Tests
// ============================================================================

func TestRemoveAt(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		index int
		want  []int
	}{
		{"remove first", []int{1, 2, 3}, 0, []int{2, 3}},
		{"remove middle", []int{1, 2, 3}, 1, []int{1, 3}},
		{"remove last", []int{1, 2, 3}, 2, []int{1, 2}},
		{"out of bounds negative", []int{1, 2, 3}, -1, nil},
		{"out of bounds positive", []int{1, 2, 3}, 5, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveAt(tt.input, tt.index)
			if !Equal(got, tt.want) {
				t.Errorf("RemoveAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRemoveValue(t *testing.T) {
	got := RemoveValue([]int{1, 2, 3, 2, 4}, 2)
	want := []int{1, 3, 4}

	if !Equal(got, want) {
		t.Errorf("RemoveValue() = %v, want %v", got, want)
	}
}

func TestRemoveFirst(t *testing.T) {
	got := RemoveFirst([]int{1, 2, 3, 2, 4}, 2)
	want := []int{1, 3, 2, 4}

	if !Equal(got, want) {
		t.Errorf("RemoveFirst() = %v, want %v", got, want)
	}
}

// ============================================================================
// Insert Tests
// ============================================================================

func TestInsert(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		index int
		value int
		want  []int
	}{
		{"insert at start", []int{2, 3}, 0, 1, []int{1, 2, 3}},
		{"insert in middle", []int{1, 3}, 1, 2, []int{1, 2, 3}},
		{"insert at end", []int{1, 2}, 2, 3, []int{1, 2, 3}},
		{"insert beyond end", []int{1, 2}, 10, 3, []int{1, 2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Insert(tt.input, tt.index, tt.value)
			if !Equal(got, tt.want) {
				t.Errorf("Insert() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ============================================================================
// Compact Tests
// ============================================================================

func TestCompact(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"with consecutive dups", []int{1, 1, 2, 2, 2, 3, 1}, []int{1, 2, 3, 1}},
		{"no consecutive dups", []int{1, 2, 3}, []int{1, 2, 3}},
		{"all same", []int{1, 1, 1, 1}, []int{1}},
		{"single", []int{1}, []int{1}},
		{"empty", []int{}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compact(tt.input)
			if !Equal(got, tt.want) {
				t.Errorf("Compact() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ============================================================================
// FlatMap Tests
// ============================================================================

func TestFlatMap(t *testing.T) {
	input := []string{"hello world", "foo bar"}
	got := FlatMap(input, func(s string) []byte { return []byte(s) })

	// Just verify length - content would be ASCII values
	if len(got) != 18 { // "hello world" (11) + "foo bar" (7) = 18 bytes
		t.Errorf("FlatMap() length = %v, want 18", len(got))
	}
}

// ============================================================================
// Associate Tests
// ============================================================================

func TestAssociate(t *testing.T) {
	type user struct {
		id   int
		name string
	}
	users := []user{{1, "Alice"}, {2, "Bob"}}
	got := Associate(users, func(u user) int { return u.id })

	if got[1].name != "Alice" || got[2].name != "Bob" {
		t.Errorf("Associate() = %v", got)
	}
}

func TestAssociateWith(t *testing.T) {
	keys := []string{"a", "bb", "ccc"}
	got := AssociateWith(keys, func(s string) int { return len(s) })

	if got["a"] != 1 || got["bb"] != 2 || got["ccc"] != 3 {
		t.Errorf("AssociateWith() = %v", got)
	}
}

// ============================================================================
// ZipWith Tests
// ============================================================================

func TestZipWith(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5, 6}
	got := ZipWith(a, b, func(x, y int) int { return x + y })
	want := []int{5, 7, 9}

	if !Equal(got, want) {
		t.Errorf("ZipWith() = %v, want %v", got, want)
	}
}

func TestZipWithDifferentLengths(t *testing.T) {
	a := []int{1, 2, 3, 4}
	b := []int{10, 20}
	got := ZipWith(a, b, func(x, y int) int { return x + y })
	want := []int{11, 22}

	if !Equal(got, want) {
		t.Errorf("ZipWith() = %v, want %v", got, want)
	}
}

// ============================================================================
// Shuffle Tests
// ============================================================================

func TestShuffle(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	original := slices.Clone(input)
	shuffled := Shuffle(input)

	// Verify original unchanged
	if !Equal(input, original) {
		t.Error("Shuffle modified original slice")
	}

	// Verify same elements
	slices.Sort(shuffled)
	if !Equal(shuffled, original) {
		t.Error("Shuffle changed elements")
	}
}

func TestShuffleDeterministic(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}

	SeedShuffle(42)
	first := Shuffle(input)

	SeedShuffle(42)
	second := Shuffle(input)

	if !Equal(first, second) {
		t.Error("Shuffle with same seed should produce same result")
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkFilter(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Filter(data, func(n int) bool { return n%2 == 0 })
	}
}

func BenchmarkFilterInPlace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		data := make([]int, 10000)
		for j := range data {
			data[j] = j
		}
		b.StartTimer()
		FilterInPlace(data, func(n int) bool { return n%2 == 0 })
	}
}

func BenchmarkMap(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Map(data, func(n int) int { return n * 2 })
	}
}

func BenchmarkContains(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Contains(data, 9999)
	}
}

func BenchmarkUnique(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i % 100 // Lots of duplicates
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unique(data)
	}
}

func BenchmarkChunk(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Chunk(data, 100)
	}
}

func BenchmarkGroupBy(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GroupBy(data, func(n int) int { return n % 10 })
	}
}

func BenchmarkIntersect(b *testing.B) {
	a := make([]int, 5000)
	bSlice := make([]int, 5000)
	for i := range a {
		a[i] = i
		bSlice[i] = i + 2500
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Intersect(a, bSlice)
	}
}
