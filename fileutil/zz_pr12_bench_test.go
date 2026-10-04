package fileutil

import (
	"bytes"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

// benchmarkPR12LineBreaks places a terminator at the end of each eight-byte
// block, cycling through terms. n is the total byte length, including breaks.
func benchmarkPR12LineBreaks(n int, terms ...string) []byte {
	data := []byte(benchkit.Text(n, benchkit.ASCII))
	for i := 7; i < len(data); i += 8 {
		term := terms[(i/8)%len(terms)]
		copy(data[i+1-len(term):i+1], term)
	}

	return data
}

// BenchmarkPR12NormalizeLineTerminators measures normalization of n ASCII bytes
// with one break every eight bytes. mixed cycles CRLF, CR, LF; input is immutable
// and all setup is outside the loop. LF input also exercises owned-copy paths.
func BenchmarkPR12NormalizeLineTerminators(b *testing.B) {
	for _, target := range []struct {
		name string
		term LineTerminator
	}{{"LF", LF}, {"CR", CR}} {
		b.Run("target="+target.name, func(b *testing.B) {
			for _, input := range []struct {
				name  string
				terms []string
			}{{"lf", []string{"\n"}}, {"cr", []string{"\r"}}, {"mixed", []string{"\r\n", "\r", "\n"}}} {
				b.Run("input="+input.name, func(b *testing.B) {
					benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
						data := benchmarkPR12LineBreaks(n, input.terms...)
						b.ReportAllocs()
						for b.Loop() {
							NormalizeLineTerminators(data, target.term)
						}
					})
				})
			}
		})
	}
}

// BenchmarkPR12DetectLineTerminator measures detection in n immutable ASCII
// bytes. lf has LF every eight bytes and no CR; mixed starts with LF at byte 7
// then CRLF every eight bytes, so the second break decides Mixed. Setup is outside
// the loop and no restoration is needed.
func BenchmarkPR12DetectLineTerminator(b *testing.B) {
	for _, input := range []string{"lf", "mixed"} {
		b.Run("input="+input, func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				data := benchmarkPR12LineBreaks(n, "\n")
				want := LF
				if input == "mixed" {
					for i := 14; i+1 < len(data); i += 8 {
						data[i], data[i+1] = '\r', '\n'
					}
					want = Mixed
				}
				b.ReportAllocs()
				for b.Loop() {
					if got := DetectLineTerminator(data); got != want {
						b.Fatalf("terminator = %v, want %v", got, want)
					}
				}
			})
		})
	}
}

// benchmarkPR12SparseLineBreaks builds 80-byte lines of ASCII letters ending in
// term, cut to n bytes. An empty term produces n letters without line breaks.
func benchmarkPR12SparseLineBreaks(n int, term string) []byte {
	data := bytes.Repeat([]byte{'a'}, n)
	if term != "" {
		for i := 80 - len(term); i < len(data); i += 80 {
			copy(data[i:], term)
		}
	}

	return data
}

// BenchmarkPR12NormalizeSparse measures normalization with 80-byte LF/CRLF
// lines or no terminators. Input is immutable and setup is outside the loop.
func BenchmarkPR12NormalizeSparse(b *testing.B) {
	for _, target := range []struct {
		name string
		term LineTerminator
	}{{"LF", LF}, {"CRLF", CRLF}} {
		b.Run("target="+target.name, func(b *testing.B) {
			for _, input := range []struct {
				name string
				term string
			}{{"lf80", "\n"}, {"crlf80", "\r\n"}, {"none", ""}} {
				b.Run("input="+input.name, func(b *testing.B) {
					benchkit.Run(b, []int{1024, 65536}, func(b *testing.B, n int) {
						data := benchmarkPR12SparseLineBreaks(n, input.term)
						b.ReportAllocs()
						for b.Loop() {
							NormalizeLineTerminators(data, target.term)
						}
					})
				})
			}
		})
	}
}
