package cryptoutil

import (
	"strings"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

// BenchmarkPR13DecryptWrapped measures decoding and decryption of wrapped
// base64. n is the plaintext byte length; the key is 32 bytes. Encryption and
// newline insertion happen before the loop, with no trailing newline. Each
// iteration decrypts the same input without retaining its result.
func BenchmarkPR13DecryptWrapped(b *testing.B) {
	for _, mode := range []struct {
		name  string
		width int
	}{{"mime76", 76}, {"newline-every-4", 4}} {
		b.Run("case="+mode.name, func(b *testing.B) {
			key := benchkit.Bytes(32)
			benchkit.Run(b, []int{1024, 65536}, func(b *testing.B, n int) {
				ciphertext, err := Encrypt(benchkit.Bytes(n), key)
				if err != nil {
					b.Fatal(err)
				}
				var wrapped strings.Builder
				for len(ciphertext) > mode.width {
					wrapped.WriteString(ciphertext[:mode.width])
					wrapped.WriteByte('\n')
					ciphertext = ciphertext[mode.width:]
				}
				wrapped.WriteString(ciphertext)
				input := wrapped.String()
				b.ReportAllocs()
				for b.Loop() {
					if _, err := Decrypt(input, key); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

// BenchmarkPR13HashBoundary measures hashing and hex encoding at the short
// input cutoff and chunk boundaries. benchkit.Text trims at rune boundaries
// and pads to exactly n bytes. The deterministic mixed text is reused.
func BenchmarkPR13HashBoundary(b *testing.B) {
	benchkit.Run(b, []int{31, 32, 33, 511, 512, 513, 1023, 1025}, func(b *testing.B, n int) {
		text := benchkit.Text(n, benchkit.Mixed)
		b.ReportAllocs()
		for b.Loop() {
			HashString(text)
		}
	})
}
