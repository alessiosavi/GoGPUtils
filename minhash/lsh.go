package minhash

import (
	"fmt"
	"math"
	"slices"
	"sync"
)

// BandKeys returns one stable key per band. bands must be positive and divide
// the nonzero signature length. Store (signature_version, band, key) rows and
// index (band, key), or (signature_version, band, key); equal (band, key) pairs
// identify candidates. Keys use all 64 bits: store int64(key) in signed BIGINT
// (a lossless bit reinterpretation preserving equality). Isolate incompatible
// seeds, shingle sizes, signature lengths and band layouts in separate namespaces.
func BandKeys(sig Signature, bands int) ([]uint64, error) {
	if len(sig) == 0 || bands < 1 || bands > len(sig) {
		return nil, fmt.Errorf("minhash: bands must be in [1, len(sig)] for a nonempty signature: %w", ErrInvalidParameter)
	}
	if len(sig)%bands != 0 {
		return nil, fmt.Errorf("minhash: bands must divide signature length: %w", ErrInvalidParameter)
	}
	keys := make([]uint64, bands)
	rows := len(sig) / bands
	for j := range keys {
		keys[j] = bandKey(sig[j*rows:(j+1)*rows], j)
	}
	return keys, nil
}

func bandKey(sig Signature, band int) uint64 {
	h := fnvUint64(fnvOffset, uint64(band)) //nolint:gosec // band indices are nonnegative.
	for _, x := range sig {
		h = fnvUint64(h, x)
	}
	return mix64(h)
}

// Index is a small in-memory LSH index, safe for concurrent use. Create it with
// NewIndex; do not copy it. Add and Candidates never retain caller signatures.
// All entries must use the same hashing configuration. Concurrent Adds are
// ordered by acquisition of the write lock, not goroutine launch time.
type Index struct {
	mu          sync.RWMutex
	bands, rows int
	ids         []string
	seen        map[string]struct{}
	buckets     []map[uint64][]int
}

// NewIndex constructs an index with positive bands and rows and at most
// MaxNumPerm total signature positions.
func NewIndex(bands, rows int) (*Index, error) {
	if bands < 1 || rows < 1 || bands > MaxNumPerm/rows {
		return nil, fmt.Errorf("minhash: require positive bands and rows with product <= %d: %w", MaxNumPerm, ErrInvalidParameter)
	}
	idx := &Index{bands: bands, rows: rows, seen: make(map[string]struct{}), buckets: make([]map[uint64][]int, bands)}
	for j := range idx.buckets {
		idx.buckets[j] = make(map[uint64][]int)
	}
	return idx, nil
}

func (idx *Index) validate(sig Signature) error {
	if len(sig) != idx.bands*idx.rows {
		return fmt.Errorf("minhash: signature length %d, want %d: %w", len(sig), idx.bands*idx.rows, ErrLengthMismatch)
	}
	if idx.bands == 0 {
		return fmt.Errorf("minhash: use NewIndex to construct Index: %w", ErrInvalidParameter)
	}
	return nil
}

// Add stores id in each matching band bucket. It returns ErrLengthMismatch for
// a wrong signature length and ErrDuplicateID for an already added id.
func (idx *Index) Add(id string, sig Signature) error {
	if err := idx.validate(sig); err != nil {
		return err
	}
	idx.mu.Lock()
	defer idx.mu.Unlock()
	if _, ok := idx.seen[id]; ok {
		return fmt.Errorf("minhash: id %q: %w", id, ErrDuplicateID)
	}
	seq := len(idx.ids)
	idx.ids = append(idx.ids, id)
	idx.seen[id] = struct{}{}
	for j, bucket := range idx.buckets {
		key := bandKey(sig[j*idx.rows:(j+1)*idx.rows], j)
		bucket[key] = append(bucket[key], seq)
	}
	return nil
}

// Candidates returns IDs sharing at least one band with sig, deduplicated in
// successful insertion order. A wrong length returns ErrLengthMismatch.
// No matches returns nil, nil. Verify candidates with exact Jaccard.
func (idx *Index) Candidates(sig Signature) ([]string, error) {
	if err := idx.validate(sig); err != nil {
		return nil, err
	}
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	seen := make(map[int]struct{})
	var seqs []int
	for j, bucket := range idx.buckets {
		key := bandKey(sig[j*idx.rows:(j+1)*idx.rows], j)
		for _, seq := range bucket[key] {
			if _, ok := seen[seq]; !ok {
				seen[seq] = struct{}{}
				seqs = append(seqs, seq)
			}
		}
	}
	slices.Sort(seqs)
	var ids []string
	for _, seq := range seqs {
		ids = append(ids, idx.ids[seq])
	}
	return ids, nil
}

// Bands chooses a divisor pair balancing integrated false positives and false
// negatives with equal weights around threshold. It uses composite Simpson
// integration with 1000 intervals per integral, considering ascending bands
// and keeping the smaller band count on ties. threshold must be in (0,1).
// For 64/0.85 it chooses 4x16 (~27% recall at J=0.85). Candidate-then-verify
// dedup often needs more bands, e.g. 16x4; see the package documentation.
func Bands(numPerm int, threshold float64) (bands, rows int, err error) {
	if numPerm < 1 || numPerm > MaxNumPerm || !(threshold > 0 && threshold < 1) {
		return 0, 0, fmt.Errorf("minhash: invalid numPerm or threshold outside (0,1): %w", ErrInvalidParameter)
	}
	best := math.Inf(1)
	for b := 1; b <= numPerm; b++ {
		if numPerm%b != 0 {
			continue
		}
		r := numPerm / b
		fp := simpson(0, threshold, func(s float64) float64 { return 1 - math.Pow(1-math.Pow(s, float64(r)), float64(b)) })
		fn := simpson(threshold, 1, func(s float64) float64 { return math.Pow(1-math.Pow(s, float64(r)), float64(b)) })
		if score := 0.5 * (fp + fn); score < best {
			best, bands, rows = score, b, r
		}
	}
	return bands, rows, nil
}

func simpson(lo, hi float64, f func(float64) float64) float64 {
	const n = 1000
	h := (hi - lo) / n
	sum := f(lo) + f(hi)
	for i := 1; i < n; i++ {
		sum += float64(2+2*(i%2)) * f(lo+float64(i)*h)
	}
	return sum * h / 3
}
