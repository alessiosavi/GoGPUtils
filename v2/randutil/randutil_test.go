package randutil

import (
	"slices"
	"testing"
)

// ============================================================================
// Secure Random Tests
// ============================================================================

func TestSecureBytes(t *testing.T) {
	sizes := []int{1, 16, 32, 64, 128}

	for _, size := range sizes {
		b, err := SecureBytes(size)
		if err != nil {
			t.Errorf("SecureBytes(%d) error: %v", size, err)
		}
		if len(b) != size {
			t.Errorf("SecureBytes(%d) returned %d bytes", size, len(b))
		}
	}

	// Two calls should produce different output
	b1, _ := SecureBytes(32)
	b2, _ := SecureBytes(32)
	if slices.Equal(b1, b2) {
		t.Error("SecureBytes should produce different output")
	}

	// Invalid size
	_, err := SecureBytes(0)
	if err != ErrInvalidLength {
		t.Errorf("SecureBytes(0) error = %v, want ErrInvalidLength", err)
	}
}

func TestSecureString(t *testing.T) {
	tests := []struct {
		length  int
		charset string
	}{
		{8, Digits},
		{16, Lowercase},
		{32, AlphaNumeric},
		{64, Hex},
	}

	for _, tt := range tests {
		s, err := SecureString(tt.length, tt.charset)
		if err != nil {
			t.Errorf("SecureString(%d, %q) error: %v", tt.length, tt.charset, err)
		}
		if len(s) != tt.length {
			t.Errorf("SecureString(%d, %q) returned len %d", tt.length, tt.charset, len(s))
		}

		// Verify all characters are in charset
		for _, c := range s {
			if !containsRune(tt.charset, c) {
				t.Errorf("SecureString produced char %q not in charset", c)
			}
		}
	}

	// Invalid inputs
	_, err := SecureString(0, AlphaNumeric)
	if err != ErrInvalidLength {
		t.Errorf("SecureString(0, ...) error = %v, want ErrInvalidLength", err)
	}

	_, err = SecureString(10, "")
	if err != ErrEmptyCharset {
		t.Errorf("SecureString(10, '') error = %v, want ErrEmptyCharset", err)
	}
}

func TestSecureInt(t *testing.T) {
	max := 100

	for i := 0; i < 1000; i++ {
		n, err := SecureInt(max)
		if err != nil {
			t.Fatalf("SecureInt error: %v", err)
		}
		if n < 0 || n >= max {
			t.Errorf("SecureInt(%d) = %d, out of range", max, n)
		}
	}

	// Distribution test (rough check)
	counts := make(map[int]int)
	for i := 0; i < 10000; i++ {
		n, _ := SecureInt(10)
		counts[n]++
	}

	// Each bucket should have roughly 1000 hits (±300)
	for i := 0; i < 10; i++ {
		if counts[i] < 700 || counts[i] > 1300 {
			t.Logf("Warning: SecureInt distribution may be biased: bucket %d has %d hits", i, counts[i])
		}
	}
}

func TestSecureID(t *testing.T) {
	id, err := SecureID()
	if err != nil {
		t.Fatalf("SecureID error: %v", err)
	}

	if len(id) != 32 {
		t.Errorf("SecureID length = %d, want 32", len(id))
	}

	// All characters should be hex
	for _, c := range id {
		if !containsRune(Hex, c) {
			t.Errorf("SecureID contains non-hex char: %q", c)
		}
	}

	// Two IDs should be different
	id2, _ := SecureID()
	if id == id2 {
		t.Error("SecureID should produce unique IDs")
	}
}

func TestSecureChoice(t *testing.T) {
	items := []string{"a", "b", "c", "d", "e"}

	for i := 0; i < 100; i++ {
		choice, err := SecureChoice(items)
		if err != nil {
			t.Fatalf("SecureChoice error: %v", err)
		}
		if !slices.Contains(items, choice) {
			t.Errorf("SecureChoice returned %q, not in items", choice)
		}
	}

	// Empty slice
	_, err := SecureChoice([]int{})
	if err != ErrEmptySlice {
		t.Errorf("SecureChoice(empty) error = %v, want ErrEmptySlice", err)
	}
}

// ============================================================================
// Generator Tests
// ============================================================================

