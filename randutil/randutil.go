package randutil

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"math/bits"
	"math/rand/v2"
)

// Character sets for random string generation.
const (
	Digits       = "0123456789"
	Lowercase    = "abcdefghijklmnopqrstuvwxyz"
	Uppercase    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Letters      = Lowercase + Uppercase
	AlphaNumeric = Letters + Digits
	Hex          = "0123456789abcdef"
	Symbols      = "!@#$%^&*()_+-=[]{}|;:,.<>?"
	All          = AlphaNumeric + Symbols
)

// Common errors.
var (
	ErrInvalidLength = errors.New("length must be positive")
	ErrEmptyCharset  = errors.New("charset cannot be empty")
	ErrEmptySlice    = errors.New("slice cannot be empty")
)

// ============================================================================
// Secure Random Generation (crypto/rand)
// ============================================================================

// SecureBytes returns n cryptographically secure random bytes.
//
// Example:
//
//	key, err := SecureBytes(32) // 32 random bytes
func SecureBytes(n int) ([]byte, error) {
	if n <= 0 {
		return nil, ErrInvalidLength
	}

	b := make([]byte, n)
	if _, err := io.ReadFull(cryptorand.Reader, b); err != nil {
		return nil, err
	}

	return b, nil
}

// SecureString returns a cryptographically secure random string
// of the specified byte length. The charset is treated as a set of bytes;
// use ASCII charsets. Every byte position is selected with equal probability,
// so repeated bytes in charset have proportionally greater probability.
//
// Reads may be batched. Output, errors, and consumed bytes match repeated
// crypto/rand.Int draws with bound len(charset) over an identical byte stream
// for readers whose behavior depends only on stream position. This excludes
// readers whose output or errors depend on Read-call granularity, readers that
// return data with a non-EOF error and later deliver more data, and readers that
// replace crypto/rand.Reader during the call; the reader is captured once per
// batch.
//
// Example:
//
//	password, err := SecureString(16, AlphaNumeric)
//	token, err := SecureString(32, Hex)
func SecureString(length int, charset string) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}

	if len(charset) == 0 {
		return "", ErrEmptyCharset
	}

	result := make([]byte, length)
	if len(charset) == 1 {
		for i := range result {
			result[i] = charset[0]
		}
		return string(result), nil
	}
	k, mask := candidateParams(uint64(len(charset)))
	var scratch []byte
	if k >= 2 {
		scratch = make([]byte, min(length, 4096/k)*k)
	}
	filled := 0
	for filled < length {
		remaining := length - filled
		m := min(remaining, 4096/k)
		var candidates []byte
		if k == 1 {
			candidates = result[filled : filled+m]
		} else {
			candidates = scratch[:m*k]
		}
		if err := fillSecureCandidates(cryptorand.Reader, candidates, k); err != nil {
			return "", err
		}
		for i := 0; i < len(candidates); i += k {
			v := uint64(candidates[i] & mask)
			for j := 1; j < k; j++ {
				v = v<<8 | uint64(candidates[i+j])
			}
			if v < uint64(len(charset)) {
				result[filled] = charset[v]
				filled++
			}
		}
	}
	return string(result), nil
}

// fillSecureCandidates keeps the underlying error distinct from ReadFull's
// synthesized UnexpectedEOF while preserving candidate-sized error handling.
func fillSecureCandidates(r io.Reader, p []byte, k int) error {
	n := 0
	for n < len(p) {
		nr, err := r.Read(p[n:])
		n += nr
		if n == len(p) {
			return nil
		}
		if err != nil {
			// Like a candidate-sized ReadFull, ignore an error accompanying
			// a completed candidate, then ask the same reader for more.
			if nr > 0 && n%k == 0 {
				continue
			}
			//nolint:errorlint // Only exact EOF is translated by io.ReadFull.
			if err == io.EOF && n%k != 0 {
				return io.ErrUnexpectedEOF
			}
			return err
		}
	}
	return nil
}

// SecureInt returns a uniformly distributed cryptographically secure random int in [0, max).
// Returns ErrInvalidLength if max <= 0.
func SecureInt(max int) (int, error) {
	if max <= 0 {
		return 0, ErrInvalidLength
	}

	n, err := secureUint64n(uint64(max))
	//nolint:gosec // n is zero on error or less than the positive int bound max.
	return int(n), err
}

// SecureInt64 returns a uniformly distributed cryptographically secure random int64 in [0, max).
// Returns ErrInvalidLength if max <= 0.
func SecureInt64(max int64) (int64, error) {
	if max <= 0 {
		return 0, ErrInvalidLength
	}

	n, err := secureUint64n(uint64(max))
	//nolint:gosec // n is zero on error or less than the positive int64 bound max.
	return int64(n), err
}

