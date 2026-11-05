package minhash

import (
	"errors"
	"fmt"
	"math"
	"slices"
	"sync"
	"testing"
)

func mustNew(t testing.TB, n int, seed uint64) *Hasher {
	t.Helper()
	h, err := New(n, seed)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func mustSign(t testing.TB, h *Hasher, xs []uint64) Signature {
	t.Helper()
	sig, err := h.Sign(xs)
	if err != nil {
		t.Fatal(err)
	}
	return sig
}

func wantError(t *testing.T, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("error = %v; want %v", err, target)
	}
}

func TestShinglesEdges(t *testing.T) {
	for _, k := range []int{0, -1, math.MinInt} {
		_, err := Shingles(nil, k)
		wantError(t, err, ErrInvalidParameter)
	}
	empty, err := Shingles(nil, 5)
	if err != nil || empty != nil {
		t.Fatalf("empty = %v, %v", empty, err)
	}
	for _, tc := range []struct {
		tokens   []string
		k, count int
	}{
		{[]string{"a", "b"}, 5, 1}, {[]string{"a", "a", "a"}, 1, 1}, {[]string{"", ""}, 1, 1}, {[]string{""}, math.MaxInt, 1},
	} {
		got, err := Shingles(tc.tokens, tc.k)
		if err != nil || len(got) != tc.count {
			t.Fatalf("%v/%d: %v %v", tc.tokens, tc.k, got, err)
		}
	}
	for _, pair := range [][2][]string{
		{{"ab", "c"}, {"a", "bc"}}, {{"è"}, {"e\u0300"}}, {{"a", ""}, {"a"}}, {{""}, {"", ""}},
	} {
		a, _ := Shingles(pair[0], 8)
		b, _ := Shingles(pair[1], 8)
		if slices.Equal(a, b) {
			t.Fatalf("ambiguous tokens: %q", pair)
		}
	}
}

func TestValidation(t *testing.T) {
	for _, n := range []int{0, -1, 1025, math.MaxInt} {
		_, err := New(n, 0)
		wantError(t, err, ErrInvalidParameter)
	}
	h := mustNew(t, 64, 42)
	if h.NumPerm() != 64 {
		t.Fatal("NumPerm")
	}
	sig, err := h.Sign(nil)
	wantError(t, err, ErrEmptySet)
	if sig != nil {
		t.Fatal("empty signature must be nil on error")
	}
	dst := make(Signature, 64)
	for i := range dst {
		dst[i] = 77
	}
	wantError(t, h.SignInto(dst, nil), ErrEmptySet)
	for _, x := range dst {
		if x != 77 {
			t.Fatal("dst changed on error")
		}
	}
	wantError(t, h.SignInto(dst[:63], []uint64{1}), ErrLengthMismatch)
	_, err = Similarity(nil, Signature{1})
	wantError(t, err, ErrLengthMismatch)
	_, err = Similarity(nil, nil)
	wantError(t, err, ErrInvalidParameter)
	for _, bands := range []int{0, -1, 65, 3} {
		_, err = BandKeys(dst, bands)
		wantError(t, err, ErrInvalidParameter)
	}
	_, err = BandKeys(nil, 1)
	wantError(t, err, ErrInvalidParameter)
	for _, pair := range [][2]int{{0, 1}, {1, 0}, {-1, 1}, {1, -1}, {1024, 2}, {math.MaxInt, math.MaxInt}} {
		_, err = NewIndex(pair[0], pair[1])
		wantError(t, err, ErrInvalidParameter)
	}
	for _, threshold := range []float64{0, 1, -0.1, math.NaN(), math.Inf(1), math.Inf(-1)} {
		_, _, err = Bands(64, threshold)
		wantError(t, err, ErrInvalidParameter)
	}
	for _, n := range []int{0, -1, 1025} {
		_, _, err = Bands(n, 0.5)
		wantError(t, err, ErrInvalidParameter)
	}
	var zero Hasher
	_, err = zero.Sign([]uint64{1})
	wantError(t, err, ErrInvalidParameter)
	var idx Index
	wantError(t, idx.Add("", nil), ErrInvalidParameter)
	_, err = idx.Candidates(nil)
	wantError(t, err, ErrInvalidParameter)
}

func TestSignProperties(t *testing.T) {
	xs := []uint64{math.MaxUint64, prime, 0, 1, prime - 1, 7, 7}
	for _, seed := range []uint64{0, 42, math.MaxUint64} {
		large := mustNew(t, MaxNumPerm, seed)
		full := mustSign(t, large, xs)
		for _, n := range []int{1, 8, 64, 128, 1024} {
			h := mustNew(t, n, seed)
			if !slices.Equal(h.permutations, large.permutations[:n]) {
				t.Fatal("permutation prefix")
			}
			sig := mustSign(t, h, xs)
			if !slices.Equal(sig, full[:n]) {
				t.Fatal("signature prefix")
			}
			for _, x := range sig {
				if x >= prime {
					t.Fatal("signature range")
				}
			}
		}
	}
	h := mustNew(t, len(xs), 42)
	want := mustSign(t, h, xs)
	aliased := slices.Clone(xs)
	if err := h.SignInto(aliased, aliased); err != nil || !slices.Equal(aliased, want) {
		t.Fatal("alias support", err)
	}
	slices.Sort(xs)
	xs = slices.Compact(xs)
	if !slices.Equal(mustSign(t, h, xs), want) {
		t.Fatal("set order or duplicates changed signature")
	}
}

