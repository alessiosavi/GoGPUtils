package randutil

import "testing"

// BenchmarkSecureStringWideCharset samples a 300-byte charset with the real
// crypto reader. The repeated-byte fixture is built outside the timed loop.
func BenchmarkSecureStringWideCharset(b *testing.B) {
	charset := make([]byte, 300)
	for i := range charset {
		charset[i] = byte(i)
	}
	chars := string(charset)
	b.Run("n=1024", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := SecureString(1024, chars); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkSecureIntWideBound samples a wide int bound with the real crypto reader.
func BenchmarkSecureIntWideBound(b *testing.B) {
	b.Run("max=2147483647", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := SecureInt(2147483647); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkSecureInt64WideBound samples an eight-byte candidate with heavy
// rejection using the real crypto reader.
func BenchmarkSecureInt64WideBound(b *testing.B) {
	b.Run("max=4611686018427387905", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := SecureInt64(4611686018427387905); err != nil {
				b.Fatal(err)
			}
		}
	})
}
