package minhash

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"testing"
)

// randomPair draws disjoint random shared/unique sets with J=c/(c+2u).
func randomPair(rng *rand.Rand, c, u int) ([]uint64, []uint64) {
	seen := map[uint64]bool{}
	values := make([]uint64, 0, c+2*u)
	for len(values) < cap(values) {
		x := rng.Uint64()
		if !seen[x] {
			seen[x] = true
			values = append(values, x)
		}
	}
	a := slices.Clone(values[:c+u])
	b := append(slices.Clone(values[:c]), values[c+u:]...)
	slices.Sort(a)
	slices.Sort(b)
	return a, b
}

func TestStatisticalBounds(t *testing.T) {
	const trials = 200 // Per J/seed/numPerm configuration, also under -race.
	for _, n := range []int{64, 128} {
		for _, seed := range []uint64{0, 42, 2026} {
			for _, tenth := range []int{1, 3, 5, 7, 9} {
				t.Run(fmt.Sprintf("perm=%d/seed=%d/J=%d", n, seed, tenth), func(t *testing.T) {
					h := mustNew(t, n, seed)
					rng := rand.New(rand.NewPCG(123+seed, uint64(tenth)))
					within := 0
					for range trials {
						a, b := randomPair(rng, 2*tenth, 10-tenth)
						j := float64(tenth) / 10
						if Jaccard(a, b) != j {
							t.Fatal("incorrect fixture Jaccard")
						}
						estimate, err := Similarity(mustSign(t, h, a), mustSign(t, h, b))
						if err != nil {
							t.Fatal(err)
						}
						if math.Abs(estimate-j) <= 3/math.Sqrt(float64(n)) {
							within++
						}
					}
					if within < 198 {
						t.Fatalf("only %d/%d within bound", within, trials)
					}
				})
			}
		}
	}
}

func TestLSHBehaviour(t *testing.T) {
	// Measured with these fixed PCG streams (1000 pairs each): J=0.9 recall
	// 1000/1000 = 1.000; J=0.3 false-candidate rate 126/1000 = 0.126.
	// BandKeys and Index agree on every pair, not just the aggregate rates.
	const trials = 1000
	h := mustNew(t, 64, 42)
	for _, tenth := range []int{9, 3} {
		rng := rand.New(rand.NewPCG(8181, uint64(tenth)))
		hits := 0
		for range trials {
			a, b := randomPair(rng, 2*tenth, 10-tenth)
			sa, sb := mustSign(t, h, a), mustSign(t, h, b)
			ka, err := BandKeys(sa, 16)
			if err != nil {
				t.Fatal(err)
			}
			kb, err := BandKeys(sb, 16)
			if err != nil {
				t.Fatal(err)
			}
			match := false
			for i := range ka {
				match = match || ka[i] == kb[i]
			}
			idx, err := NewIndex(16, 4)
			if err != nil {
				t.Fatal(err)
			}
			if err := idx.Add("a", sa); err != nil {
				t.Fatal(err)
			}
			candidates, err := idx.Candidates(sb)
			if err != nil || (len(candidates) == 1) != match {
				t.Fatal("Index/BandKeys disagree", err)
			}
			if match {
				hits++
			}
		}
		rate := float64(hits) / trials
		t.Logf("J=%.1f: BandKeys and Index %d/%d = %.4f", float64(tenth)/10, hits, trials, rate)
		if tenth == 9 && rate < 0.95 || tenth == 3 && rate > 0.20 {
			t.Fatalf("LSH rate %.4f out of bounds", rate)
		}
	}
}