func TestSignIntoAllocs(t *testing.T) {
	h := mustNew(t, 128, 42)
	xs := make([]uint64, 300)
	for i := range xs {
		xs[i] = uint64(i) * 7919
	}
	dst := make(Signature, h.NumPerm())
	var err error
	allocs := testing.AllocsPerRun(100, func() { err = h.SignInto(dst, xs) })
	if err != nil || allocs != 0 {
		t.Fatalf("SignInto: %g allocations, %v", allocs, err)
	}
}

func TestJaccardAndSimilarity(t *testing.T) {
	for _, tc := range []struct {
		a, b []uint64
		want float64
	}{
		{nil, nil, 0}, {nil, []uint64{1}, 0}, {[]uint64{1, 2}, []uint64{1, 2}, 1}, {[]uint64{1, 2}, []uint64{2, 3}, 1.0 / 3},
	} {
		if got := Jaccard(tc.a, tc.b); got != tc.want {
			t.Fatalf("Jaccard = %g want %g", got, tc.want)
		}
	}
	a, b := Signature{1, 2, 3, 4}, Signature{1, 0, 3, 0}
	got, err := Similarity(a, b)
	if got != 0.5 || err != nil {
		t.Fatal(got, err)
	}
	for _, pair := range [][2][]uint64{{{1, 1, 1}, {1}}, {{3, 2, 1}, {1, 2, 3}}, {{0, 0}, {0, 0, 0}}} {
		if x := Jaccard(pair[0], pair[1]); x < 0 || x > 1 {
			t.Fatal("unspecified input escaped bounds")
		}
	}
}

func TestIndex(t *testing.T) {
	idx, err := NewIndex(2, 2)
	if err != nil {
		t.Fatal(err)
	}
	query := Signature{1, 2, 3, 4}
	// The first inserted match is in the second band: bucket traversal must
	// not determine output order. Empty IDs are legal.
	for _, entry := range []struct {
		id  string
		sig Signature
	}{
		{"z", Signature{9, 9, 3, 4}}, {"", query}, {"a", Signature{1, 2, 8, 8}}, {"none", Signature{7, 7, 7, 7}},
	} {
		if err := idx.Add(entry.id, entry.sig); err != nil {
			t.Fatal(err)
		}
	}
	wantError(t, idx.Add("z", query), ErrDuplicateID)
	wantError(t, idx.Add("bad", query[:3]), ErrLengthMismatch)
	_, err = idx.Candidates(query[:3])
	wantError(t, err, ErrLengthMismatch)
	got, err := idx.Candidates(query)
	if err != nil || !slices.Equal(got, []string{"z", "", "a"}) {
		t.Fatalf("candidates = %q %v", got, err)
	}
	// Index retains keys, not the caller's signature or result storage.
	query[0] = 99
	got[0] = "mutated"
	got, err = idx.Candidates(Signature{1, 2, 3, 4})
	if err != nil || !slices.Equal(got, []string{"z", "", "a"}) {
		t.Fatal("caller mutation changed index")
	}
	got, err = idx.Candidates(Signature{42, 42, 42, 42})
	if got != nil || err != nil {
		t.Fatal("no-match result", got, err)
	}
}

func TestBandsKnown(t *testing.T) {
	for _, tc := range []struct {
		threshold float64
		b, r      int
	}{{0.5, 16, 4}, {0.7, 8, 8}, {0.8, 4, 16}, {0.85, 4, 16}, {0.9, 4, 16}} {
		b, r, err := Bands(64, tc.threshold)
		if err != nil || b != tc.b || r != tc.r {
			t.Fatalf("Bands(64,%g)=%dx%d %v", tc.threshold, b, r, err)
		}
		t.Logf("Bands(64,%g)=%dx%d", tc.threshold, b, r)
	}
}

func TestConcurrentHasherAndIndex(t *testing.T) {
	h := mustNew(t, 64, 42)
	xs := []uint64{1, 4, 8, 9, 13, 21}
	want := mustSign(t, h, xs)
	idx, err := NewIndex(16, 4)
	if err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for worker := range 24 {
		wg.Go(func() {
			for j := range 8 {
				sig, err := h.Sign(xs)
				if err != nil || !slices.Equal(sig, want) {
					t.Error("concurrent signature", err)
					return
				}
				id := fmt.Sprintf("%d/%d", worker, j)
				if err := idx.Add(id, sig); err != nil {
					t.Error(err)
					return
				}
				candidates, err := idx.Candidates(sig)
				if err != nil || !slices.Contains(candidates, id) {
					t.Error("missing concurrent insertion", err)
					return
				}
			}
		})
	}
	wg.Wait()
	candidates, err := idx.Candidates(want)
	if err != nil || len(candidates) != 24*8 || !slices.Equal(candidates, idx.ids) {
		t.Fatal("final insertion order", err)
	}
	// Concurrent attempts at the same ID must produce exactly one insertion.
	var successes int
	var mu sync.Mutex
	for range 16 {
		wg.Go(func() {
			err := idx.Add("same", want)
			if err == nil {
				mu.Lock()
				successes++
				mu.Unlock()
			} else if !errors.Is(err, ErrDuplicateID) {
				t.Error(err)
			}
		})
	}
	wg.Wait()
	if successes != 1 {
		t.Fatalf("duplicate race: %d successes", successes)
	}
}
