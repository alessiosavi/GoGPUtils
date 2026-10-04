package randutil

import (
	"bytes"
	cryptorand "crypto/rand"
	"fmt"
	"io"
	"math"
	"math/big"
	"math/bits"
	"math/rand/v2"
	"slices"
	"testing"
	"testing/iotest"
)

// Private BASE 4c05787 bodies; only the function identifiers are renamed.
func baseSecureString(length int, charset string) (string, error) {
	if length <= 0 {
		return "", ErrInvalidLength
	}

	if len(charset) == 0 {
		return "", ErrEmptyCharset
	}

	result := make([]byte, length)
	charsetLen := big.NewInt(int64(len(charset)))

	for i := range result {
		// crypto/rand.Int uses rejection sampling to avoid modulo bias.
		index, err := cryptorand.Int(cryptorand.Reader, charsetLen)
		if err != nil {
			return "", err
		}

		result[i] = charset[index.Int64()]
	}

	return string(result), nil
}

//nolint:revive // Preserve the BASE function body and parameter name verbatim.
func baseSecureInt(max int) (int, error) {
	if max <= 0 {
		return 0, ErrInvalidLength
	}

	n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}

	return int(n.Int64()), nil
}

//nolint:revive // Preserve the BASE function body and parameter name verbatim.
func baseSecureInt64(max int64) (int64, error) {
	if max <= 0 {
		return 0, ErrInvalidLength
	}

	n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(max))
	if err != nil {
		return 0, err
	}

	return n.Int64(), nil
}

func baseRangeStep(start, end, step int) []int {
	if step == 0 || (step > 0 && start >= end) || (step < 0 && start <= end) {
		return nil
	}

	var result []int

	for i := start; (step > 0 && i < end) || (step < 0 && i > end); {
		result = append(result, i)

		next := i + step
		if (step > 0 && next <= i) || (step < 0 && next >= i) {
			break
		}

		i = next
	}

	return result
}

// A recorder outside the short-read adapters measures bytes delivered to the
// sampler, including DataErrReader's buffered bytes, and every request size.
type oracleRecorder struct {
	reader   io.Reader
	position int
	requests []int
}

func (r *oracleRecorder) Read(p []byte) (int, error) {
	r.requests = append(r.requests, len(p))
	n, e := r.reader.Read(p)
	r.position += n
	return n, e
}

var errOracleSentinel = fmt.Errorf("persistent sentinel")
var errOracleWrappedEOF = fmt.Errorf("wrap: %w", io.EOF)
var errOracleWrappedUnexpected = fmt.Errorf("wrap: %w", io.ErrUnexpectedEOF)

// terminalReader attaches an error to the last data, then persists with it
// or returns EOF forever. Neither policy delivers any more data.
// Unlike "error on the last call", its endpoint is independent of request size.
type oracleTerminalReader struct {
	data     []byte
	position int
	err      error
	eofAfter bool
}

func (r *oracleTerminalReader) Read(p []byte) (int, error) {
	if r.eofAfter && r.position == len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.position:])
	r.position += n
	if r.position == len(r.data) {
		return n, r.err
	}
	return n, nil
}
func oracleReader(kind int, data []byte) *oracleRecorder {
	var r io.Reader = bytes.NewReader(bytes.Clone(data))
	switch kind {
	case 1:
		r = iotest.OneByteReader(r)
	case 2:
		r = iotest.HalfReader(r)
	case 3:
		r = iotest.DataErrReader(r)
	case 4:
		r = iotest.ErrReader(io.ErrUnexpectedEOF)
	case 5:
		r = iotest.ErrReader(errOracleWrappedEOF)
	case 6:
		r = iotest.ErrReader(errOracleSentinel)
	case 7:
		r = &oracleTerminalReader{data: bytes.Clone(data), err: errOracleSentinel}
	}
	return &oracleRecorder{reader: r}
}
func oracleStringCase(t *testing.T, length int, charset string, data []byte, kind int) {
	t.Helper()
	baseReader, headReader := oracleReader(kind, data), oracleReader(kind, data)
	oracleStringReaders(t, length, charset, baseReader, headReader)
}

