package textnorm

import (
	"runtime"
	"strings"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

// Mixed case is made explicit here because benchkit.ASCII is lowercase-only:
// uppercase the first letter and then every other letter, preserving byte length.
// late-nonascii replaces the final two ASCII bytes with the two-byte rune é.
func zzFastPathInput(n int, kind string) string {
	if kind == "lower" {
		return benchkit.Text(n, benchkit.ASCII)
	}
	text := []byte(benchkit.Text(n, benchkit.ASCII))
	letter := 0
	for i, c := range text {
		if c >= 'a' && c <= 'z' {
			if letter%2 == 0 {
				text[i] = c - ('a' - 'A')
			}
			letter++
		}
	}
	if kind == "late-nonascii" {
		copy(text[n-2:], "é")
	}
	return string(text)
}
func zzBenchmarkFastPath(b *testing.B, p Pipeline, kinds ...string) {
	for _, kind := range kinds {
		b.Run("case="+kind, func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				in := zzFastPathInput(n, kind)
				b.ReportAllocs()
				for b.Loop() {
					if _, err := p.Run(in); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}
func BenchmarkFastPath_FoldCase(b *testing.B) {
	zzBenchmarkFastPath(b, New().FoldCase(), "ascii", "lower", "late-nonascii")
}
func BenchmarkFastPath_Lower(b *testing.B) {
	zzBenchmarkFastPath(b, New().Lower(), "ascii", "lower", "late-nonascii")
}
func BenchmarkFastPath_NormalizeUnicode(b *testing.B) {
	zzBenchmarkFastPath(b, New().NormalizeUnicode(), "ascii", "late-nonascii")
}
func BenchmarkFastPath_NormalizeUnicodeLatin(b *testing.B) {
	zzBenchmarkFastPath(b, New().NormalizeUnicodeLatin(), "ascii", "late-nonascii")
}

// This is a fixture check, not a benchmark execution.
func TestZZFastPathFixtures(t *testing.T) {
	for _, n := range benchkit.Sizes {
		original := benchkit.Text(n, benchkit.ASCII)
		if original != strings.ToLower(original) {
			t.Fatal("benchkit.ASCII premise changed")
		}
		for _, kind := range []string{"ascii", "lower", "late-nonascii"} {
			s := zzFastPathInput(n, kind)
			if kind == "lower" {
				if s != original || len(s) != n || !zzAllASCII(s) {
					t.Fatalf("lower fixture n=%d differs from benchkit.ASCII", n)
				}
				continue
			}
			if len(s) != n || strings.ToLower(s) == s {
				t.Fatalf("fixture %s n=%d has wrong length or lacks uppercase", kind, n)
			}
			if kind == "late-nonascii" && (!strings.HasSuffix(s, "é") || !zzAllASCII(s[:n-2])) {
				t.Fatal("non-ASCII is not exclusively at end")
			}
			if kind == "ascii" && !zzAllASCII(s) {
				t.Fatal("non-ASCII in ASCII fixture")
			}
			letter := 0
			for i := range n {
				if kind == "late-nonascii" && i >= n-2 {
					break
				}
				want := original[i]
				if want >= 'a' && want <= 'z' {
					if letter%2 == 0 {
						want -= 'a' - 'A'
					}
					letter++
				}
				if s[i] != want {
					t.Fatalf("%s n=%d byte=%d recipe mismatch", kind, n, i)
				}
			}
		}
	}
}
func zzAllASCII(s string) bool {
	for i := range len(s) {
		if s[i] >= 128 {
			return false
		}
	}
	return true
}

func TestZZMatrixAllocationDiagnostics(t *testing.T) {
	for _, n := range benchkit.Sizes {
		for _, kind := range []string{"ascii", "lower"} {
			in := zzFastPathInput(n, kind)
			for _, p := range []struct {
				name     string
				pipeline Pipeline
			}{
				{"FoldCase", New().FoldCase()}, {"Lower", New().Lower()},
				{"NormalizeUnicode", New().NormalizeUnicode()}, {"NormalizeUnicodeLatin", New().NormalizeUnicodeLatin()},
			} {
				if kind == "lower" && p.name != "FoldCase" && p.name != "Lower" {
					continue
				}
				var sink string
				allocs := testing.AllocsPerRun(100, func() {
					var err error
					sink, err = p.pipeline.Run(in)
					if err != nil {
						panic(err)
					}
				})
				runtime.KeepAlive(sink)
				target := 1.0
				if kind == "lower" || (p.name == "NormalizeUnicode" && n == 16) {
					target = 0
				}
				t.Logf("%s case=%s n=%d allocs/run=%g target<=%g meets=%v", p.name, kind, n, allocs, target, allocs <= target)
			}
		}
	}
}
