package cryptoutil

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/alessiosavi/GoGPUtils/internal/benchkit"
)

// Private, verbatim function bodies from immutable BASE 62dc6b5.
func baseDecrypt(ciphertextB64 string, key []byte) ([]byte, error) {
	if !isValidKeySize(len(key)) {
		return nil, ErrInvalidKeySize
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return nil, ErrInvalidCiphertext
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrInvalidCiphertext
	}

	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	return plaintext, nil
}

func baseHashString(s string) string {
	h := sha256.Sum256([]byte(s))

	return encodeHex(h[:])
}

func zzWrapCiphertext(s string, width int, separator string) string {
	var b strings.Builder
	for len(s) > width {
		b.WriteString(s[:width])
		b.WriteString(separator)
		s = s[width:]
	}
	b.WriteString(s)
	return b.String()
}

func zzCheckDecrypt(t *testing.T, ciphertext string, key []byte) ([]byte, error) {
	t.Helper()
	keyBefore := bytes.Clone(key)
	ciphertextBefore := strings.Clone(ciphertext)
	want, wantErr := baseDecrypt(ciphertext, key)
	got, gotErr := Decrypt(ciphertext, key)
	if gotErr != wantErr || !bytes.Equal(got, want) || (got == nil) != (want == nil) { //nolint:errorlint // Error identity is part of the BASE contract.
		t.Fatalf("Decrypt mismatch: len=%d/%d nil=%v/%v error=%v/%v", len(got), len(want), got == nil, want == nil, gotErr, wantErr)
	}
	if !bytes.Equal(key, keyBefore) || ciphertext != ciphertextBefore {
		t.Fatal("Decrypt mutated its inputs")
	}
	if cap(got) != len(got) {
		t.Fatalf("result capacity=%d, want length=%d", cap(got), len(got))
	}
	text, textErr := DecryptString(ciphertext, key)
	if textErr != wantErr || text != string(want) { //nolint:errorlint // Check exact error identity, including nil.
		t.Fatal("DecryptString differs from BASE plaintext/error")
	}
	return got, gotErr
}

func TestZZDecryptDifferential(t *testing.T) {
	for _, size := range []int{16, 24, 32} {
		key := benchkit.Bytes(size)
		for _, n := range []int{0, 1, 16, 1024, 65536} {
			plaintext := benchkit.Bytes(n)
			ciphertext, err := Encrypt(plaintext, key)
			if err != nil {
				t.Fatal(err)
			}
			for _, format := range []struct{ name, text string }{
				{"standard", ciphertext},
				{"mime76", zzWrapCiphertext(ciphertext, 76, "\n")},
				{"newline-every-4", zzWrapCiphertext(ciphertext, 4, "\n")},
				{"crlf", zzWrapCiphertext(ciphertext, 76, "\r\n")},
				{"trailing-newline", ciphertext + "\n"},
				{"large-slack", ciphertext + strings.Repeat("\r\n", 4096)},
			} {
				t.Run(fmt.Sprintf("key=%d/n=%d/%s", size, n, format.name), func(t *testing.T) {
					got, err := zzCheckDecrypt(t, format.text, key)
					if err != nil || !bytes.Equal(got, plaintext) || (n == 0 && got != nil) {
						t.Fatal("valid ciphertext did not reproduce plaintext/nilness")
					}
					other, err := Decrypt(format.text, key)
					if err != nil {
						t.Fatal(err)
					}
					if len(got) > 0 {
						got[0] ^= 0xff
					}
					got = append(got, 0xff)
					for i := range got {
						got[i] = 0
					}
					if !bytes.Equal(other, plaintext) {
						t.Fatal("mutating/appending a result changed another result")
					}
					again, _ := zzCheckDecrypt(t, format.text, key)
					if !bytes.Equal(again, plaintext) {
						t.Fatal("mutating a result changed a later decryption")
					}
				})
			}
		}
	}
}

