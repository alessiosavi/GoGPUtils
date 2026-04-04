// Package randutil provides random number generation utilities.
//
// The package provides two types of generators:
//
// 1. Secure generators using crypto/rand (for security-sensitive applications)
// 2. Fast generators using math/rand/v2 (for non-security applications)
//
// # Secure Random Generation
//
// For cryptographic purposes, password generation, or security tokens:
//
//	bytes, err := randutil.SecureBytes(32)
//	str, err := randutil.SecureString(16, randutil.AlphaNumeric)
//	id, err := randutil.SecureID()
//
// # Fast Random Generation
//
// For non-security purposes like tests, sampling, or simulations:
//
//	rng := randutil.NewGenerator()
//	n := rng.Int(100)             // [0, 100)
//	f := rng.Float64()            // [0.0, 1.0)
//	item := randutil.Choice(rng, items)
//	randutil.Shuffle(rng, items)  // Shuffle in place
//
// # Deterministic Generation
//
// For reproducible random sequences in tests:
//
//	rng := randutil.NewGeneratorWithSeed(42)
//	// Always produces the same sequence
//
// # Character Sets
//
// Predefined character sets for string generation:
//
//	Digits      - "0123456789"
//	Lowercase   - "abcdefghijklmnopqrstuvwxyz"
//	Uppercase   - "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
//	Letters     - Lowercase + Uppercase
//	AlphaNumeric - Letters + Digits
//	Hex         - "0123456789abcdef"
package randutil
