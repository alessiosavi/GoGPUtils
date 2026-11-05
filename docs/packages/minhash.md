---
title: minhash
parent: Packages
nav_order: 10
---

# Package `minhash`

Deterministic word shingles, MinHash signatures, and locality-sensitive hashing (LSH) for near-duplicate text. The production package imports only the standard library. Hashes are noncryptographic and can collide; use exact Jaccard of shingle sets to verify candidates.

```go
import "github.com/alessiosavi/GoGPUtils/minhash"
```

## API

```go
const SignatureVersion = 1
const MaxNumPerm = 1024

type Signature []uint64

func Shingles(tokens []string, k int) ([]uint64, error)
func New(numPerm int, seed uint64) (*Hasher, error)
func (h *Hasher) NumPerm() int
func (h *Hasher) Sign(shingles []uint64) (Signature, error)
func (h *Hasher) SignInto(dst Signature, shingles []uint64) error
func Similarity(a, b Signature) (float64, error)
func Jaccard(a, b []uint64) float64
func BandKeys(sig Signature, bands int) ([]uint64, error)
func Bands(numPerm int, threshold float64) (bands, rows int, err error)
func NewIndex(bands, rows int) (*Index, error)
func (idx *Index) Add(id string, sig Signature) error
func (idx *Index) Candidates(sig Signature) ([]string, error)
```

| Operation | Contract |
|---|---|
| `Shingles` | k must be positive. Returns sorted, unique hashes. Empty tokens input returns `nil, nil`; nonempty input shorter than k produces one shingle over all tokens. Empty-string tokens are legal. |
| `New` | numPerm must be 1..1024. Same seed and size produce the same permutations; a larger size shares the smaller one's prefix. |
| `Sign` | Rejects zero shingles. Accepts unsorted input with duplicates. Returns a newly allocated signature. |
| `SignInto` | Destination must have exactly NumPerm positions. Zero allocations on success; errors leave the destination untouched. Overlapping input and output are supported. |
| `Similarity` | Fraction of equal positions. Requires nonempty, equal-length signatures from the same seed, version, k and normalization. |
| `Jaccard` | Linear merge of sorted, unique sets. Both empty returns 0: empties are never duplicates. Otherwise intersection/union. With invalid ordering or duplicates the value is unspecified but still in [0,1]. |
| `BandKeys` | Positive bands must divide a nonempty signature length. One full-width 64-bit key per band. |
| `NewIndex` | Positive bands and rows, with product at most 1024. |
| `Index.Add` | Unique ID, signature length bands×rows. Empty-string IDs are legal. Does not retain the signature. |
| `Index.Candidates` | Signature length bands×rows. Deduplicated IDs sharing at least one band, in successful insertion order. No matches returns `nil, nil`. |
| `Bands` | Valid numPerm and finite threshold strictly between 0 and 1. Equal-weight integrated error minimization, not a guaranteed recall target. |

Match errors with `errors.Is`: `ErrEmptySet` for signing zero shingles; `ErrInvalidParameter` for invalid configuration or a required empty signature; `ErrLengthMismatch` for incompatible signature/destination lengths; `ErrDuplicateID` for an ID already in the index. Construct `Hasher` and `Index` through their constructors; their zero values cannot sign/index.

## Text preprocessing and example

No normalization occurs inside minhash. Tokens are raw bytes: NFC `è` and NFD `e\u0300` differ. Choose and keep a consistent normalization pipeline, such as `textnorm.MeaningPreset().SplitTokens()`. Length-prefixing distinguishes `["ab", "c"]` from `["a", "bc"]` and preserves empty-string tokens.

