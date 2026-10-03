package randutil

import (
	"bytes"
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"slices"
	"testing"
	"time"
)

// Tests replacing Reader must remain sequential and restore it before returning.
func useSecureReader(t *testing.T, reader io.Reader) {
	t.Helper()

	original := cryptorand.Reader
	cryptorand.Reader = reader
	t.Cleanup(func() {
		cryptorand.Reader = original
	})
}

func TestSecureString_UsesEveryCharsetPosition(t *testing.T) {
	for _, size := range []int{256, 300} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			defer func() {
				if recovered := recover(); recovered != nil {
					t.Errorf("SecureString panicked with charset length %d: %v", size, recovered)
				}
			}()

			reader := new(bytes.Reader)
			useSecureReader(t, reader)
			for position := range size {
				charset := bytes.Repeat([]byte{'.'}, size)
				charset[position] = '|'
				candidate := []byte{byte(position >> 8), byte(position)}
				if size == 256 {
					candidate = candidate[1:]
				}

				reader.Reset(candidate)
				got, err := SecureString(1, string(charset))
				if err != nil || got != "|" {
					t.Fatalf("charset position %d: got %q, %v; want |, nil", position, got, err)
				}
			}
		})
	}
}

type cyclingByteReader struct {
	next byte
}

func (r *cyclingByteReader) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = r.next
		r.next++
	}

	return len(p), nil
}

func TestSecureString_SelectsBytesWithoutModuloBias(t *testing.T) {
	useSecureReader(t, new(cyclingByteReader))

	charset := make([]byte, 88)
	for i := range charset {
		charset[i] = byte(i)
	}

	const perBucket = 256
	got, err := SecureString(len(charset)*perBucket, string(charset))
	if err != nil {
		t.Fatal(err)
	}

	var counts [88]int
	for i := range got {
		counts[got[i]]++
	}

	for i, count := range counts {
		if count != perBucket {
			t.Errorf("byte %d selected %d times, want %d", i, count, perBucket)
		}
	}
}

func TestSecureInt_RejectsBiasedCandidates(t *testing.T) {
	reader := new(bytes.Reader)
	useSecureReader(t, reader)

	// An all-ones candidate must be rejected; the following zero is valid.
	var candidates [16]byte
	binary.BigEndian.PutUint64(candidates[:8], math.MaxUint64)
	reader.Reset(candidates[:])
	got, err := SecureInt(math.MaxInt/2 + 2)
	if err != nil || got != 0 {
		t.Errorf("SecureInt() = %d, %v; want 0, nil after rejection", got, err)
	}
}

func TestSecureInt64_RejectsBiasedCandidates(t *testing.T) {
	var candidates [16]byte
	binary.BigEndian.PutUint64(candidates[:8], math.MaxUint64)
	useSecureReader(t, bytes.NewReader(candidates[:]))

	got, err := SecureInt64(math.MaxInt64/2 + 2)
	if err != nil || got != 0 {
		t.Errorf("SecureInt64() = %d, %v; want 0, nil after rejection", got, err)
	}
}

func TestRangeStep_StopsAtIntegerOverflow(t *testing.T) {
	tests := []struct {
		name             string
		start, end, step int
		want             []int
	}{
		{"ascending", math.MaxInt - 5, math.MaxInt, 4, []int{math.MaxInt - 5, math.MaxInt - 1}},
		{"descending", math.MinInt + 5, math.MinInt, -4, []int{math.MinInt + 5, math.MinInt + 1}},
		{"minimum step", 0, math.MinInt, math.MinInt, []int{0}},
		{"cross zero", math.MinInt, math.MaxInt, math.MaxInt, []int{math.MinInt, -1, math.MaxInt - 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan []int, 1)
			go func() {
				done <- RangeStep(tt.start, tt.end, tt.step)
			}()

			select {
			case got := <-done:
				if !slices.Equal(got, tt.want) {
					t.Errorf("RangeStep() = %v, want %v", got, tt.want)
				}
			case <-time.After(time.Second):
				t.Fatal("RangeStep did not terminate within one second")
			}
		})
	}
}
