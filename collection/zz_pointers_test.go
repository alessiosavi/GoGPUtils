package collection

import (
	"testing"
	"unsafe"
)

type myInt int

func checkPointers[T any](t *testing.T, want bool) {
	t.Helper()
	same(t, mayHavePointers[T](), want)
	n := testing.AllocsPerRun(10, func() { _ = genericPointers[T]() })
	t.Logf("mayHavePointers type=%T AllocsPerRun=%g", (*T)(nil), n)
	if n != 0 {
		t.Fatalf("helper allocates: %v", n)
	}
}

func TestMayHavePointers(t *testing.T) {
	checkPointers[int](t, false)
	checkPointers[int8](t, false)
	checkPointers[int16](t, false)
	checkPointers[int32](t, false)
	checkPointers[int64](t, false)
	checkPointers[uint](t, false)
	checkPointers[uint8](t, false)
	checkPointers[uint16](t, false)
	checkPointers[uint32](t, false)
	checkPointers[uint64](t, false)
	checkPointers[uintptr](t, false)
	checkPointers[float32](t, false)
	checkPointers[float64](t, false)
	checkPointers[complex64](t, false)
	checkPointers[complex128](t, false)
	checkPointers[bool](t, false)
	checkPointers[string](t, true)
	checkPointers[*int](t, true)
	checkPointers[[]int](t, true)
	checkPointers[any](t, true)
	checkPointers[struct{ p *int }](t, true)
	checkPointers[myInt](t, true)
	checkPointers[[4096]int](t, true)
	checkPointers[map[int]int](t, true)
	checkPointers[func()](t, true)
	checkPointers[chan int](t, true)
	checkPointers[unsafe.Pointer](t, true)
}

//go:noinline
func genericPointers[T any]() bool { return mayHavePointers[T]() }

func TestPointerDictionary(t *testing.T) {
	same(t, genericPointers[int](), false)
	same(t, genericPointers[myInt](), true)
	same(t, genericPointers[*int](), true)
	same(t, genericPointers[*string](), true)
	same(t, genericPointers[any](), true)
	same(t, genericPointers[[4096]int](), true)
	if a := testing.AllocsPerRun(100, func() { _ = genericPointers[[4096]int]() }); a != 0 {
		t.Fatal(a)
	}
}