func TestGenerator_Int(t *testing.T) {
	rng := NewGenerator()

	for i := 0; i < 1000; i++ {
		n := rng.Int(100)
		if n < 0 || n >= 100 {
			t.Errorf("Int(100) = %d, out of range", n)
		}
	}
}

func TestGenerator_IntRange(t *testing.T) {
	rng := NewGenerator()

	for i := 0; i < 1000; i++ {
		n := rng.IntRange(10, 20)
		if n < 10 || n > 20 {
			t.Errorf("IntRange(10, 20) = %d, out of range", n)
		}
	}

	// Same min and max
	if rng.IntRange(5, 5) != 5 {
		t.Error("IntRange(5, 5) should return 5")
	}
}

func TestGenerator_Float64(t *testing.T) {
	rng := NewGenerator()

	for i := 0; i < 1000; i++ {
		f := rng.Float64()
		if f < 0 || f >= 1 {
			t.Errorf("Float64() = %v, out of range [0, 1)", f)
		}
	}
}

func TestGenerator_Float64Range(t *testing.T) {
	rng := NewGenerator()

	for i := 0; i < 1000; i++ {
		f := rng.Float64Range(10.0, 20.0)
		if f < 10.0 || f >= 20.0 {
			t.Errorf("Float64Range(10, 20) = %v, out of range", f)
		}
	}
}

func TestGenerator_Bool(t *testing.T) {
	rng := NewGenerator()

	trueCount := 0
	for i := 0; i < 10000; i++ {
		if rng.Bool() {
			trueCount++
		}
	}

	// Should be roughly 50% (±5%)
	if trueCount < 4500 || trueCount > 5500 {
		t.Logf("Warning: Bool distribution may be biased: %d true out of 10000", trueCount)
	}
}

func TestGenerator_String(t *testing.T) {
	rng := NewGenerator()

	s := rng.String(20, AlphaNumeric)
	if len(s) != 20 {
		t.Errorf("String length = %d, want 20", len(s))
	}

	for _, c := range s {
		if !containsRune(AlphaNumeric, c) {
			t.Errorf("String contains invalid char: %q", c)
		}
	}
}

func TestGenerator_Choice(t *testing.T) {
	rng := NewGenerator()
	items := []int{1, 2, 3, 4, 5}

	for i := 0; i < 100; i++ {
		choice := rng.ChoiceInt(items)
		if !slices.Contains(items, choice) {
			t.Errorf("ChoiceInt returned %d, not in items", choice)
		}
	}
}

func TestGenerator_Sample(t *testing.T) {
	rng := NewGenerator()
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	sample := rng.SampleInts(items, 5)
	if len(sample) != 5 {
		t.Errorf("Sample length = %d, want 5", len(sample))
	}

	// All items should be from original
	for _, v := range sample {
		if !slices.Contains(items, v) {
			t.Errorf("Sample contains %d, not in original", v)
		}
	}

	// No duplicates
	seen := make(map[int]bool)
	for _, v := range sample {
		if seen[v] {
			t.Error("Sample contains duplicates")
		}
		seen[v] = true
	}

	// Request more than available
	fullSample := rng.SampleInts(items, 20)
	if len(fullSample) != 10 {
		t.Errorf("Sample(20) on 10 items returned %d items", len(fullSample))
	}
}

func TestGenerator_Shuffle(t *testing.T) {
	rng := NewGenerator()
	original := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	shuffled := ShuffleCopy(rng, original)

	// Same length
	if len(shuffled) != len(original) {
		t.Errorf("Shuffle length = %d, want %d", len(shuffled), len(original))
	}

	// Same elements
	slices.Sort(shuffled)
	if !slices.Equal(shuffled, original) {
		t.Error("Shuffle changed elements")
	}

	// Original unchanged
	if !slices.Equal(original, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}) {
		t.Error("Shuffle modified original slice")
	}
}

func TestGenerator_ShuffleSlice(t *testing.T) {
	rng := NewGenerator()
	items := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	originalCopy := slices.Clone(items)

	rng.ShuffleInts(items)

	// Same elements
	slices.Sort(items)
	slices.Sort(originalCopy)
	if !slices.Equal(items, originalCopy) {
		t.Error("ShuffleInts changed elements")
	}
}

