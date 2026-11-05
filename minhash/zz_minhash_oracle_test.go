package minhash

import (
	"encoding/binary"
	"hash/fnv"
	"math"
	"math/big"
	"math/rand/v2"
	"slices"
	"strings"
	"testing"
)

func oracleMix(z uint64) uint64 {
	z ^= z >> 30
	z *= 0xbf58476d1ce4e5b9
	z ^= z >> 27
	z *= 0x94d049bb133111eb
	return z ^ (z >> 31)
}

func oracleHash(data []byte) uint64 {
	h := fnv.New64a()
	_, _ = h.Write(data)
	return h.Sum64()
}

func oracleShingles(tokens []string, k int) []uint64 {
	if len(tokens) == 0 {
		return nil
	}
	k = min(k, len(tokens))
	set := map[uint64]bool{}
	for i := 0; i+k <= len(tokens); i++ {
		var enc []byte
		for _, token := range tokens[i : i+k] {
			enc = binary.AppendUvarint(enc, uint64(len(token)))
			enc = append(enc, token...)
		}
		set[oracleMix(oracleHash(enc))] = true
	}
	out := make([]uint64, 0, len(set))
	for x := range set {
		out = append(out, x)
	}
	slices.Sort(out)
	return out
}

func oraclePermutation(a, x, b uint64) uint64 {
	p := new(big.Int).SetUint64(1<<61 - 1)
	n := new(big.Int).Mul(new(big.Int).SetUint64(a), new(big.Int).SetUint64(x))
	n.Add(n, new(big.Int).SetUint64(b))
	return n.Mod(n, p).Uint64()
}

func oracleSign(h *Hasher, xs []uint64) Signature {
	out := make(Signature, h.NumPerm())
	for i, p := range h.permutations {
		out[i] = math.MaxUint64
		for _, x := range xs {
			out[i] = min(out[i], oraclePermutation(p.a, x, p.b))
		}
	}
	return out
}

func oracleJaccard(a, b []uint64) float64 {
	aa, bb := map[uint64]bool{}, map[uint64]bool{}
	for _, x := range a {
		aa[x] = true
	}
	for _, x := range b {
		bb[x] = true
	}
	inter := 0
	for x := range aa {
		if bb[x] {
			inter++
		}
	}
	if len(aa)+len(bb) == 0 {
		return 0
	}
	return float64(inter) / float64(len(aa)+len(bb)-inter)
}

func TestArithmeticOracle(t *testing.T) {
	edges := []uint64{0, 1, prime - 1, prime, math.MaxUint64}
	for _, x := range edges {
		if reduce(x) != x%prime {
			t.Fatalf("reduce(%d)", x)
		}
		for _, a := range edges[:3] {
			for _, b := range edges[:3] {
				if got, want := permute(a, reduce(x), b), oraclePermutation(a, x, b); got != want {
					t.Fatalf("a=%d x=%d b=%d: %d != %d", a, x, b, got, want)
				}
			}
		}
	}
	rng := rand.New(rand.NewPCG(7, 11))
	for range 10000 {
		a, x, b := rng.Uint64()%prime, rng.Uint64(), rng.Uint64()%prime
		if reduce(x) != x%prime || permute(a, reduce(x), b) != oraclePermutation(a, x, b) {
			t.Fatalf("arithmetic: %d %d %d", a, x, b)
		}
	}
}

func TestOracles(t *testing.T) {
	rng := rand.New(rand.NewPCG(17, 29))
	words := []string{"", "è", "e\u0300", "東京", "a", "a", "20%", strings.Repeat("x", 128), strings.Repeat("y", 16384)}
	for trial := range 100 {
		tokens := make([]string, rng.IntN(25))
		for i := range tokens {
			tokens[i] = words[rng.IntN(len(words))]
		}
		k := 1 + rng.IntN(8)
		got, err := Shingles(tokens, k)
		if err != nil || !slices.Equal(got, oracleShingles(tokens, k)) {
			t.Fatalf("shingles trial %d: %v", trial, err)
		}
		if fnvBytes(fnvOffset, strings.Join(tokens, "")) != oracleHash([]byte(strings.Join(tokens, ""))) {
			t.Fatal("FNV differs from hash/fnv")
		}
		h := mustNew(t, 1+rng.IntN(32), rng.Uint64())
		if len(got) == 0 {
			continue
		}
		sig := mustSign(t, h, got)
		if !slices.Equal(sig, oracleSign(h, got)) {
			t.Fatalf("signature trial %d", trial)
		}
		for bands := 1; bands <= len(sig); bands++ {
			if len(sig)%bands != 0 {
				continue
			}
			keys, err := BandKeys(sig, bands)
			if err != nil {
				t.Fatal(err)
			}
			r := len(sig) / bands
			for j, key := range keys {
				enc := binary.LittleEndian.AppendUint64(nil, uint64(j))
				for _, x := range sig[j*r : (j+1)*r] {
					enc = binary.LittleEndian.AppendUint64(enc, x)
				}
				if key != oracleMix(oracleHash(enc)) {
					t.Fatal("band key oracle mismatch")
				}
			}
		}
		other, err := Shingles(tokens, 1)
		if err != nil {
			t.Fatal(err)
		}
		if Jaccard(got, other) != oracleJaccard(got, other) {
			t.Fatal("Jaccard oracle mismatch")
		}
	}
}

// Independent midpoint quadrature (not Simpson), including prime numPerm.
func TestBandsOracle(t *testing.T) {
	for _, n := range []int{1, 7, 16, 64, 128} {
		for _, threshold := range []float64{0.1, 0.5, 0.7, 0.8, 0.85, 0.9, 0.99} {
			best, wantB := math.Inf(1), 0
			for b := 1; b <= n; b++ {
				if n%b != 0 {
					continue
				}
				score := 0.0
				const steps = 4000
				for i := range steps {
					x := (float64(i) + 0.5) / steps
					// Compute integer powers by repeated multiplication, independently
					// of the production math.Pow path.
					probability := func(s float64) float64 {
						row := 1.0
						for range n / b {
							row *= s
						}
						miss := 1.0
						for range b {
							miss *= 1 - row
						}
						return miss
					}
					score += threshold*(1-probability(x*threshold)) + (1-threshold)*probability(threshold+x*(1-threshold))
				}
				if score < best {
					best, wantB = score, b
				}
			}
			b, r, err := Bands(n, threshold)
			if err != nil || b != wantB || b*r != n {
				t.Fatalf("Bands(%d,%g) = %dx%d %v; want %d bands", n, threshold, b, r, err, wantB)
			}
		}
	}
}
