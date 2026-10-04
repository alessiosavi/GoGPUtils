package textnorm

import (
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

const pr5Drop = "\u200b"

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
func BenchmarkPR5RemoveFormatChars(b *testing.B) {
	p := New().RemoveFormatChars()
	for _, name := range []string{"clean", "late-removal", "invalid"} {
		b.Run("case="+name, func(b *testing.B) {
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				text := pr5BenchInput(name, n)
				b.ReportAllocs()
				for b.Loop() {
					if _, err := p.Run(text); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}
