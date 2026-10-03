package sliceutil

import (
	"math"
	"slices"
	"testing"
)

func TestChunk_AppendDoesNotOverwriteInput(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6}
	chunks := Chunk(input[:5], 2)
	for i, chunk := range chunks {
		if cap(chunk) != len(chunk) {
			t.Errorf("chunk %d capacity = %d, want %d", i, cap(chunk), len(chunk))
		}

		grown := append(chunk, 99)
		if grown[len(chunk)] != 99 {
			t.Fatal("append did not add the new element")
		}
	}

	if !slices.Equal(input, []int{1, 2, 3, 4, 5, 6}) {
		t.Errorf("appending to chunks changed input to %v", input)
	}

	// Element assignments still intentionally share the input's storage.
	chunks[0][0] = 42
	if input[0] != 42 {
		t.Error("chunk elements no longer share input storage")
	}
}

func TestChunk_AcceptsMaxIntSize(t *testing.T) {
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Errorf("Chunk panicked with MaxInt size: %v", recovered)
		}
	}()

	got := Chunk([]int{1, 2}, math.MaxInt)
	if len(got) != 1 || !slices.Equal(got[0], []int{1, 2}) {
		t.Errorf("Chunk with MaxInt size = %v, want [[1 2]]", got)
	}
}
