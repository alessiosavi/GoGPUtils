package minhash_test

import (
	"fmt"

	"github.com/alessiosavi/GoGPUtils/minhash"
	"github.com/alessiosavi/GoGPUtils/textnorm"
)

func Example() {
	texts := []string{
		"Il caffè della città è disponibile online con consegna gratuita. Our shop offers fresh coffee every morning with helpful service and delivery today.",
		"Il caffè della città è disponibile online con consegna gratuita. Our shop offers fresh coffee every morning with helpful service and delivery tomorrow.",
	}
	h, err := minhash.New(64, 42)
	if err != nil {
		fmt.Println(err)
		return
	}
	var sets [][]uint64
	var signatures []minhash.Signature
	for _, text := range texts {
		tokens, err := textnorm.MeaningPreset().SplitTokens().Run(text)
		if err != nil {
			fmt.Println(err)
			return
		}
		shingles, err := minhash.Shingles(tokens, 5)
		if err != nil {
			fmt.Println(err)
			return
		}
		sig, err := h.Sign(shingles)
		if err != nil {
			fmt.Println(err)
			return
		}
		sets = append(sets, shingles)
		signatures = append(signatures, sig)
	}
	estimate, err := minhash.Similarity(signatures[0], signatures[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("estimate: %.3f; exact: %.3f\n", estimate, minhash.Jaccard(sets[0], sets[1]))
	keys, err := minhash.BandKeys(signatures[0], 16)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("persisted bands:", len(keys))
	idx, err := minhash.NewIndex(16, 4)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := idx.Add("coffee-page", signatures[0]); err != nil {
		fmt.Println(err)
		return
	}
	candidates, err := idx.Candidates(signatures[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("candidates:", candidates)
	// Fetch candidate shingle sets and verify Jaccard >= 0.85 before deduping.
	fmt.Println("verified duplicate:", minhash.Jaccard(sets[0], sets[1]) >= 0.85)
	// Output:
	// estimate: 0.859; exact: 0.900
	// persisted bands: 16
	// candidates: [coffee-page]
	// verified duplicate: true
}

func ExampleBands() {
	bands, rows, err := minhash.Bands(64, 0.85)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Balanced error; use 16x4 for higher candidate recall, then verify Jaccard.
	fmt.Println(bands, rows)
	// Output: 4 16
}

func ExampleBandKeys() {
	h, err := minhash.New(64, 42)
	if err != nil {
		fmt.Println(err)
		return
	}
	sig, err := h.Sign([]uint64{1, 2, 3})
	if err != nil {
		fmt.Println(err)
		return
	}
	keys, err := minhash.BandKeys(sig, 16)
	if err != nil {
		fmt.Println(err)
		return
	}
	stored := int64(keys[0]) // Preserve all bits in a signed SQL BIGINT.
	fmt.Println("version:", minhash.SignatureVersion, "bands:", len(keys))
	fmt.Println("round trip:", uint64(stored) == keys[0])
	// Output:
	// version: 1 bands: 16
	// round trip: true
}

func ExampleHasher_SignInto() {
	h, err := minhash.New(64, 42)
	if err != nil {
		fmt.Println(err)
		return
	}
	dst := make(minhash.Signature, h.NumPerm())
	if err := h.SignInto(dst, []uint64{1, 2, 3}); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("positions:", len(dst))
	// Output: positions: 64
}

func ExampleShingles() {
	set, err := minhash.Shingles([]string{"caffè", "in", "città"}, 5)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("short text shingles:", len(set))
	empty, err := minhash.Shingles(nil, 5)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("empty text shingles:", len(empty))
	// Output:
	// short text shingles: 1
	// empty text shingles: 0
}
