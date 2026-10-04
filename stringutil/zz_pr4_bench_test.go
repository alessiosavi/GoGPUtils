package stringutil

import (
	"strings"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

// BenchmarkPR4Invalid places an invalid byte inside the emitted span.
func BenchmarkPR4Invalid(b *testing.B) {
	for _, name := range []string{"Reverse", "Truncate", "TruncateWords", "TruncateRunes", "SafeSlice"} {
		b.Run("func="+name, func(b *testing.B) {
			benchkit.Run(b, []int{1024, 65536}, func(b *testing.B, n int) {
				text := benchkit.Text(n, benchkit.ASCII)
				limit := n / 2
				pos := limit - 1
				switch name {
				case "Reverse":
					pos = n - 1
				case "Truncate":
					pos = limit - 4
				case "TruncateWords":
					cut := limit - 3
					if space := strings.LastIndexByte(text[:cut], ' '); space > 0 {
						cut = space
					}
					pos = cut - 1
				case "SafeSlice":
					pos = 3*n/4 - 1
				}
				text = text[:pos] + "\xff" + text[pos+1:]
				b.ReportAllocs()
				for b.Loop() {
					switch name {
					case "Reverse":
						Reverse(text)
					case "Truncate":
						Truncate(text, limit, "...")
					case "TruncateWords":
						TruncateWords(text, limit, "...")
					case "TruncateRunes":
						TruncateRunes(text, limit)
					case "SafeSlice":
						SafeSlice(text, n/4, 3*n/4)
					}
				}
			})
		})
	}
}