// candidateParams requires bound > 1.
func candidateParams(bound uint64) (k int, mask byte) {
	bitLen := bits.Len64(bound - 1)
	k = (bitLen + 7) / 8
	top := bitLen - 8*(k-1)
	return k, 0xff >> (8 - top)
}

func secureUint64n(bound uint64) (uint64, error) {
	if bound == 1 {
		return 0, nil
	}
	k, mask := candidateParams(bound)
	var buf [8]byte
	r := cryptorand.Reader
	for {
		if _, err := io.ReadFull(r, buf[:k]); err != nil {
			return 0, err
		}
		buf[0] &= mask
		var v uint64
		for _, b := range buf[:k] {
			v = v<<8 | uint64(b)
		}
		if v < bound {
			return v, nil
		}
	}
}

// SecureID generates a cryptographically secure random ID.
// Returns a 32-character hex string (128 bits of entropy).
//
// Example:
//
//	id := SecureID() // "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
func SecureID() (string, error) {
	b := make([]byte, 16)
	if _, err := io.ReadFull(cryptorand.Reader, b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

// SecureChoice returns a cryptographically secure random element from the slice.
func SecureChoice[T any](s []T) (T, error) {
	if len(s) == 0 {
		var zero T

		return zero, ErrEmptySlice
	}

	idx, err := SecureInt(len(s))
	if err != nil {
		var zero T

		return zero, err
	}

	return s[idx], nil
}

// ============================================================================
// Fast Random Generation (math/rand/v2)
// ============================================================================

// Generator provides fast (non-cryptographic) random generation.
// Use NewGenerator() for automatic seeding or NewGeneratorWithSeed() for reproducibility.
type Generator struct {
	rng *rand.Rand
}

// NewGenerator creates a new Generator with automatic seeding.
// Uses a cryptographically secure seed for good randomness.
func NewGenerator() *Generator {
	return &Generator{
		rng: rand.New(rand.NewPCG(secureUint64(), secureUint64())),
	}
}

// NewGeneratorWithSeed creates a Generator with a specific seed.
// Useful for reproducible random sequences in tests.
//
// Example:
//
//	rng := NewGeneratorWithSeed(42)
//	// Always produces the same sequence
func NewGeneratorWithSeed(seed uint64) *Generator {
	return &Generator{
		rng: rand.New(rand.NewPCG(seed, seed)),
	}
}

// secureUint64 generates a random uint64 using crypto/rand.
func secureUint64() uint64 {
	var b [8]byte
	if _, err := io.ReadFull(cryptorand.Reader, b[:]); err != nil {
		// Fallback to less random but functional seed
		return uint64(0xDEADBEEF)
	}

	return binary.BigEndian.Uint64(b[:])
}

// Int returns a random int in [0, max).
// Panics if max <= 0.
func (g *Generator) Int(max int) int {
	return g.rng.IntN(max)
}

// IntRange returns a random int in [min, max].
// Panics if min > max.
func (g *Generator) IntRange(min, max int) int {
	if min == max {
		return min
	}

	return min + g.rng.IntN(max-min+1)
}

// Int64 returns a random int64 in [0, max).
func (g *Generator) Int64(max int64) int64 {
	return g.rng.Int64N(max)
}

// Int64Range returns a random int64 in [min, max].
func (g *Generator) Int64Range(min, max int64) int64 {
	if min == max {
		return min
	}

	return min + g.rng.Int64N(max-min+1)
}

// Float64 returns a random float64 in [0.0, 1.0).
func (g *Generator) Float64() float64 {
	return g.rng.Float64()
}

// Float64Range returns a random float64 in [min, max).
func (g *Generator) Float64Range(min, max float64) float64 {
	return min + g.rng.Float64()*(max-min)
}

// Float32 returns a random float32 in [0.0, 1.0).
func (g *Generator) Float32() float32 {
	return g.rng.Float32()
}

// Bool returns a random boolean.
func (g *Generator) Bool() bool {
	return g.rng.IntN(2) == 1
}

// Bytes returns n random bytes.
func (g *Generator) Bytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(g.rng.IntN(256))
	}

	return b
}

// String returns a random string of the specified length using the charset.
func (g *Generator) String(length int, charset string) string {
	if length <= 0 || len(charset) == 0 {
		return ""
	}

	result := make([]byte, length)
	for i := range result {
		result[i] = charset[g.rng.IntN(len(charset))]
	}

	return string(result)
}

// AlphaNumericString returns a random alphanumeric string.
func (g *Generator) AlphaNumericString(length int) string {
	return g.String(length, AlphaNumeric)
}

