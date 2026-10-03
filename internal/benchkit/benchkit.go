// Package benchkit provides deterministic inputs and size sweeps for the
// repository's benchmarks. Only _test.go files import it.
//
// Every generator builds a fresh pseudo-random source from a fixed seed and
// its own arguments, so a result depends only on those arguments and never on
// which benchmarks ran before it.
package benchkit

import (
	"encoding/binary"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"
)

// Sizes are the input sizes for linear and n·log n workloads.
var Sizes = []int{16, 1024, 65536}

// SmallSizes are the input sizes for quadratic or worse workloads, including
// square-matrix dimensions.
var SmallSizes = []int{8, 64, 512}

// TextKind selects the characters that Text produces.
type TextKind int

const (
	// ASCII text uses lowercase ASCII words, spaces, and full stops.
	ASCII TextKind = iota
	// Unicode text uses accented Latin (precomposed and decomposed), Greek,
	// Cyrillic, CJK, Korean, and emoji words, separated by ASCII spaces and
	// full stops.
	Unicode
	// Mixed text alternates ASCII and Unicode words.
	Mixed
)

const seed = 0x9e3779b97f4a7c15

const (
	streamInts uint64 = iota + 1
	streamFloats
	streamStrings
	streamText
	streamWords
	streamBytes
	streamMatrix
)

var asciiWords = []string{
	"the", "quick", "brown", "fox", "jumps", "over", "lazy", "dog", "cache",
	"slice", "string", "token", "value", "index", "golang", "bench", "random",
	"search", "normal", "simple", "vector", "matrix", "queue", "stack", "set",
	"tree", "hash", "key", "data", "stream",
}

var unicodeWords = []string{
	"café", "café", "naïve", "résumé", "Ærøskøbing", "straße", "Zürich",
	"ἀλήθεια", "λόγος", "Ελλάδα", "мир", "привет", "東京", "数据", "日本語",
	"キャッシュ", "한국어", "😀", "👍🏽", "ﬁne",
}

// newRand returns a fresh source for one generator call.
func newRand(stream uint64, arg int) *rand.Rand {
	return rand.New(rand.NewPCG(seed, stream<<32^uint64(arg))) //nolint:gosec // G404, G115: deterministic benchmark data, not security; arg is a non-negative size.
}

func mix(a, b int) int { return a*1_000_003 + b }

// Ints returns n pseudo-random ints in [0, 2³¹).
func Ints(n int) []int {
	r := newRand(streamInts, n)
	s := make([]int, n)
	for i := range s {
		s[i] = int(r.Int32())
	}

	return s
}

// SortedInts returns Ints(n) in ascending order.
func SortedInts(n int) []int {
	s := Ints(n)
	slices.Sort(s)

	return s
}

// Floats returns n pseudo-random float64 values in [0, 1).
func Floats(n int) []float64 {
	r := newRand(streamFloats, n)
	s := make([]float64, n)
	for i := range s {
		s[i] = r.Float64()
	}

	return s
}

// Strings returns n strings of exactly length runes, each rune drawn from
// alphabet. It panics if alphabet is empty and length > 0.
func Strings(n, length int, alphabet string) []string {
	runes := []rune(alphabet)
	if len(runes) == 0 && length > 0 {
		panic("benchkit: Strings needs a non-empty alphabet")
	}
	r := newRand(streamStrings, mix(n, length))
	out := make([]string, n)
	var sb strings.Builder
	for i := range out {
		sb.Reset()
		for range length {
			sb.WriteRune(runes[r.IntN(len(runes))])
		}
		out[i] = sb.String()
	}

	return out
}

// Words returns n lowercase ASCII words drawn from a fixed vocabulary.
func Words(n int) []string {
	r := newRand(streamWords, n)
	out := make([]string, n)
	for i := range out {
		out[i] = asciiWords[r.IntN(len(asciiWords))]
	}

	return out
}

// Text returns exactly nBytes bytes of valid UTF-8 text of the given kind:
// words separated by spaces, with a full stop every twelve words. The text is
// cut at a rune boundary and padded with ASCII spaces to the exact size.
func Text(nBytes int, kind TextKind) string {
	r := newRand(streamText, mix(nBytes, int(kind)))
	var sb strings.Builder
	sb.Grow(nBytes + 16)
	for i := 0; sb.Len() < nBytes; i++ {
		if i > 0 {
			if i%12 == 0 {
				sb.WriteString(". ")
			} else {
				sb.WriteByte(' ')
			}
		}
		words := asciiWords
		if kind == Unicode || (kind == Mixed && i%2 == 1) {
			words = unicodeWords
		}
		sb.WriteString(words[r.IntN(len(words))])
	}
	s := sb.String()
	cut := min(nBytes, len(s))
	for cut > 0 && cut < len(s) && !utf8.RuneStart(s[cut]) {
		cut--
	}

	return s[:cut] + strings.Repeat(" ", nBytes-cut)
}

// Bytes returns n pseudo-random bytes.
func Bytes(n int) []byte {
	r := newRand(streamBytes, n)
	b := make([]byte, n+7)
	for i := 0; i < n; i += 8 {
		binary.LittleEndian.PutUint64(b[i:], r.Uint64())
	}

	return b[:n:n]
}

// Matrix returns a rows×cols matrix of pseudo-random values in [0, 1).
func Matrix(rows, cols int) [][]float64 {
	r := newRand(streamMatrix, mix(rows, cols))
	m := make([][]float64, rows)
	for i := range m {
		m[i] = make([]float64, cols)
		for j := range m[i] {
			m[i][j] = r.Float64()
		}
	}

	return m
}

// Run runs fn once per size, as a sub-benchmark named "n=<size>".
func Run(b *testing.B, sizes []int, fn func(b *testing.B, n int)) {
	b.Helper()
	for _, n := range sizes {
		b.Run("n="+strconv.Itoa(n), func(b *testing.B) { fn(b, n) })
	}
}

// TypeCase is one type instantiation of a generic benchmark.
type TypeCase struct {
	// Name is the element type, e.g. "int"; the sub-benchmark is "type=<Name>".
	Name string
	// Fn runs the benchmark for that type and owns its sub-benchmarks and loops.
	Fn func(b *testing.B)
}

// Each runs every case as a sub-benchmark named "type=<Name>".
func Each(b *testing.B, cases ...TypeCase) {
	b.Helper()
	for _, c := range cases {
		b.Run("type="+c.Name, c.Fn)
	}
}
