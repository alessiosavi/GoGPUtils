package cryptoutil

import (
	"strconv"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

// BenchmarkEncrypt measures AES-GCM encryption and base64 encoding of n
// plaintext bytes with a key of 16, 24, or 32 bytes. The deterministic key
// and plaintext are reused; each operation generates a nonce with the real RNG.
func BenchmarkEncrypt(b *testing.B) {
	for _, size := range []int{16, 24, 32} {
		b.Run("key="+strconv.Itoa(size), func(b *testing.B) {
			key := benchkit.Bytes(size)
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				plaintext := benchkit.Bytes(n)
				b.ReportAllocs()
				for b.Loop() {
					if _, err := Encrypt(plaintext, key); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

// BenchmarkDecrypt measures base64 decoding and AES-GCM decryption with a key
// of 16, 24, or 32 bytes; n is the original plaintext length in bytes.
// Valid ciphertext is prepared from deterministic bytes before the loop and
// reused without mutation, so encryption setup is not measured.
func BenchmarkDecrypt(b *testing.B) {
	for _, size := range []int{16, 24, 32} {
		b.Run("key="+strconv.Itoa(size), func(b *testing.B) {
			key := benchkit.Bytes(size)
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				ciphertext, err := Encrypt(benchkit.Bytes(n), key)
				if err != nil {
					b.Fatal(err)
				}
				b.ReportAllocs()
				for b.Loop() {
					if _, err := Decrypt(ciphertext, key); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

// BenchmarkHash measures SHA-256 over n input bytes. The deterministic input
// is reused unchanged on every iteration.
func BenchmarkHash(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		data := benchkit.Bytes(n)
		b.ReportAllocs()
		for b.Loop() {
			Hash(data)
		}
	})
}

// BenchmarkGenerateKey measures generating a fresh key with the real RNG.
// key is the output length in bytes and covers the full accepted domain:
// 16, 24, and 32. No generated keys are retained between iterations.
func BenchmarkGenerateKey(b *testing.B) {
	for _, size := range []int{16, 24, 32} {
		b.Run("key="+strconv.Itoa(size), func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				if _, err := GenerateKey(size); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkEncryptString measures string conversion, AES-GCM encryption, and
// base64 encoding of n plaintext bytes with a key of 16, 24, or 32 bytes.
// Deterministic mixed UTF-8 text and the key are reused; nonce generation
// uses the real RNG on every iteration.
func BenchmarkEncryptString(b *testing.B) {
	for _, size := range []int{16, 24, 32} {
		b.Run("key="+strconv.Itoa(size), func(b *testing.B) {
			key := benchkit.Bytes(size)
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				plaintext := benchkit.Text(n, benchkit.Mixed)
				b.ReportAllocs()
				for b.Loop() {
					if _, err := EncryptString(plaintext, key); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

// BenchmarkDecryptString measures base64 decoding, AES-GCM decryption, and
// conversion to a string with a key of 16, 24, or 32 bytes; n is the original
// plaintext length in bytes. Valid ciphertext from deterministic mixed UTF-8
// text is prepared before the loop and reused; encryption setup is not measured.
func BenchmarkDecryptString(b *testing.B) {
	for _, size := range []int{16, 24, 32} {
		b.Run("key="+strconv.Itoa(size), func(b *testing.B) {
			key := benchkit.Bytes(size)
			benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
				ciphertext, err := EncryptString(benchkit.Text(n, benchkit.Mixed), key)
				if err != nil {
					b.Fatal(err)
				}
				b.ReportAllocs()
				for b.Loop() {
					if _, err := DecryptString(ciphertext, key); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}

// BenchmarkDeriveKey measures deriving a 32-byte key from the fixed password
// "password" and salt "salt". These inputs are reused unchanged each iteration.
func BenchmarkDeriveKey(b *testing.B) {
	const password, salt = "password", "salt"
	b.ReportAllocs()
	for b.Loop() {
		DeriveKey(password, salt)
	}
}

// BenchmarkGenerateNonce measures generating a fresh AES-GCM nonce with the
// real RNG; n is the standard nonce length of 12 bytes. No generated nonces
// are retained between iterations.
func BenchmarkGenerateNonce(b *testing.B) {
	b.Run("case=aes-gcm/n=12", func(b *testing.B) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := GenerateNonce(12); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkHashString measures SHA-256 and hex encoding over n input bytes
// of deterministic mixed UTF-8 text, reused unchanged each iteration.
func BenchmarkHashString(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		text := benchkit.Text(n, benchkit.Mixed)
		b.ReportAllocs()
		for b.Loop() {
			HashString(text)
		}
	})
}

// BenchmarkCompareHash measures hashing n input bytes and comparing against
// their matching SHA-256 digest. The deterministic data and expected digest
// are prepared before the loop and reused unchanged; every comparison must match.
func BenchmarkCompareHash(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		data := benchkit.Bytes(n)
		expected := Hash(data)
		b.ReportAllocs()
		for b.Loop() {
			if !CompareHash(data, expected) {
				b.Fatal("matching hash comparison failed")
			}
		}
	})
}

// BenchmarkRandomBytes measures generating n output bytes with the real RNG.
// Each iteration produces fresh bytes without retaining earlier results.
func BenchmarkRandomBytes(b *testing.B) {
	benchkit.Run(b, benchkit.Sizes, func(b *testing.B, n int) {
		b.ReportAllocs()
		for b.Loop() {
			if _, err := RandomBytes(n); err != nil {
				b.Fatal(err)
			}
		}
	})
}
