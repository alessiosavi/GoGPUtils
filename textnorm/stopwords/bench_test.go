package stopwords

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

var benchmarkSet map[string]struct{}

// BenchmarkEnglish measures cached access to the English set. The first call
// parses the embedded list during setup; the shared map is never mutated.
func BenchmarkEnglish(b *testing.B) {
	English()
	b.ReportAllocs()
	for b.Loop() {
		benchmarkSet = English()
	}
}

// BenchmarkFrench measures cached access to the French set. The first call
// parses the embedded list during setup; the shared map is never mutated.
func BenchmarkFrench(b *testing.B) {
	French()
	b.ReportAllocs()
	for b.Loop() {
		benchmarkSet = French()
	}
}

// BenchmarkItalian measures cached access to the Italian set. The first call
// parses the embedded list during setup; the shared map is never mutated.
func BenchmarkItalian(b *testing.B) {
	Italian()
	b.ReportAllocs()
	for b.Loop() {
		benchmarkSet = Italian()
	}
}

// BenchmarkLoadFromFile reads and parses 500 distinct lowercase words from
// benchkit, with no duplicates; the resulting set has 500 keys. File creation in
// b.TempDir is setup; opening, reading, parsing and closing the same file are
// measured each iteration. The informational language is "english"; the file
// is not mutated and no restoration is measured.
func BenchmarkLoadFromFile(b *testing.B) {
	b.Run("words=500", func(b *testing.B) {
		words := benchkit.Strings(500, 8, "abcdefghijklmnopqrstuvwxyz")
		distinct := make(map[string]struct{}, len(words))
		for _, word := range words {
			distinct[word] = struct{}{}
		}
		if len(distinct) != 500 {
			b.Fatalf("setup has %d distinct words, want 500", len(distinct))
		}
		path := filepath.Join(b.TempDir(), "english.txt")
		if err := os.WriteFile(path, []byte(strings.Join(words, "\n")+"\n"), 0o600); err != nil {
			b.Fatal(err)
		}
		b.ReportAllocs()
		for b.Loop() {
			set, err := LoadFromFile("english", path)
			if err != nil {
				b.Fatal(err)
			}
			benchmarkSet = set
		}
	})
}

// BenchmarkLoadFromList constructs a 500-key set from 500 distinct lowercase
// words from benchkit, with no duplicates. The informational language is "english";
// the input slice is immutable and no restoration is measured.
func BenchmarkLoadFromList(b *testing.B) {
	b.Run("words=500", func(b *testing.B) {
		words := benchkit.Strings(500, 8, "abcdefghijklmnopqrstuvwxyz")
		distinct := make(map[string]struct{}, len(words))
		for _, word := range words {
			distinct[word] = struct{}{}
		}
		if len(distinct) != 500 {
			b.Fatalf("setup has %d distinct words, want 500", len(distinct))
		}
		b.ReportAllocs()
		for b.Loop() {
			benchmarkSet = LoadFromList("english", words)
		}
	})
}

// BenchmarkCleanAllStopwords builds a language-set union, not cleaned text.
// languages=1/2/3 selects successive prefixes of english, french, italian.
// All built-in sets are warmed before timing and never mutated. Each
// iteration constructs a new map; no restoration is measured.
func BenchmarkCleanAllStopwords(b *testing.B) {
	English()
	French()
	Italian()
	languages := []string{"english", "french", "italian"}
	for count := 1; count <= len(languages); count++ {
		b.Run("languages="+strconv.Itoa(count), func(b *testing.B) {
			selected := languages[:count]
			b.ReportAllocs()
			for b.Loop() {
				benchmarkSet = CleanAllStopwords(selected)
			}
		})
	}
}

// BenchmarkUnion builds a fresh union of the three built-in language sets.
// Embedded lists are loaded during setup; the input maps remain unchanged
// and no restoration is measured.
func BenchmarkUnion(b *testing.B) {
	b.Run("languages=3", func(b *testing.B) {
		sets := []map[string]struct{}{English(), French(), Italian()}
		b.ReportAllocs()
		for b.Loop() {
			benchmarkSet = Union(sets...)
		}
	})
}
