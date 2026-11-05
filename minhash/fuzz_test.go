package minhash

import (
	"encoding/binary"
	"errors"
	"slices"
	"strings"
	"testing"
)

func FuzzShingles(f *testing.F) {
	f.Add([]byte("perché\x00città\x00\x00e\u0300"), 3)
	f.Add([]byte{}, 1)
	f.Add([]byte("a\x00b"), -1)
	f.Add([]byte(strings.Repeat("x", 128)), 2)
	f.Fuzz(func(t *testing.T, data []byte, k int) {
		// Bound quadratic oracle work while preserving arbitrary bytes and k.
		if len(data) > 2048 {
			return
		}
		var tokens []string
		if len(data) > 0 {
			tokens = strings.Split(string(data), "\x00")
		}
		got, err := Shingles(tokens, k)
		if k < 1 {
			wantError(t, err, ErrInvalidParameter)
			return
		}
		if err != nil || !slices.Equal(got, oracleShingles(tokens, k)) {
			t.Fatal("oracle mismatch", err)
		}
		for i := 1; i < len(got); i++ {
			if got[i-1] >= got[i] {
				t.Fatal("not sorted and unique")
			}
		}
	})
}

func fuzzValues(data []byte) []uint64 {
	var out []uint64
	for len(data) >= 8 {
		out = append(out, binary.LittleEndian.Uint64(data))
		data = data[8:]
	}
	if len(data) > 0 {
		var last [8]byte
		copy(last[:], data)
		out = append(out, binary.LittleEndian.Uint64(last[:]))
	}
	return out
}

func FuzzSign(f *testing.F) {
	f.Add([]byte("abcdefgh01234567"), 64, uint64(42))
	f.Add([]byte{}, 1, uint64(0))
	f.Add([]byte{255}, 0, uint64(1))
	f.Fuzz(func(t *testing.T, data []byte, n int, seed uint64) {
		if len(data) > 2048 {
			return
		}
		h, err := New(n, seed)
		if n < 1 || n > MaxNumPerm {
			wantError(t, err, ErrInvalidParameter)
			return
		}
		if err != nil {
			t.Fatal(err)
		}
		xs := fuzzValues(data)
		sig, err := h.Sign(xs)
		dst := make(Signature, n)
		intoErr := h.SignInto(dst, xs)
		if len(xs) == 0 {
			wantError(t, err, ErrEmptySet)
			wantError(t, intoErr, ErrEmptySet)
			return
		}
		again := mustSign(t, mustNew(t, n, seed), xs)
		if err != nil || intoErr != nil || !slices.Equal(sig, dst) || !slices.Equal(sig, again) {
			t.Fatal("sign mismatch", err, intoErr)
		}
	})
}

func FuzzSimilarity(f *testing.F) {
	f.Add([]byte("abcdefgh"), []byte("abcdefgi"))
	f.Add([]byte{}, []byte{})
	f.Add([]byte{1}, []byte("123456789"))
	f.Fuzz(func(t *testing.T, aa, bb []byte) {
		a, b := fuzzValues(aa), fuzzValues(bb)
		ab, err := Similarity(a, b)
		ba, reverseErr := Similarity(b, a)
		if len(a) != len(b) {
			wantError(t, err, ErrLengthMismatch)
			wantError(t, reverseErr, ErrLengthMismatch)
			return
		}
		if len(a) == 0 {
			if !errors.Is(err, ErrInvalidParameter) {
				t.Fatal(err)
			}
			return
		}
		self, selfErr := Similarity(a, a)
		if err != nil || reverseErr != nil || selfErr != nil || ab != ba || ab < 0 || ab > 1 || self != 1 {
			t.Fatal("similarity invariants", ab, ba, self, err)
		}
	})
}