```go
tokens, err := textnorm.MeaningPreset().SplitTokens().Run(text)
if err != nil { return err }
shingles, err := minhash.Shingles(tokens, 5)
if err != nil { return err }
h, err := minhash.New(64, 42)
if err != nil { return err }
sig, err := h.Sign(shingles)
if err != nil { return err } // Empty text: ErrEmptySet.

idx, err := minhash.NewIndex(16, 4)
if err != nil { return err }
if err := idx.Add("document-1", sig); err != nil { return err }
candidates, err := idx.Candidates(sig)
if err != nil { return err }
// candidates == []string{"document-1"}
// For a new document, query its signature, fetch candidate shingle sets,
// and require minhash.Jaccard(newSet, candidateSet) >= 0.85 before deduping.
```

Reuse a destination with `dst := make(minhash.Signature, h.NumPerm())` followed by `h.SignInto(dst, shingles)` for allocation-free signing. See the [runnable examples](https://pkg.go.dev/github.com/alessiosavi/GoGPUtils/minhash#pkg-examples): the Italian/English coffee-page pair has estimated similarity 0.859 and exact Jaccard 0.900, and the index retrieves its candidate.

## Persisted format version 1

All integer operations wrap as unsigned 64-bit arithmetic unless explicitly reduced modulo p.

1. Encode every token in a window as `binary.AppendUvarint(uint64(len(token)))` followed by its raw bytes, in order. Length counts bytes, not Unicode runes. There is no window-length header or separator.
2. Hash those bytes with FNV-1a 64: initialize h to **14695981039346656037**, then for every byte compute `(h XOR byte) * 1099511628211`. The implementation hashes incrementally without constructing the encoding.
3. Apply the splitmix64 finalizer to that FNV hash:

   ```text
   z = (z XOR (z >> 30)) * 0xbf58476d1ce4e5b9
   z = (z XOR (z >> 27)) * 0x94d049bb133111eb
   return z XOR (z >> 31)
   ```

4. Use p = **2^61−1**. Start generator state at the caller's seed. Each `next()` first adds **0x9e3779b97f4a7c15** to state, then returns `mix64(state)`. A draw is `next() >> 3`. For each permutation i in ascending order, draw a until **1 ≤ a < p**, then draw b until **0 ≤ b < p**. Rejected draws consume generator state.
5. Reduce every shingle x modulo p. Evaluate **(a·x+b) mod p** at full precision; keep the minimum for each permutation. Production code uses `bits.Mul64` and Mersenne reduction; tests independently check against `math/big`. Signature values are below p and fit in signed int64 / SQL BIGINT unchanged.
6. For band j with r rows, encode **uint64LE(j) ‖ uint64LE(sig[j·r]) ‖ … ‖ uint64LE(sig[j·r+r−1])**. The key is **mix64(FNV-1a-64(encoded bytes))**. The zero-based band index is part of the hash.

`SignatureVersion` must change if token encoding, shingle hashing, permutation derivation, the permutation formula, or band keys change. Add new golden vectors under the new version while retaining the v1 literals. Existing persisted records need recomputation from their source text before comparison with a new format. The version identifies the algorithm; it does **not** identify the caller's configuration.

## Persistence recipe

Store `(signature_version, band, key)` for every document and index `(band, key)` or `(signature_version, band, key)`. Scope queries to one compatible configuration: **seed, shingle k, numPerm, normalization, and bands/rows** as well as version. A namespace can identify that configuration. Store signatures and shingle sets with the document so retrieved candidates can be verified without reprocessing text.

Band keys use the full uint64 range. In Go, `int64(key)` reinterprets the bits for signed BIGINT; `uint64(stored)` recovers them. Negative stored keys are valid and equality is preserved. Signatures need no such reinterpretation because their values are below 2^61−1.

A generic SQL sketch (bind parameters using your database driver's syntax):

```sql
-- namespace identifies a fixed seed, k, numPerm, normalization and band layout.
CREATE TABLE document_band (
    namespace VARCHAR(128) NOT NULL,
    signature_version INTEGER NOT NULL,
    document_id VARCHAR(256) NOT NULL,
    band INTEGER NOT NULL,
    band_key BIGINT NOT NULL,
    PRIMARY KEY (namespace, signature_version, document_id, band)
);
CREATE INDEX document_band_lookup
    ON document_band (namespace, signature_version, band, band_key);

-- Populate query_band with the new document's (band, band_key) values.
SELECT DISTINCT stored.document_id
FROM document_band AS stored
JOIN query_band AS query
  ON stored.band = query.band AND stored.band_key = query.band_key
WHERE stored.namespace = :namespace
  AND stored.signature_version = :signature_version;
```

SQL result order requires an explicit `ORDER BY` if needed; only the in-memory Index promises insertion order. Fetch each candidate's sorted shingle set and apply `Jaccard` before accepting a duplicate. LSH can miss duplicates; exact verification cannot recover a candidate that was never retrieved.

## Choosing bands and rows

For similarity s, b bands and r rows, the ideal candidate probability is the S-curve:

```text
P(s) = 1 − (1 − s^r)^b
FP = integral from 0 to threshold of P(s) ds
FN = integral from threshold to 1 of (1 − P(s)) ds
balanced error = 0.5·FP + 0.5·FN
```

These integrals weight the similarity axis uniformly; they are not measured error rates on a particular corpus. `Bands` enumerates exact divisor pairs b·r=numPerm in ascending b, integrates each interval with composite Simpson using 1000 intervals, and selects strictly smaller error. Ties retain the smaller b.

Theoretical candidate probabilities for all 64-permutation divisor pairs:

| Bands × rows | J=0.3 | J=0.5 | J=0.7 | J=0.85 | J=0.9 |
|---|---:|---:|---:|---:|---:|
| 1 × 64 | 0.000000 | 0.000000 | 0.000000 | 0.000030 | 0.001179 |
| 2 × 32 | 0.000000 | 0.000000 | 0.000022 | 0.010996 | 0.067495 |
| 4 × 16 | 0.000000 | 0.000061 | 0.013227 | 0.265532 | 0.559459 |
| 8 × 8 | 0.000525 | 0.030826 | 0.378122 | 0.921529 | 0.988930 |
| 16 × 4 | 0.122017 | 0.643926 | 0.987638 | 0.999993 | 1.000000 |
| 32 × 2 | 0.951098 | 0.999900 | 1.000000 | 1.000000 | 1.000000 |
| 64 × 1 | 1.000000 | 1.000000 | 1.000000 | 1.000000 | 1.000000 |

Values are rounded to six decimals; 0.000000 and 1.000000 need not be exact endpoints.

For 64 permutations, `Bands` chooses 16×4 at threshold 0.5, 8×8 at 0.7, and 4×16 at 0.8, 0.85, or 0.9. At 0.85 the balanced **4×16 catches only ~27%** of pairs exactly at J=0.85 (P≈0.265532). The bake-off's **16×4 catches ≈99.999%** at J=0.85, at the cost of ~12% false candidates at J=0.3 (P≈0.122017). Exact Jaccard verification filters those false candidates.

For candidate-then-verify dedup where a missed duplicate costs more than an extra Jaccard check, **use more bands, such as 16×4 for 64 permutations at threshold 0.85, rather than `Bands`**. Measure recall and candidate volume on the actual corpus. The fixed random test fixture measured 1000/1000 candidates at J=0.9 and 126/1000 at J=0.3 with seed 42, consistently through BandKeys and Index.

## Concurrency

Hasher is immutable after New and safe to share. Index guards its mutable state with a `sync.RWMutex`; do not copy an Index. Candidates are deduplicated in successful Add order; for simultaneous Adds this means write-lock acquisition order. Index retains IDs and bucket sequence numbers, not caller signature slices. Returned candidate slices belong to the caller.

Callers must not mutate inputs during calls or use the same writable SignInto destination concurrently. Independent destinations can share one Hasher. Index is intentionally small and in-memory, with no deletion or persistence API; use BandKeys for external database storage.