func oracleStringReaders(t *testing.T, length int, charset string, baseReader, headReader *oracleRecorder) {
	t.Helper()
	useSecureReader(t, baseReader)
	want, we := baseSecureString(length, charset)
	useSecureReader(t, headReader)
	got, ge := SecureString(length, charset)
	//nolint:errorlint // Compare exact error identity, including wrapped errors.
	if got != want || ge != we || headReader.position != baseReader.position {
		t.Fatalf("len=%d charset=%d: got len=%d err=%v pos=%d; BASE len=%d err=%v pos=%d",
			length, len(charset), len(got), ge, headReader.position, len(want), we, baseReader.position)
	}
	// Further reads observe the same logical stream, not merely its byte count.
	a, b := make([]byte, 19), make([]byte, 19)
	an, ae := headReader.Read(a)
	bn, be := baseReader.Read(b)
	//nolint:errorlint // Compare exact error identity, including wrapped errors.
	if an != bn || ae != be || !bytes.Equal(a[:an], b[:bn]) {
		t.Fatal("remaining reader stream differs")
	}
}
func TestOracleSecureStringStreams(t *testing.T) {
	rng := rand.New(rand.NewPCG(913, 17))
	lengths := []int{1, 2, 3, 61, 62, 63, 64, 65, 255, 256, 257, 300, 4095, 4097, 65535, 65536, 65537, 135536}
	// Reuse charsets outside cases; duplicates are unavoidable above 256.
	charsets := make([]string, len(lengths)+12)
	for i := range charsets {
		n := rng.IntN(150000) + 1
		if i < len(lengths) {
			n = lengths[i]
		}
		b := make([]byte, n)
		for j := range b {
			b[j] = byte(rng.Uint32())
		}
		charsets[i] = string(b)
	}
	for i := range 20000 {
		cs := charsets[i%len(charsets)]
		n := rng.IntN(300) + 1
		if i%131 == 0 {
			n = 5000
		}
		k := (bits.Len64(uint64(len(cs)-1)) + 7) / 8
		size := rng.IntN(2*n*max(k, 1) + 2)
		data := make([]byte, size)
		for j := range data {
			data[j] = byte(rng.Uint32())
		}
		t.Run(fmt.Sprint(i), func(t *testing.T) { oracleStringCase(t, n, cs, data, (i/len(charsets))%8) })
	}
	// Named deterministic cases make rare error/boundary paths compulsory.
	for _, n := range []int{-1, 0, 1, 2, 4095, 4096, 4097, 5000} {
		for _, cs := range []string{"", "x", "ab", "abc", string(bytes.Repeat([]byte("abcd"), 75))} {
			for _, kind := range []int{0, 1, 2, 3, 4, 5, 6, 7} {
				t.Run(fmt.Sprintf("edge/%d/%d/%d", n, len(cs), kind), func(t *testing.T) {
					oracleStringCase(t, n, cs, bytes.Repeat([]byte{255, 0, 1, 2}, 6000), kind)
				})
			}
		}
	}
	for _, data := range [][]byte{{0}, {0, 0}, {255}, {255, 0}, {255, 255, 0, 0}, {0, 0, 0}} {
		for _, cs := range []string{"ab", "abc", string(bytes.Repeat([]byte("abcd"), 75))} {
			for _, kind := range []int{0, 3, 7} {
				for _, n := range []int{1, 2, 3} {
					t.Run(fmt.Sprintf("terminal/%x/%d/%d/%d", data, len(cs), kind, n), func(t *testing.T) { oracleStringCase(t, n, cs, data, kind) })
				}
			}
		}
	}
}

// Error objects are shared by both fresh readers; equality checks identity.
func oracleErrors() []error {
	return []error{io.EOF, io.ErrUnexpectedEOF, errOracleSentinel, errOracleWrappedEOF, errOracleWrappedUnexpected}
}

func oracleTerminalPair(data []byte, err error, eofAfter bool, adapter int) (*oracleRecorder, *oracleRecorder) {
	newReader := func() *oracleRecorder {
		var r io.Reader = &oracleTerminalReader{data: bytes.Clone(data), err: err, eofAfter: eofAfter}
		switch adapter {
		case 1:
			r = iotest.OneByteReader(r)
		case 2:
			r = iotest.HalfReader(r)
		}
		return &oracleRecorder{reader: r}
	}
	return newReader(), newReader()
}

func TestOracleSecureStringErrorMatrix(t *testing.T) {
	for _, charsetLen := range []int{3, 300, 65537} {
		charset := make([]byte, charsetLen)
		for i := range charset {
			charset[i] = byte(i)
		}
		k := (bits.Len64(uint64(charsetLen-1)) + 7) / 8
		// With three requested outputs, 3*k is a complete first batch;
		// 1..3*k-1 cover every incomplete boundary and mid-candidate ending.
		// 4*k also exercises another fill after full-batch rejection.
		for size := 0; size <= 4*k; size++ {
			for _, pattern := range []string{"accept", "reject", "reject-then-accept"} {
				data := make([]byte, size)
				for i := range data {
					if pattern == "reject" || (pattern == "reject-then-accept" && i < k) {
						data[i] = 255
					}
				}
				for _, err := range oracleErrors() {
					for _, eofAfter := range []bool{false, true} {
						for adapter := range 3 {
							t.Run(fmt.Sprintf("k=%d/bytes=%d/%s/%v/eofAfter=%v/adapter=%d", k, size, pattern, err, eofAfter, adapter), func(t *testing.T) {
								a, b := oracleTerminalPair(data, err, eofAfter, adapter)
								oracleStringReaders(t, 3, string(charset), a, b)
							})
						}
					}
				}
			}
		}
	}
}

