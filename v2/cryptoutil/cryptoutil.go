package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// Common errors.
var (
	ErrInvalidKeySize    = errors.New("key must be 16, 24, or 32 bytes")
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
	ErrDecryptFailed     = errors.New("decryption failed: authentication error")
)

// Encrypt encrypts plaintext using AES-GCM with the given key.
// The key must be 16, 24, or 32 bytes (AES-128, AES-192, AES-256).
// Returns base64-encoded ciphertext containing nonce + encrypted data + auth tag.
//
// Example:
//
//	key := make([]byte, 32) // Generate a proper key
//	io.ReadFull(rand.Reader, key)
//	ciphertext, err := Encrypt([]byte("secret message"), key)
func Encrypt(plaintext, key []byte) (string, error) {
	if !isValidKeySize(len(key)) {
		return "", ErrInvalidKeySize
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and prepend nonce
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts base64-encoded ciphertext using AES-GCM with the given key.
// Returns ErrDecryptFailed if authentication fails (tampered or wrong key).
//
// Example:
//
//	plaintext, err := Decrypt(ciphertext, key)
//	if errors.Is(err, ErrDecryptFailed) {
//	    // Wrong key or tampered data
//	}
func Decrypt(ciphertextB64 string, key []byte) ([]byte, error) {
	if !isValidKeySize(len(key)) {
		return nil, ErrInvalidKeySize
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return nil, ErrInvalidCiphertext
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, ErrInvalidCiphertext
	}

	nonce, encryptedData := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return nil, ErrDecryptFailed
	}

	return plaintext, nil
}

// EncryptString encrypts a string and returns base64-encoded ciphertext.
func EncryptString(plaintext string, key []byte) (string, error) {
	return Encrypt([]byte(plaintext), key)
}

// DecryptString decrypts base64-encoded ciphertext and returns the plaintext string.
func DecryptString(ciphertextB64 string, key []byte) (string, error) {
	plaintext, err := Decrypt(ciphertextB64, key)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// DeriveKey derives a key from a password using SHA-256.
// This is a simple key derivation for non-security-critical applications.
// For password storage or high-security applications, use argon2 or bcrypt.
//
// Example:
//
//	key := DeriveKey("mypassword", "mysalt")
//	ciphertext, err := Encrypt(data, key)
func DeriveKey(password, salt string) []byte {
	h := sha256.New()
	h.Write([]byte(password))
	h.Write([]byte(salt))
	return h.Sum(nil) // Returns 32 bytes (AES-256 key)
}

// GenerateKey generates a cryptographically secure random key of the specified size.
// Size must be 16, 24, or 32 bytes.
//
// Example:
//
//	key, err := GenerateKey(32) // AES-256
//	if err != nil { ... }
func GenerateKey(size int) ([]byte, error) {
	if !isValidKeySize(size) {
		return nil, ErrInvalidKeySize
	}

	key := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// GenerateNonce generates a random nonce of the specified size.
// For AES-GCM, the standard nonce size is 12 bytes.
func GenerateNonce(size int) ([]byte, error) {
	if size <= 0 {
		return nil, errors.New("nonce size must be positive")
	}

	nonce := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return nonce, nil
}

// Hash returns the SHA-256 hash of data as a byte slice.
func Hash(data []byte) []byte {
	h := sha256.Sum256(data)
	return h[:]
}

// HashString returns the SHA-256 hash of a string as a hex-encoded string.
func HashString(s string) string {
	h := sha256.Sum256([]byte(s))
	return encodeHex(h[:])
}

// CompareHash compares a hash with a computed hash of data in constant time.
// This prevents timing attacks.
func CompareHash(data, expectedHash []byte) bool {
	computed := Hash(data)
	if len(computed) != len(expectedHash) {
		return false
	}

	// Constant-time comparison
	result := byte(0)
	for i := range computed {
		result |= computed[i] ^ expectedHash[i]
	}
	return result == 0
}

// RandomBytes returns n cryptographically secure random bytes.
func RandomBytes(n int) ([]byte, error) {
	if n < 0 {
		return nil, errors.New("size must be non-negative")
	}
	if n == 0 {
		return []byte{}, nil
	}

	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return nil, err
	}
	return b, nil
}

// isValidKeySize checks if the key size is valid for AES.
func isValidKeySize(size int) bool {
	return size == 16 || size == 24 || size == 32
}

// encodeHex encodes bytes to hexadecimal string.
func encodeHex(data []byte) string {
	const hexChars = "0123456789abcdef"
	buf := make([]byte, len(data)*2)
	for i, b := range data {
		buf[i*2] = hexChars[b>>4]
		buf[i*2+1] = hexChars[b&0x0f]
	}
	return string(buf)
}

// decodeHex decodes a hexadecimal string to bytes.
func decodeHex(s string) ([]byte, error) {
	if len(s)%2 != 0 {
		return nil, errors.New("hex string must have even length")
	}

	result := make([]byte, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		high, err := hexCharToByte(s[i])
		if err != nil {
			return nil, err
		}
		low, err := hexCharToByte(s[i+1])
		if err != nil {
			return nil, err
		}
		result[i/2] = (high << 4) | low
	}
	return result, nil
}

func hexCharToByte(c byte) (byte, error) {
	switch {
	case c >= '0' && c <= '9':
		return c - '0', nil
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10, nil
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10, nil
	default:
		return 0, errors.New("invalid hex character")
	}
}
