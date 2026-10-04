package textnorm

import (
	"fmt"
	"math/rand/v2"
	"testing"
)

func zzASCIIByteLoop(s string) bool {
	for i := range len(s) {
		if s[i] >= 0x80 {
			return false
		}
	}
	return true
}

func TestZZASCIIScanEquivalence(t *testing.T) {
	check := func(t *testing.T, s string) {
		t.Helper()
		if got, want := isASCII(s), zzASCIIByteLoop(s); got != want {
			t.Fatalf("isASCII(%x) = %v, want %v", s, got, want)
		}
	}
	rng := rand.New(rand.NewPCG(8, 80))
	for n := range 81 {
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			for offset := range 8 {
				// Exercise unaligned substrings with non-ASCII bytes just outside
				// the input, including empty strings and every word/tail boundary.
				buf := make([]byte, n+15)
				for i := range buf {
					buf[i] = 0xff
				}
				for i := range n {
					buf[offset+i] = byte((i*17 + n) % 128)
				}
				check(t, string(buf)[offset:offset+n])
				for pos := range n {
					ascii := buf[offset+pos]
					for _, high := range []byte{0x80, 0xc3, 0xff} {
						buf[offset+pos] = high
						check(t, string(buf)[offset:offset+n])
					}
					buf[offset+pos] = ascii
				}
			}
			buf := make([]byte, n)
			for range 100 {
				for i := range buf {
					buf[i] = byte(rng.Uint32())
				}
				check(t, string(buf))
				for i := range buf {
					buf[i] &= 0x7f
				}
				check(t, string(buf))
			}
		})
	}
}

func TestZZASCIIScanSingleBytes(t *testing.T) {
	for b := range 256 {
		s := string([]byte{byte(b)})
		if got, want := isASCII(s), b < 0x80; got != want {
			t.Fatalf("isASCII(%x) = %v, want %v", s, got, want)
		}
	}
}
