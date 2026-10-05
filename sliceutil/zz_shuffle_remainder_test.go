package sliceutil

import (
	"fmt"
	"slices"
	"testing"
)

// Keep reference state local and take the remainder in uint64 so the oracle
// describes the existing 64-bit permutation on every architecture.
func shuffleReferenceDraw(state *uint64) uint64 {
	*state ^= *state << 13
	*state ^= *state >> 7
	*state ^= *state << 17
	return *state & 0xffffffff
}

func shuffleReferenceInPlace(s []int, state *uint64) {
	for i := len(s) - 1; i > 0; i-- {
		j := int(shuffleReferenceDraw(state) % uint64(i+1))
		s[i], s[j] = s[j], s[i]
	}
}

func TestShuffleUnsignedRemainder(t *testing.T) {
	saved, _ := states()
	t.Cleanup(func() { SeedShuffle(saved) })
	const highBitSeed = uint64(2)
	state := highBitSeed
	if draw := shuffleReferenceDraw(&state); draw != 0x81044082 {
		t.Fatalf("high-bit fixture: first draw=%#x, want 0x81044082", draw)
	}
	for _, seed := range []uint64{0, 1, highBitSeed, 42} {
		// -1 represents nil; zero represents a non-nil empty slice.
		for _, n := range []int{-1, 0, 1, 2, 16, 17, 32, 33} {
			for _, copying := range []bool{false, true} {
				t.Run(fmt.Sprintf("seed=%d/n=%d/copy=%t", seed, n, copying), func(t *testing.T) {
					var input []int
					if n >= 0 {
						input = make([]int, n)
						for i := range input {
							input[i] = i
						}
					}
					original := slices.Clone(input)
					want := slices.Clone(input)
					wantState := seed
					if wantState == 0 {
						wantState = 1
					}
					shuffleReferenceInPlace(want, &wantState)
					SeedShuffle(seed)
					got, panicked, msg := outcome(func() []int {
						if copying {
							return Shuffle(input)
						}
						ShuffleInPlace(input)
						return input
					})
					if panicked {
						t.Fatalf("shuffle panicked: %s", msg)
					}
					if !slices.Equal(got, want) || (got == nil) != (original == nil) {
						t.Errorf("permutation=%v (nil=%t), want %v (nil=%t)", got, got == nil, want, original == nil)
					}
					gotState, _ := states()
					if gotState != wantState {
						t.Errorf("state=%#x, want %#x after %d draws", gotState, wantState, max(n-1, 0))
					}
					if copying {
						if !slices.Equal(input, original) || (input == nil) != (original == nil) {
							t.Error("Shuffle modified its input")
						}
						if len(got) > 0 {
							got[0] = -1
							if !slices.Equal(input, original) {
								t.Error("Shuffle result aliases its input")
							}
						}
					}
				})
			}
		}
	}
}

func TestShuffleGoldenPermutations(t *testing.T) {
	saved, _ := states()
	t.Cleanup(func() { SeedShuffle(saved) })
	// Captured from unmodified master 82496f7 on a 64-bit runtime.
	cases := []struct {
		seed uint64
		want []int
	}{
		{0, []int{13, 14, 15, 2, 3, 12, 4, 10, 0, 5, 7, 9, 6, 11, 8, 1}},
		{2, []int{14, 11, 7, 6, 2, 0, 9, 10, 1, 8, 13, 15, 5, 12, 16, 3, 4}},
		{42, []int{7, 1, 2, 4, 11, 14, 9, 0, 13, 3, 6, 15, 5, 12, 8, 10}},
		{42, []int{10, 21, 26, 15, 28, 13, 27, 32, 8, 24, 3, 2, 25, 30, 16, 17, 9, 19, 1, 29, 14, 6, 23, 5, 0, 20, 11, 7, 12, 22, 4, 31, 18}},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("seed=%d/n=%d", tc.seed, len(tc.want)), func(t *testing.T) {
			got := make([]int, len(tc.want))
			for i := range got {
				got[i] = i
			}
			reference := slices.Clone(got)
			state := tc.seed
			if state == 0 {
				state = 1
			}
			shuffleReferenceInPlace(reference, &state)
			if !slices.Equal(reference, tc.want) {
				t.Fatalf("reference=%v, want golden %v", reference, tc.want)
			}
			SeedShuffle(tc.seed)
			if panicked, _, msg := panicDetails(func() { ShuffleInPlace(got) }); panicked {
				t.Fatalf("shuffle panicked: %s", msg)
			}
			if !slices.Equal(got, tc.want) {
				t.Fatalf("permutation=%v, want golden %v", got, tc.want)
			}
		})
	}
}