func TestOracleSecureIntErrorMatrix(t *testing.T) {
	// Every candidate width, with accepted, rejected and truncated candidates.
	for k := 1; k <= 8; k++ {
		bound := int64(1)<<uint(8*(k-1)+1) + 1
		for size := 0; size <= 3*k; size++ {
			for _, pattern := range []string{"accept", "reject", "reject-then-accept"} {
				data := make([]byte, size)
				for i := range data {
					if pattern == "reject" || (pattern == "reject-then-accept" && i < k) {
						data[i] = 255
					}
				}
				for _, err := range oracleErrors() {
					for _, eofAfter := range []bool{false, true} {
						for adapter := range 3 {
							for _, useInt := range []bool{false, true} {
								if useInt && bound > int64(math.MaxInt) {
									continue
								}
								t.Run(fmt.Sprintf("k=%d/bytes=%d/%s/%v/eofAfter=%v/adapter=%d/int=%v", k, size, pattern, err, eofAfter, adapter, useInt), func(t *testing.T) {
									a, b := oracleTerminalPair(data, err, eofAfter, adapter)
									oracleIntReaders(t, bound, a, b, useInt)
								})
							}
						}
					}
				}
			}
		}
	}
}

// The replacement happens synchronously on the first Read, without a race.
type oracleSwitchReader struct {
	reader      io.Reader
	replacement io.Reader
	switched    bool
}

func (r *oracleSwitchReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)
	if !r.switched {
		cryptorand.Reader = r.replacement
		r.switched = true
	}
	return n, err
}

func TestOracleSecureIntCapturesReader(t *testing.T) {
	for _, bound := range []int64{3, 257, 65537, 1<<62 + 1} {
		k := (bits.Len64(uint64(bound-1)) + 7) / 8
		for _, short := range []bool{false, true} {
			for _, useInt := range []bool{false, true} {
				if useInt && bound > int64(math.MaxInt) {
					continue
				}
				t.Run(fmt.Sprintf("bound=%d/short=%v/int=%v", bound, short, useInt), func(t *testing.T) {
					data := append(bytes.Repeat([]byte{255}, k), make([]byte, k)...)
					newReader := func() (*oracleRecorder, *oracleRecorder) {
						replacementData := make([]byte, k)
						replacementData[k-1] = 1
						replacement := &oracleRecorder{reader: bytes.NewReader(replacementData)}
						var source io.Reader = bytes.NewReader(bytes.Clone(data))
						if short {
							source = iotest.OneByteReader(source)
						}
						return &oracleRecorder{reader: &oracleSwitchReader{reader: source, replacement: replacement}}, replacement
					}
					a, ar := newReader()
					b, br := newReader()
					oracleIntReaders(t, bound, a, b, useInt)
					if a.position != 2*k || len(ar.requests) != 0 || len(br.requests) != 0 {
						t.Fatalf("original position=%d, replacement requests BASE=%v candidate=%v", a.position, ar.requests, br.requests)
					}
				})
			}
		}
	}
}

type oracleRecoveringReader struct{ position int }

func (r *oracleRecoveringReader) Read(p []byte) (int, error) {
	if r.position == 0 {
		p[0] = 0
		r.position++
		return 1, errOracleSentinel
	}
	for i := range p {
		p[i] = 1
	}
	r.position += len(p)
	return len(p), nil
}
func TestSecureStringExcludedRecoveringReader(t *testing.T) {
	useSecureReader(t, new(oracleRecoveringReader))
	want, we := baseSecureString(4, "ab")
	useSecureReader(t, new(oracleRecoveringReader))
	got, ge := SecureString(4, "ab")
	// This reader is outside the guarantee. Validate each outcome independently;
	// do not require either implementation to match the other.
	for _, result := range []struct {
		value string
		err   error
	}{{want, we}, {got, ge}} {
		if result.err == nil && len(result.value) != 4 {
			t.Fatalf("successful result has length %d, want 4", len(result.value))
		}
		if result.err != nil && result.value != "" {
			t.Fatalf("error result must be empty, got %q", result.value)
		}
	}
	t.Logf("excluded recovering reader: BASE %q, %v; candidate %q, %v", want, we, got, ge)
}
func oracleIntCase(t *testing.T, bound int64, data []byte, kind int, useInt bool) {
	t.Helper()
	a, b := oracleReader(kind, data), oracleReader(kind, data)
	oracleIntReaders(t, bound, a, b, useInt)
}

