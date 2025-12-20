// Package cryptoutil provides secure cryptographic operations.
//
// This package uses only secure, modern cryptographic primitives.
// Specifically, it uses AES-GCM for authenticated encryption,
// which provides both confidentiality and integrity protection.
//
// # Key Points
//
//   - Uses AES-GCM (Galois/Counter Mode) - authenticated encryption
//   - Keys must be 16, 24, or 32 bytes (AES-128, AES-192, AES-256)
//   - Nonces are generated automatically using crypto/rand
//   - Output is base64-encoded for safe transport
//
// # Why AES-GCM?
//
// The original library used AES-ECB mode which has serious weaknesses:
//   - Identical plaintext blocks produce identical ciphertext
//   - No authentication (ciphertext can be modified undetected)
//
// AES-GCM solves both problems:
//   - Each message is unique due to random nonce
//   - Authentication tag prevents tampering
//
// # Usage
//
//	key := make([]byte, 32) // 256-bit key
//	rand.Read(key)
//
//	ciphertext, err := cryptoutil.Encrypt(plaintext, key)
//	if err != nil { ... }
//
//	plaintext, err := cryptoutil.Decrypt(ciphertext, key)
//	if err != nil { ... }
//
// # Key Management
//
// This package does not handle key management. It is the caller's
// responsibility to:
//   - Generate keys securely (use crypto/rand)
//   - Store keys securely (use secret management systems)
//   - Rotate keys periodically
//
// # Security Considerations
//
//   - Never reuse a key with the same nonce
//   - Use at least 128-bit keys (32 bytes for AES-256)
//   - Do not use this for password hashing (use bcrypt/argon2)
package cryptoutil
