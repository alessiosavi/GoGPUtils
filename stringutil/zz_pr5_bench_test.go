package stringutil

import (
	"strings"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

const pr5Drop = "\x07"

// n is the exact input byte length; suffixes fit within that budget.
func pr5BenchInput(name string, n int) string {
	switch name {
	case "clean":
		return benchkit.Text(n, benchkit.Mixed)
	case "late-removal":
		return benchkit.Text(n-len(pr5Drop), benchkit.Mixed) + pr5Drop
	case "invalid":
		return benchkit.Text(n-1, benchkit.Mixed) + "\xff"
	}
	panic(name)
}
func BenchmarkPR5RemoveNonPrintable(b *testing.B) {
	for _, name := range []string{"clean", "late-removal", "invalid"} {
		b.Run("case="+name, func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				text := pr5BenchInput(name, n)
				b.ReportAllocs()
				for b.Loop() {
					RemoveNonPrintable(text)
				}
			})
		})
	}
}

func pr5LineInput(name string, n int) string {
	if name == "cr" {
		return benchWordText(n, "\r")
	}
	var b strings.Builder
	for i, w := range benchkit.Words(n) {
		if i > 0 {
			b.WriteString([]string{"\r\n", "\r", "\n"}[(i-1)%3])
		}
		b.WriteString(w)
		if b.Len() >= n {
			break
		}
	}
	return b.String()[:n]
}
func BenchmarkPR5Lines(b *testing.B) {
	for _, name := range []string{"cr", "mixed"} {
		b.Run("case="+name, func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				text := pr5LineInput(name, n)
				b.ReportAllocs()
				for b.Loop() {
					Lines(text)
				}
			})
		})
	}
}
func TestZZPR5MixedFixture(t *testing.T) {
	for _, n := range benchkit.Sizes {
		s := pr5LineInput("mixed", n)
		crlf := strings.Count(s, "\r\n")
		lone := strings.Count(s, "\r") - crlf
		if len(s) != n || crlf == 0 || lone == 0 {
			t.Fatalf("n=%d len=%d CRLF=%d lone-CR=%d", n, len(s), crlf, lone)
		}
		t.Logf("n=%d CRLF=%d lone-CR=%d", n, crlf, lone)
	}
}
