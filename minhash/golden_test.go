package minhash

import (
	"slices"
	"strings"
	"testing"
)

func TestGolden(t *testing.T) {
	type seeded struct {
		prefix       Signature
		keys         []uint64
		permutations []permutation
	}
	type vectors struct {
		shingles     map[int][]uint64
		longShingles map[int][]uint64
		seeds        map[uint64]seeded
	}
	// Persisted v1 literals: never regenerate these to accommodate a change.
	// Bump SignatureVersion and add a new entry while retaining the old vectors.
	golden := map[int]vectors{
		1: {
			shingles: map[int][]uint64{
				1:  {0x1473af777f100235, 0x25fc6dd36ce04b20, 0x2afe41a27620d808, 0x4c516bd2b0d9620b, 0x514e766b853e15f6, 0x67532d91b1124643, 0x6b1c27841a51dcb5, 0x924db29458a9bb28, 0xddd1354f6030b3fe, 0xe3cf63f073cd53ff},
				3:  {0x32ad6ca5cd5484a7, 0x34a3620a21071c14, 0x69e4a190376693be, 0xa5f7bb9ffa8d0d20, 0xc8422f42b00a5b67, 0xe021d87b3cddb0f4, 0xe0cfd21c18ce6a80, 0xe6ef20002456c71e},
				20: {0x62ffaaf20edb78aa},
			},
			// UTF-8 tokens with two- and three-byte uvarint length prefixes.
			longShingles: map[int][]uint64{
				70:   {0xc23f760b490595ae},
				8192: {0x5d54540ac171994d},
			},
			seeds: map[uint64]seeded{
				42: {
					prefix:       Signature{0x3235ca8c7598996, 0x3e34f4d04d64150, 0x42fad0d55b895e7, 0x49ee3d43414f92, 0x188fcbe90da1e8b, 0x6e1552d434ee99c, 0x3e834b406d8801, 0xb2e93d9b4e77ff},
					keys:         []uint64{0x81cb0d0231b5a0ec, 0x4697aa9812b37382, 0x28727825ce7e0820, 0xb46cca5aa248de53, 0xf734c15617a7339a, 0x8fecccf683ef07a2, 0xe72b603834a0aab5, 0xec26a7199a6981c7, 0x7bd4eb397e3f51fe, 0xcd6ecf3083d60102, 0x459c5a72ee191374, 0xd741ec4834d32d0e, 0x570f570b1d58da8, 0xeeef007ae38f831, 0x42a1c7e143dd271c, 0x3f6e62752a314a34},
					permutations: []permutation{{1709932191594409426, 368728261515861536}, {642410468557845482, 793649757532281970}, {87691598267745406, 2001997640707873632}, {503608089097203115, 1846131415873471988}},
				},
				0: {
					prefix:       Signature{0x361c27030b95d7f, 0x6348fa053c2581, 0x5910594b2f77f7a, 0x53c9eb96982c1c9, 0xa54479454a309f4, 0x5e7bd61243f1687, 0x66499f1d8a2d9fb, 0x65b36ae4d8a2f8d},
					keys:         []uint64{0x4025937a874613ef, 0x9d34a6bd9ecd9aaf, 0x41dee562458cfafe, 0x468b1836fca77a58, 0x47a258aa271e0a3, 0x80516efc9d99440f, 0x91934a9863b25d71, 0x4b012307e7c45587, 0x578345e7d025124b, 0x4401e1b3c2968e4c, 0x3f2746eaee1c4475, 0xb3ef35311212226a, 0x2c94c24086a80188, 0x161b40278b1f7931, 0x6f9edfcf2a785c4e, 0x8fe254f7f263c261},
					permutations: []permutation{{2036776052082325941, 995035815274294462}, {60952127433943209, 2238701422097567805}, {245218775303261843, 754761825157895261}, {400912003250038364, 1779065233200043367}},
				},
			},
		},
	}
	v, ok := golden[SignatureVersion]
	if !ok {
		t.Fatal("missing golden vectors: bump SignatureVersion and add new vectors; never edit v1")
	}
	tokens := []string{"perché", "città", "caffè", "è", "e\u0300", "4.5", "20%", "1,000", "", "fine"}
	for k, want := range v.shingles {
		got, err := Shingles(tokens, k)
		if err != nil || !slices.Equal(got, want) {
			t.Errorf("v%d shingles k=%d changed: got %x; want %x; %v", SignatureVersion, k, got, want, err)
		}
	}
	for repeats, want := range v.longShingles {
		tokens := []string{strings.Repeat("ò", repeats), "fine"}
		got, err := Shingles(tokens, 2)
		if err != nil || !slices.Equal(got, want) {
			t.Errorf("v%d shingles tokenBytes=%d changed: got %x; want %x; %v", SignatureVersion, len(tokens[0]), got, want, err)
		}
		if oracle := oracleShingles(tokens, 2); !slices.Equal(oracle, want) {
			t.Errorf("v%d shingles tokenBytes=%d oracle mismatch: got %x; want %x", SignatureVersion, len(tokens[0]), oracle, want)
		}
	}
	xs, err := Shingles(tokens, 3)
	if err != nil {
		t.Fatal(err)
	}
	for seed, want := range v.seeds {
		h := mustNew(t, 64, seed)
		if !slices.Equal(h.permutations[:len(want.permutations)], want.permutations) {
			t.Errorf("v%d seed=%d permutation derivation changed", SignatureVersion, seed)
		}
		sig := mustSign(t, h, xs)
		if !slices.Equal(sig[:8], want.prefix) {
			t.Errorf("v%d seed=%d signature prefix changed: got %x", SignatureVersion, seed, sig[:8])
		}
		keys, err := BandKeys(sig, 16)
		if err != nil || !slices.Equal(keys, want.keys) {
			t.Errorf("v%d seed=%d band keys changed: got %x; %v", SignatureVersion, seed, keys, err)
		}
	}
}