func oracleIntReaders(t *testing.T, bound int64, a, b *oracleRecorder, useInt bool) {
	t.Helper()
	useSecureReader(t, a)
	var want, got int64
	var we, ge error
	if useInt {
		w, e := baseSecureInt(int(bound))
		want, we = int64(w), e
	} else {
		want, we = baseSecureInt64(bound)
	}
	useSecureReader(t, b)
	if useInt {
		g, e := SecureInt(int(bound))
		got, ge = int64(g), e
	} else {
		got, ge = SecureInt64(bound)
	}
	//nolint:errorlint // Compare exact error identity, including wrapped errors.
	if got != want || ge != we || a.position != b.position || !slices.Equal(a.requests, b.requests) {
		t.Fatalf("bound=%d int=%v: got (%d,%v) pos=%d reads=%v; BASE (%d,%v) pos=%d reads=%v", bound, useInt, got, ge, b.position, b.requests, want, we, a.position, a.requests)
	}
}
func TestOracleSecureIntStreams(t *testing.T) {
	bounds := []int64{-1, 0, 1, 2, 3, 255, 256, 257, 65535, 65536, 65537, 1<<31 - 1, 1 << 31, 1 << 32, 1<<53 + 1, math.MaxInt64/2 + 2, math.MaxInt64}
	// Exercise every candidate width, especially the otherwise-missing k=5,6,7.
	for shift := 1; shift < 63; shift++ {
		bounds = append(bounds, 1<<shift, 1<<shift+1)
	}
	rng := rand.New(rand.NewPCG(9, 14))
	for range 1000 {
		bounds = append(bounds, rng.Int64N(math.MaxInt64)+1)
	}
	for i, bound := range bounds {
		k := max(1, (bits.Len64(uint64(max(bound, 1)-1))+7)/8)
		streams := [][]byte{nil, bytes.Repeat([]byte{0}, 4*k), bytes.Repeat([]byte{255}, 3*k)}
		rejection := append(bytes.Repeat([]byte{255}, 3*k), bytes.Repeat([]byte{0}, 2*k)...)
		streams = append(streams, rejection)
		for size := 1; size <= k; size++ {
			streams = append(streams, make([]byte, size))
		}
		data := make([]byte, 4*k)
		for j := range data {
			data[j] = byte(rng.Uint32())
		}
		streams = append(streams, data)
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			for _, data := range streams {
				for kind := range 8 {
					oracleIntCase(t, bound, data, kind, false)
					oracleIntCase(t, min(bound, int64(math.MaxInt)), data, kind, true)
				}
			}
		})
	}
}
func TestSecureChoiceEdges(t *testing.T) {
	r := oracleReader(6, nil)
	useSecureReader(t, r)
	g, e := SecureChoice[int](nil)
	//nolint:errorlint // The public sentinel must be returned unchanged.
	if g != 0 || e != ErrEmptySlice {
		t.Fatalf("empty (%d,%v)", g, e)
	}
	g, e = SecureChoice([]int{37})
	if g != 37 || e != nil || len(r.requests) != 0 {
		t.Fatalf("singleton (%d,%v), reads=%v", g, e, r.requests)
	}
}
func TestOracleRangeStep(t *testing.T) {
	edges := []int{math.MinInt, math.MinInt + 1, math.MinInt + 5, -10, -1, 0, 1, 10, math.MaxInt - 5, math.MaxInt - 1, math.MaxInt}
	steps := []int{math.MinInt, -math.MaxInt, -4, -1, 0, 1, 4, math.MaxInt}
	check := func(start, end, step int) {
		t.Helper()
		a, b := baseRangeStep(start, end, step), RangeStep(start, end, step)
		if (a == nil) != (b == nil) || !slices.Equal(a, b) {
			t.Fatalf("(%d,%d,%d): got %v; BASE %v", start, end, step, b, a)
		}
	}
	for _, start := range edges {
		for _, end := range edges {
			for _, step := range steps {
				if step != 0 && ((step > 0 && start < end) || (step < 0 && start > end)) {
					distance, stride := uint64(end)-uint64(start), uint64(step)
					if step < 0 {
						distance = uint64(start) - uint64(end)
						stride = -stride
					}
					if (distance-1)/stride+1 > 10000 {
						continue
					} // Do not try to allocate exabytes.
				}
				check(start, end, step)
			}
		}
	}
	rng := rand.New(rand.NewPCG(234, 41))
	for range 20000 {
		check(rng.IntN(1001)-500, rng.IntN(1001)-500, rng.IntN(101)-50)
	}
}