func TestGenerator_WeightedChoice(t *testing.T) {
	rng := NewGenerator()

	// Heavily weighted towards first option
	weights := []float64{0.9, 0.05, 0.05}
	counts := make([]int, 3)

	for i := 0; i < 10000; i++ {
		idx := rng.WeightedChoice(weights)
		counts[idx]++
	}

	// First option should be chosen ~90% of the time
	if counts[0] < 8500 || counts[0] > 9500 {
		t.Logf("Warning: WeightedChoice distribution may be biased: %v", counts)
	}
}

func TestGenerator_Probability(t *testing.T) {
	rng := NewGenerator()

	// 0% probability
	for i := 0; i < 100; i++ {
		if rng.Probability(0) {
			t.Error("Probability(0) should always return false")
		}
	}

	// 100% probability
	for i := 0; i < 100; i++ {
		if !rng.Probability(1) {
			t.Error("Probability(1) should always return true")
		}
	}

	// 50% probability (rough check)
	trueCount := 0
	for i := 0; i < 10000; i++ {
		if rng.Probability(0.5) {
			trueCount++
		}
	}

	if trueCount < 4500 || trueCount > 5500 {
		t.Logf("Warning: Probability(0.5) distribution: %d true out of 10000", trueCount)
	}
}

func TestGeneratorDeterministic(t *testing.T) {
	seed := uint64(42)

	rng1 := NewGeneratorWithSeed(seed)
	rng2 := NewGeneratorWithSeed(seed)

	for i := 0; i < 100; i++ {
		n1 := rng1.Int(1000)
		n2 := rng2.Int(1000)

		if n1 != n2 {
			t.Errorf("Same seed produced different results: %d vs %d", n1, n2)
		}
	}
}

// ============================================================================
// Sequence Generation Tests
// ============================================================================

func TestSequence(t *testing.T) {
	tests := []struct {
		start, n int
		want     []int
	}{
		{0, 5, []int{0, 1, 2, 3, 4}},
		{5, 3, []int{5, 6, 7}},
		{-2, 4, []int{-2, -1, 0, 1}},
		{0, 0, nil},
		{0, -1, nil},
	}

	for _, tt := range tests {
		got := Sequence(tt.start, tt.n)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Sequence(%d, %d) = %v, want %v", tt.start, tt.n, got, tt.want)
		}
	}
}

func TestRange(t *testing.T) {
	tests := []struct {
		start, end int
		want       []int
	}{
		{0, 5, []int{0, 1, 2, 3, 4}},
		{3, 7, []int{3, 4, 5, 6}},
		{5, 5, nil},
		{5, 3, nil},
	}

	for _, tt := range tests {
		got := Range(tt.start, tt.end)
		if !slices.Equal(got, tt.want) {
			t.Errorf("Range(%d, %d) = %v, want %v", tt.start, tt.end, got, tt.want)
		}
	}
}

func TestRangeStep(t *testing.T) {
	tests := []struct {
		start, end, step int
		want             []int
	}{
		{0, 10, 2, []int{0, 2, 4, 6, 8}},
		{0, 10, 3, []int{0, 3, 6, 9}},
		{10, 0, -2, []int{10, 8, 6, 4, 2}},
		{0, 10, 0, nil},
		{10, 0, 1, nil},
	}

	for _, tt := range tests {
		got := RangeStep(tt.start, tt.end, tt.step)
		if !slices.Equal(got, tt.want) {
			t.Errorf("RangeStep(%d, %d, %d) = %v, want %v", tt.start, tt.end, tt.step, got, tt.want)
		}
	}
}

// ============================================================================
// Helpers
// ============================================================================

func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkSecureBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SecureBytes(32)
	}
}

func BenchmarkSecureString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SecureString(32, AlphaNumeric)
	}
}

func BenchmarkGenerator_Int(b *testing.B) {
	rng := NewGenerator()
	for i := 0; i < b.N; i++ {
		rng.Int(1000)
	}
}

func BenchmarkGenerator_String(b *testing.B) {
	rng := NewGenerator()
	for i := 0; i < b.N; i++ {
		rng.String(32, AlphaNumeric)
	}
}

func BenchmarkGenerator_Shuffle(b *testing.B) {
	rng := NewGenerator()
	items := make([]int, 1000)
	for i := range items {
		items[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rng.ShuffleInts(items)
	}
}
