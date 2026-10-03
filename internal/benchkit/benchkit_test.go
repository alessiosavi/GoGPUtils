package benchkit

import (
	"os"
	"os/exec"
	"reflect"
	"slices"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestGeneratorsAreDeterministic(t *testing.T) {
	gens := map[string]func() any{
		"Ints":       func() any { return Ints(100) },
		"SortedInts": func() any { return SortedInts(100) },
		"Floats":     func() any { return Floats(100) },
		"Strings":    func() any { return Strings(10, 7, "abcé") },
		"Words":      func() any { return Words(50) },
		"Text":       func() any { return Text(500, Mixed) },
		"Bytes":      func() any { return Bytes(33) },
		"Matrix":     func() any { return Matrix(3, 4) },
	}
	for name, gen := range gens {
		if a, b := gen(), gen(); !reflect.DeepEqual(a, b) {
			t.Errorf("%s returned different results for the same arguments", name)
		}
	}
	first := Ints(64)
	_ = Floats(64)
	_ = Text(64, ASCII)
	if !slices.Equal(first, Ints(64)) {
		t.Error("Ints depends on earlier generator calls")
	}
}

func TestShapes(t *testing.T) {
	for _, n := range []int{0, 1, 17, 1000} {
		ints := Ints(n)
		if len(ints) != n {
			t.Fatalf("len(Ints(%d)) = %d", n, len(ints))
		}
		for _, v := range ints {
			if v < 0 {
				t.Fatalf("Ints(%d) contains negative %d", n, v)
			}
		}
		if s := SortedInts(n); len(s) != n || !slices.IsSorted(s) {
			t.Fatalf("SortedInts(%d) = len %d, sorted %v", n, len(s), slices.IsSorted(s))
		}
		floats := Floats(n)
		if len(floats) != n {
			t.Fatalf("len(Floats(%d)) = %d", n, len(floats))
		}
		for _, f := range floats {
			if f < 0 || f >= 1 {
				t.Fatalf("Floats(%d) contains %v outside [0, 1)", n, f)
			}
		}
		if b := Bytes(n); len(b) != n || cap(b) != n {
			t.Fatalf("Bytes(%d) has len %d cap %d", n, len(b), cap(b))
		}
		if w := Words(n); len(w) != n {
			t.Fatalf("len(Words(%d)) = %d", n, len(w))
		}
	}
	for _, w := range Words(100) {
		if w == "" || !isASCII(w) {
			t.Fatalf("Words produced %q", w)
		}
	}
	strs := Strings(20, 9, "xyé")
	if len(strs) != 20 {
		t.Fatalf("len(Strings) = %d", len(strs))
	}
	for _, s := range strs {
		if utf8.RuneCountInString(s) != 9 {
			t.Fatalf("Strings produced %q with %d runes, want 9", s, utf8.RuneCountInString(s))
		}
		for _, r := range s {
			if !strings.ContainsRune("xyé", r) {
				t.Fatalf("Strings produced rune %q outside the alphabet", r)
			}
		}
	}
	m := Matrix(3, 5)
	if len(m) != 3 {
		t.Fatalf("Matrix rows = %d", len(m))
	}
	for _, row := range m {
		if len(row) != 5 {
			t.Fatalf("Matrix cols = %d", len(row))
		}
	}
}

func TestStringsPanicsOnEmptyAlphabet(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Strings(1, 1, \"\") did not panic")
		}
	}()
	Strings(1, 1, "")
}

func TestText(t *testing.T) {
	for _, kind := range []TextKind{ASCII, Unicode, Mixed} {
		for _, n := range []int{0, 1, 7, 64, 4096} {
			s := Text(n, kind)
			if len(s) != n {
				t.Fatalf("len(Text(%d, %d)) = %d", n, kind, len(s))
			}
			if !utf8.ValidString(s) {
				t.Fatalf("Text(%d, %d) is not valid UTF-8", n, kind)
			}
		}
	}
	if !isASCII(Text(4096, ASCII)) {
		t.Error("ASCII text contains non-ASCII bytes")
	}
	if isASCII(Text(4096, Unicode)) {
		t.Error("Unicode text is pure ASCII")
	}
	mixed := Text(4096, Mixed)
	hasLatin := strings.ContainsFunc(mixed, func(r rune) bool { return r >= 'a' && r <= 'z' })
	if isASCII(mixed) || !hasLatin {
		t.Error("Mixed text must contain both ASCII and non-ASCII words")
	}
}

// TestRunAndEachNames runs BenchmarkRunAndEachNames in a subprocess (only
// "go test -bench" builds sub-benchmark names; testing.Benchmark does not) and
// checks the generated names.
func TestRunAndEachNames(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command(executable, "-test.run=^$", "-test.bench=^BenchmarkRunAndEachNames$", "-test.benchtime=1x")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("benchmark subprocess: %v\n%s", err, output)
	}
	for _, want := range []string{"BenchmarkRunAndEachNames/type=int/n=1", "BenchmarkRunAndEachNames/type=int/n=2"} {
		if !strings.Contains(string(output), want) {
			t.Errorf("missing %q in:\n%s", want, output)
		}
	}
}

func BenchmarkRunAndEachNames(b *testing.B) {
	Each(b, TypeCase{Name: "int", Fn: func(b *testing.B) {
		Run(b, []int{1, 2}, func(b *testing.B, n int) {
			b.ReportAllocs()
			for b.Loop() {
				_ = n
			}
		})
	}})
}

func isASCII(s string) bool {
	for i := range len(s) {
		if s[i] >= utf8.RuneSelf {
			return false
		}
	}
	return true
}
