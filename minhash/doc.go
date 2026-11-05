// Package minhash provides deterministic word shingling, MinHash signatures,
// LSH band keys for persistence, and a small concurrent in-memory index.
// It uses only the standard library. Hashes are noncryptographic; collisions
// are possible. Verify LSH candidates with Jaccard on their shingle sets.
//
// # Tokens and usage
//
// No normalization is performed: tokens are raw bytes, so NFC "è" differs from
// NFD "e\u0300". Normalize first, for example with the separate textnorm package:
//
//	tokens, err := textnorm.MeaningPreset().SplitTokens().Run(text)
//	if err != nil { return err }
//	shingles, err := minhash.Shingles(tokens, 5)
//	if err != nil { return err }
//	h, err := minhash.New(64, 42)
//	if err != nil { return err }
//	sig, err := h.Sign(shingles)
//	if err != nil { return err } // empty text returns ErrEmptySet
//	keys, err := minhash.BandKeys(sig, 16)
//
// Empty tokens input produces no shingles; a nonempty input shorter than k
// produces one shingle over all tokens. Empty-string tokens are legal.
//
// # Persisted format version 1
//
// All integer arithmetic below is unsigned, wrapping modulo 2^64 except the
// explicitly stated modulo-p arithmetic. A shingle encodes each token in order
// as binary.AppendUvarint(uint64(len(token))) followed by its raw bytes. Lengths
// count bytes, not runes. The incremental byte hash is FNV-1a 64: start at
// 14695981039346656037, then for each byte h = (h XOR byte) * 1099511628211.
// The shingle hash is mix64(h), the splitmix64 finalizer:
//
//	z = (z XOR (z >> 30)) * 0xbf58476d1ce4e5b9
//	z = (z XOR (z >> 27)) * 0x94d049bb133111eb
//	return z XOR (z >> 31)
//
// Permutations use p = 2^61-1. Initialize state = seed. Each next() increments
// state by 0x9e3779b97f4a7c15 and returns mix64(state). A draw is next() >> 3.
// For each permutation in order, draw a until 1 <= a < p, then b until 0 <= b < p.
// For every shingle reduce x modulo p, then compute (a*x+b) modulo p with full
// precision. Each signature position is the minimum over the set. A larger
// hasher with the same seed shares the smaller hasher's permutation prefix.
// Values are below p and fit signed int64 / SQL BIGINT unchanged.
//
// A band j with r rows hashes uint64LE(j) followed by the r signature values
// in order, each uint64LE. Its key is mix64(FNV-1a-64(encoded bytes)). All 64
// bits are significant. Persist keys in signed BIGINT as int64(key); recover
// with uint64(stored). This preserves bits and equality even for negative keys.
//
// Store (signature_version, band, key) rows, indexed by (band, key) or
// (signature_version, band, key). Find candidates by equal (band, key), scoped
// to the same version, seed, shingle k, numPerm and band layout. Persist that
// configuration in a namespace or metadata; version alone does not identify it.
// Normalization must also remain consistent. Any change to token encoding,
// shingle hashing, permutation derivation or formula, or band keys MUST bump
// SignatureVersion and add golden vectors for that version, retaining v1 vectors.
// Never compare incompatible persisted formats without recomputing from text.
//
// # Choosing bands
//
// The ideal LSH S-curve is P(s) = 1-(1-s^r)^b, the probability of at least one
// matching band at Jaccard s. It is a probabilistic model, not a hard threshold.
// Bands minimizes 0.5*FP + 0.5*FN, where FP is the integral of P(s) from 0 to
// threshold and FN is the integral of 1-P(s) from threshold to 1. It enumerates
// exact divisor pairs in ascending b and uses composite Simpson integration
// with 1000 intervals per integral; strict improvement wins, preserving the
// smaller b on ties.
//
// For 64 permutations and threshold 0.85, the balanced choice is 4x16, catching
// only about 27% of pairs at exactly J=0.85. The bake-off's 16x4 catches about
// 99.999% at J=0.85 but also about 12% at J=0.3. Exact Jaccard verification
// filters those false candidates. For candidate-then-verify dedup where a missed
// duplicate costs more than another Jaccard check, use more bands (e.g. 16x4
// for 64 permutations at 0.85) instead of the balanced Bands recommendation.
//
// # Concurrency
//
// Hasher is immutable and safe to share. Index uses a sync.RWMutex and returns
// candidates in successful insertion order (lock acquisition order for racing
// Adds). Do not copy an Index. Callers own their slices: do not mutate inputs
// during calls or share a writable SignInto destination across goroutines.
package minhash
