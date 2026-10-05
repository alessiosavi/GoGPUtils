package cryptoutil

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

var (
	zzDecryptBytes []byte
	zzDecryptText  string
	zzDecryptError error
)

func TestZZDecryptAllocations(t *testing.T) {
	for _, size := range []int{16, 24, 32} {
		key := benchkit.Bytes(size)
		for _, n := range []int{16, 1024, 65536} {
			ciphertext, err := Encrypt(benchkit.Bytes(n), key)
			if err != nil {
				t.Fatal(err)
			}
			for _, textResult := range []bool{false, true} {
				t.Run(fmt.Sprintf("key=%d/n=%d/string=%v", size, n, textResult), func(t *testing.T) {
					baseAllocs, baseBytes := zzMeasureAllocations(func() {
						p, e := baseDecrypt(ciphertext, key)
						zzDecryptError = e
						if textResult {
							zzDecryptText = string(p)
						} else {
							zzDecryptBytes = p
						}
					})
					gotAllocs, gotBytes := zzMeasureAllocations(func() {
						if textResult {
							zzDecryptText, zzDecryptError = DecryptString(ciphertext, key)
						} else {
							zzDecryptBytes, zzDecryptError = Decrypt(ciphertext, key)
						}
					})
					if zzDecryptError != nil {
						t.Fatal(zzDecryptError)
					}
					t.Logf("allocs %.0f -> %.0f; bytes %d -> %d", baseAllocs, gotAllocs, baseBytes, gotBytes)
					// Saving the plaintext allocation may include allocator rounding.
					if gotAllocs >= baseAllocs || gotBytes > baseBytes || baseBytes-gotBytes < uint64(n) {
						t.Fatalf("want fewer allocations and at least %d bytes saved", n)
					}
				})
			}
		}
	}
}

func TestZZDecryptFallback(t *testing.T) {
	key := benchkit.Bytes(32)
	for _, n := range []int{0, 2, 16, 1024, 65536} {
		ciphertext, err := Encrypt(benchkit.Bytes(n), key)
		if err != nil {
			t.Fatal(err)
		}
		for _, mode := range []struct {
			name, text string
			fallback   bool
		}{
			{"standard", ciphertext, false},
			{"mime76", zzWrapCiphertext(ciphertext, 76, "\n"), false},
			{"newline-every-4", zzWrapCiphertext(ciphertext, 4, "\n"), n >= 1024},
			{"large-slack", ciphertext + strings.Repeat("\n", 16384), true},
		} {
			t.Run(fmt.Sprintf("n=%d/%s", n, mode.name), func(t *testing.T) {
				zzAssertDecryptPath(t, mode.text, key, n, mode.fallback)
			})
		}
	}
	// n=2 makes the threshold exactly representable by DecodeString's capacity
	// (a multiple of three). Exercise equality and the next capacity above it.
	ciphertext, err := Encrypt(benchkit.Bytes(2), key)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name       string
		encodedLen int
		fallback   bool
	}{{"at-budget", 88, false}, {"over-budget", 92, true}} {
		t.Run(tc.name, func(t *testing.T) {
			text := ciphertext + strings.Repeat("\n", tc.encodedLen-len(ciphertext))
			zzAssertDecryptPath(t, text, key, 2, tc.fallback)
		})
	}
}

func zzAssertDecryptPath(t *testing.T, ciphertext string, key []byte, n int, fallback bool) {
	t.Helper()
	decoded, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		t.Fatal(err)
	}
	excess, budget := cap(decoded)-n, 64+n/8
	if (excess > budget) != fallback {
		t.Fatalf("fixture has excess=%d budget=%d, want fallback=%v", excess, budget, fallback)
	}
	got, err := zzCheckDecrypt(t, ciphertext, key)
	if err != nil || len(got) != n || (got == nil) != (n == 0) || cap(got) != n {
		t.Fatal("plaintext/nil/capacity contract failed")
	}
	baseAllocs, baseBytes := zzMeasureAllocations(func() { zzDecryptBytes, zzDecryptError = baseDecrypt(ciphertext, key) })
	gotAllocs, gotBytes := zzMeasureAllocations(func() { zzDecryptBytes, zzDecryptError = Decrypt(ciphertext, key) })
	if zzDecryptError != nil {
		t.Fatal(zzDecryptError)
	}
	if fallback || n == 0 {
		if gotAllocs > baseAllocs || gotBytes > baseBytes {
			t.Error("fallback/empty-input allocations and bytes must not exceed BASE")
		}
	} else if gotAllocs >= baseAllocs || gotBytes > baseBytes || baseBytes-gotBytes < uint64(n) {
		t.Errorf("want fewer allocations and at least %d bytes saved", n)
	}
	t.Logf("excess=%d budget=%d fallback=%v; allocs %.0f -> %.0f; bytes %d -> %d", excess, budget, fallback, baseAllocs, gotAllocs, baseBytes, gotBytes)
}