// ChoiceInt returns a random element from an int slice.
// Panics if slice is empty.
func (g *Generator) ChoiceInt(s []int) int {
	return s[g.rng.IntN(len(s))]
}

// ChoiceString returns a random element from a string slice.
// Panics if slice is empty.
func (g *Generator) ChoiceString(s []string) string {
	return s[g.rng.IntN(len(s))]
}

// ShuffleInts shuffles an int slice in place.
func (g *Generator) ShuffleInts(s []int) {
	for i := len(s) - 1; i > 0; i-- {
		j := g.rng.IntN(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

// ShuffleStrings shuffles a string slice in place.
func (g *Generator) ShuffleStrings(s []string) {
	for i := len(s) - 1; i > 0; i-- {
		j := g.rng.IntN(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

// SampleInts returns n random ints from the slice (without replacement).
// If n > len(s), returns all elements shuffled.
func (g *Generator) SampleInts(s []int, n int) []int {
	if len(s) == 0 || n <= 0 {
		return nil
	}

	result := make([]int, len(s))
	copy(result, s)
	g.ShuffleInts(result)

	if n >= len(s) {
		return result
	}

	return result[:n]
}

// Choice is a generic function for random selection from a slice.
func Choice[T any](g *Generator, s []T) T {
	return s[g.rng.IntN(len(s))]
}

// ChoiceN is a generic function returning n random elements (with replacement).
func ChoiceN[T any](g *Generator, s []T, n int) []T {
	if len(s) == 0 || n <= 0 {
		return nil
	}

	result := make([]T, n)
	for i := range result {
		result[i] = s[g.rng.IntN(len(s))]
	}

	return result
}

// Sample is a generic function returning n random elements (without replacement).
func Sample[T any](g *Generator, s []T, n int) []T {
	if len(s) == 0 || n <= 0 {
		return nil
	}

	result := make([]T, len(s))
	copy(result, s)
	Shuffle(g, result)

	if n >= len(s) {
		return result
	}

	return result[:n]
}

// Shuffle is a generic function that shuffles a slice in place.
func Shuffle[T any](g *Generator, s []T) {
	for i := len(s) - 1; i > 0; i-- {
		j := g.rng.IntN(i + 1)
		s[i], s[j] = s[j], s[i]
	}
}

// ShuffleCopy returns a new shuffled slice (does not modify original).
func ShuffleCopy[T any](g *Generator, s []T) []T {
	result := make([]T, len(s))
	copy(result, s)
	Shuffle(g, result)

	return result
}

// WeightedChoice returns a random index based on weights.
// Weight values don't need to sum to 1.0.
func (g *Generator) WeightedChoice(weights []float64) int {
	if len(weights) == 0 {
		return -1
	}

	// Calculate total weight
	var total float64
	for _, w := range weights {
		total += w
	}

	// Pick random point
	r := g.Float64() * total

	// Find the bucket
	var cumulative float64
	for i, w := range weights {
		cumulative += w
		if r < cumulative {
			return i
		}
	}

	return len(weights) - 1
}

// Probability returns true with the given probability (0.0 to 1.0).
//
// Example:
//
//	if rng.Probability(0.75) {
//	    // 75% chance of executing
//	}
func (g *Generator) Probability(p float64) bool {
	return g.Float64() < p
}

// ============================================================================
// Sequence Generation
// ============================================================================

// Sequence returns a slice of n sequential integers starting from start.
func Sequence(start, n int) []int {
	if n <= 0 {
		return nil
	}

	result := make([]int, n)
	for i := range result {
		result[i] = start + i
	}

	return result
}

// Range returns a slice of integers from start to end (exclusive).
func Range(start, end int) []int {
	if end <= start {
		return nil
	}

	return Sequence(start, end-start)
}

// RangeStep returns a slice of integers from start to end (exclusive) with step.
// Stops before an addition would overflow.
func RangeStep(start, end, step int) []int {
	if step == 0 || (step > 0 && start >= end) || (step < 0 && start <= end) {
		return nil
	}

	var result []int
	//nolint:gosec // Unsigned subtraction intentionally handles signed extremes.
	distance, stride := uint64(end)-uint64(start), uint64(step)
	if step < 0 {
		//nolint:gosec // Modular conversion preserves the descending distance.
		distance = uint64(start) - uint64(end)
		stride = -stride
	}
	count := (distance-1)/stride + 1
	if count <= uint64(^uint(0)>>1) {
		result = make([]int, 0, int(count))
	}

	for i := start; (step > 0 && i < end) || (step < 0 && i > end); {
		result = append(result, i)

		next := i + step
		if (step > 0 && next <= i) || (step < 0 && next >= i) {
			break
		}

		i = next
	}

	return result
}