func TestZZDecryptErrors(t *testing.T) {
	type decryptCase struct {
		name, text string
		key        []byte
		want       error
	}
	for _, size := range []int{16, 24, 32} {
		key := benchkit.Bytes(size)
		ciphertext, err := Encrypt(benchkit.Bytes(16), key)
		if err != nil {
			t.Fatal(err)
		}
		raw, err := base64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			t.Fatal(err)
		}
		cases := []decryptCase{
			{"invalid-base64", "!", key, ErrInvalidCiphertext},
			{"invalid-short", "AA!", key, ErrInvalidCiphertext},
			{"invalid-key-before-base64", "!", nil, ErrInvalidKeySize},
		}
		wrongKey := bytes.Clone(key)
		wrongKey[0] ^= 1
		cases = append(cases, decryptCase{"wrong-key", ciphertext, wrongKey, ErrDecryptFailed})
		for _, n := range []int{0, 1, 11, 12, 13, 27, 28, len(raw) - 1} {
			want := ErrDecryptFailed
			if n < 12 {
				want = ErrInvalidCiphertext
			}
			cases = append(cases, decryptCase{fmt.Sprintf("truncated-%d", n), base64.StdEncoding.EncodeToString(raw[:n]), key, want})
		}
		for _, pos := range []int{0, 12, len(raw) - 1} {
			tampered := bytes.Clone(raw)
			tampered[pos] ^= 1
			cases = append(cases, decryptCase{fmt.Sprintf("tampered-%d", pos), base64.StdEncoding.EncodeToString(tampered), key, ErrDecryptFailed})
		}
		for _, n := range []int{0, 1, 15, 17, 23, 25, 31, 33, 64} {
			cases = append(cases, decryptCase{fmt.Sprintf("invalid-key-%d", n), ciphertext, benchkit.Bytes(n), ErrInvalidKeySize})
		}
		for _, tc := range cases {
			t.Run(fmt.Sprintf("key=%d/%s", size, tc.name), func(t *testing.T) {
				for _, text := range []string{tc.text, tc.text + strings.Repeat("\n", 4096)} {
					got, err := zzCheckDecrypt(t, text, tc.key)
					if err != tc.want || got != nil { //nolint:errorlint // Require the documented sentinel itself.
						t.Fatalf("got %v, %v; want nil, %v", got, err, tc.want)
					}
				}
			})
		}
	}
}

func TestZZHashStringDifferential(t *testing.T) {
	lengths := make([]int, 1101)
	for n := range lengths {
		lengths[n] = n
	}
	lengths = append(lengths, 4095, 4096, 4097, 65536)
	for _, n := range lengths {
		for _, s := range []string{benchkit.Text(n, benchkit.Mixed), string(benchkit.Bytes(n))} {
			if got, want := HashString(s), baseHashString(s); got != want {
				t.Fatalf("n=%d: got %q, want %q", n, got, want)
			}
		}
	}
	for _, s := range []string{"\xff\xfe\x00", "�", "café\x80東京"} {
		if HashString(s) != baseHashString(s) {
			t.Fatalf("mismatch on %q", s)
		}
	}
}

func FuzzZZDecryptDifferential(f *testing.F) {
	f.Add("!", []byte(nil))
	f.Add("AA!", make([]byte, 16))
	for _, size := range []int{16, 24, 32} {
		key := benchkit.Bytes(size)
		for _, n := range []int{0, 1, 16, 1024} {
			ciphertext, err := Encrypt(benchkit.Bytes(n), key)
			if err != nil {
				f.Fatal(err)
			}
			f.Add(ciphertext, key)
			f.Add(zzWrapCiphertext(ciphertext, 4, "\r\n"), key)
		}
	}
	f.Fuzz(func(t *testing.T, ciphertext string, key []byte) {
		_, _ = zzCheckDecrypt(t, ciphertext, key)
	})
}

func FuzzZZHashStringDifferential(f *testing.F) {
	for _, s := range []string{"", "\xff\xfe\x00", "� café 東京", strings.Repeat("x", 32), strings.Repeat("y", 33), strings.Repeat("z", 513)} {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, s string) {
		if got, want := HashString(s), baseHashString(s); got != want {
			t.Fatalf("got %q, want %q", got, want)
		}
	})
}

// zzMeasureAllocations measures only allocations, never elapsed time. Callers
// prepare fixtures outside fn and keep results alive in a sink. No test using
// this helper may run in parallel because it temporarily changes GC settings.
func zzMeasureAllocations(fn func()) (allocs float64, bytesPerCall uint64) {
	previousProcs := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(previousProcs)
	previousGC := debug.SetGCPercent(-1)
	defer debug.SetGCPercent(previousGC)
	const runs = 100
	allocs = testing.AllocsPerRun(runs, fn)
	var before, after runtime.MemStats
	runtime.ReadMemStats(&before)
	for range runs {
		fn()
	}
	runtime.ReadMemStats(&after)
	return allocs, (after.TotalAlloc - before.TotalAlloc) / runs
}
