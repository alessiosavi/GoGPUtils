package minhash

import (
	"errors"
	"fmt"
	"math/bits"
	"slices"
)

// SignatureVersion identifies the persisted format. Any change to token encoding,
// shingle hashing, permutation derivation, the permutation formula, or band keys
// MUST bump it. Store it alongside persisted data.
const SignatureVersion = 1

// MaxNumPerm is the largest supported number of permutations.
const MaxNumPerm = 1024

var (
	// ErrEmptySet reports that a signature cannot be made from zero shingles.
	ErrEmptySet = errors.New("minhash: empty set")
	// ErrInvalidParameter reports an invalid configuration or empty signature.
	ErrInvalidParameter = errors.New("minhash: invalid parameter")
	// ErrLengthMismatch reports incompatible signature or destination lengths.
	ErrLengthMismatch = errors.New("minhash: length mismatch")
	// ErrDuplicateID reports an ID already present in an Index.
	ErrDuplicateID = errors.New("minhash: duplicate id")
)

// Signature contains one minimum per permutation. Values produced by Hasher
// are less than 2^61-1 and fit in signed int64 / SQL BIGINT unchanged.
type Signature []uint64

// Shingles hashes consecutive windows of k tokens and returns sorted, unique
// hashes. k must be positive. Empty input returns nil, nil; fewer than k tokens
// produce one shingle containing all tokens. Tokens are raw bytes without
// normalization: NFC "è" and NFD "e\u0300" differ. Empty-string tokens are legal.
func Shingles(tokens []string, k int) ([]uint64, error) {
	if k < 1 {
		return nil, fmt.Errorf("minhash: k must be positive: %w", ErrInvalidParameter)
	}
	if len(tokens) == 0 {
		return nil, nil
	}
	k = min(k, len(tokens))
	out := make([]uint64, len(tokens)-k+1)
	for i := range out {
		h := fnvOffset
		for _, token := range tokens[i : i+k] {
			n := uint64(len(token)) //nolint:gosec // string lengths are nonnegative.
			for n >= 128 {
				h = (h ^ ((n & 127) | 128)) * fnvPrime
				n >>= 7
			}
			h = (h ^ n) * fnvPrime
			h = fnvBytes(h, token)
		}
		out[i] = mix64(h)
	}
	slices.Sort(out)
	return slices.Compact(out), nil
}

type permutation struct{ a, b uint64 }

// Hasher is immutable after New and safe for concurrent use. Construct it with
// New; the zero value cannot sign. Callers must not concurrently modify inputs
// or share writable SignInto destinations.
type Hasher struct{ permutations []permutation }

// New derives numPerm permutations deterministically from seed. numPerm must
// be in [1, MaxNumPerm]. A larger hasher with the same seed shares the prefix.
func New(numPerm int, seed uint64) (*Hasher, error) {
	if numPerm < 1 || numPerm > MaxNumPerm {
		return nil, fmt.Errorf("minhash: numPerm must be in [1, %d]: %w", MaxNumPerm, ErrInvalidParameter)
	}
	next := func() uint64 { seed += 0x9e3779b97f4a7c15; return mix64(seed) >> 3 }
	h := &Hasher{permutations: make([]permutation, numPerm)}
	for i := range h.permutations {
		a := next()
		for a == 0 || a == prime {
			a = next()
		}
		b := next()
		for b == prime {
			b = next()
		}
		h.permutations[i] = permutation{a, b}
	}
	return h, nil
}

// NumPerm returns the signature length.
func (h *Hasher) NumPerm() int { return len(h.permutations) }

// Sign allocates a signature and delegates to SignInto. Empty shingles return
// ErrEmptySet, never an all-max signature.
func (h *Hasher) Sign(shingles []uint64) (Signature, error) {
	dst := make(Signature, h.NumPerm())
	if err := h.SignInto(dst, shingles); err != nil {
		return nil, err
	}
	return dst, nil
}

// SignInto fills dst without allocating on success. len(dst) must equal NumPerm;
// shingles must be nonempty, but need not be sorted or unique. On error dst is
// untouched. Overlap between dst and shingles is supported.
func (h *Hasher) SignInto(dst Signature, shingles []uint64) error {
	if len(dst) != h.NumPerm() {
		return fmt.Errorf("minhash: dst length %d, want %d: %w", len(dst), h.NumPerm(), ErrLengthMismatch)
	}
	if h.NumPerm() == 0 {
		return fmt.Errorf("minhash: use New to construct Hasher: %w", ErrInvalidParameter)
	}
	if len(shingles) == 0 {
		return ErrEmptySet
	}
	// Stack scratch permits overlap without copying the unbounded input.
	var scratch [MaxNumPerm]uint64
	mins := scratch[:h.NumPerm()]
	for i := range mins {
		mins[i] = prime
	}
	for _, shingle := range shingles {
		x := reduce(shingle)
		for i, p := range h.permutations {
			mins[i] = min(mins[i], permute(p.a, x, p.b))
		}
	}
	copy(dst, mins)
	return nil
}

// permute requires a, x, b < p. The 122-bit product folds into at most 2p-2;
// adding b stays below 3p, so no uint64 addition can overflow.
func permute(a, x, b uint64) uint64 {
	hi, lo := bits.Mul64(a, x)
	return reduce((lo & prime) + (lo >> 61) + (hi << 3) + b)
}

// Similarity estimates Jaccard as the fraction of equal positions. It is only
// meaningful with the same numPerm, seed, SignatureVersion and shingle k.
// Unequal lengths return ErrLengthMismatch; empty signatures are invalid.
func Similarity(a, b Signature) (float64, error) {
	if len(a) != len(b) {
		return 0, fmt.Errorf("minhash: signature lengths differ: %w", ErrLengthMismatch)
	}
	if len(a) == 0 {
		return 0, fmt.Errorf("minhash: empty signature: %w", ErrInvalidParameter)
	}
	equal := 0
	for i, x := range a {
		if x == b[i] {
			equal++
		}
	}
	return float64(equal) / float64(len(a)), nil
}

// Jaccard returns exact intersection/union for sorted, deduplicated sets (as
// returned by Shingles). Both empty returns 0: empties are never duplicates.
// If the precondition is violated the result is unspecified, but remains [0,1].
func Jaccard(a, b []uint64) float64 {
	i, j, intersection := 0, 0, 0
	for i < len(a) && j < len(b) {
		switch {
		case a[i] < b[j]:
			i++
		case a[i] > b[j]:
			j++
		default:
			intersection++
			i++
			j++
		}
	}
	// Floating-point addition avoids int overflow for very large slices.
	union := float64(len(a)-intersection) + float64(len(b))
	if union == 0 {
		return 0
	}
	return float64(intersection) / union
}
