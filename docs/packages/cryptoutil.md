---
title: cryptoutil
parent: Packages
nav_order: 6
---

# cryptoutil

Secure cryptographic utilities for Go.
{: .fs-6 .fw-300 }

## Overview

The `cryptoutil` package provides secure, modern cryptographic operations with **zero external dependencies**. It uses **AES-GCM (Galois/Counter Mode)** for authenticated encryption, which provides both **confidentiality** and **integrity protection** in a single operation.

### Why AES-GCM?

AES-GCM is the industry standard for authenticated symmetric encryption. Unlike older modes like AES-ECB (which has serious weaknesses), AES-GCM:

- **Encrypts and authenticates** in one step — tampering is detected automatically
- **Uses random nonces** — identical plaintexts produce different ciphertexts
- **Is parallelizable** — high performance for large data
- **Is widely vetted** — recommended by NIST and used in TLS 1.3

### Key Features

| Feature               | Description                                               |
| --------------------- | --------------------------------------------------------- |
| AES-GCM encryption    | Authenticated encryption with confidentiality + integrity |
| Secure key generation | Cryptographically random keys via `crypto/rand`           |
| Key derivation        | SHA-256-based derivation from passwords                   |
| Hashing               | SHA-256 hashing with constant-time comparison             |
| Base64 encoding       | Ciphertext is base64-encoded for safe transport           |
| Zero dependencies     | Uses only Go standard library                             |

---

## Installation

```go
import "github.com/alessiosavi/GoGPUtils/cryptoutil"
```

---

## Constants and Errors

```go
var (
    ErrInvalidKeySize    = errors.New("key must be 16, 24, or 32 bytes")
    ErrInvalidCiphertext = errors.New("invalid ciphertext")
    ErrDecryptFailed     = errors.New("decryption failed: authentication error")
)
```

| Error                  | Description                                        |
| ---------------------- | -------------------------------------------------- |
| `ErrInvalidKeySize`    | Key is not 16, 24, or 32 bytes                     |
| `ErrInvalidCiphertext` | Ciphertext is malformed or invalid base64          |
| `ErrDecryptFailed`     | Authentication failed — wrong key or tampered data |

---

## Functions

### Encrypt

```go
func Encrypt(plaintext, key []byte) (string, error)
```

Encrypts plaintext using AES-GCM with the given key. The key must be 16, 24, or 32 bytes (AES-128, AES-192, AES-256). Returns base64-encoded ciphertext containing nonce + encrypted data + authentication tag.

**Parameters:**

- `plaintext` — The data to encrypt
- `key` — The encryption key (16, 24, or 32 bytes)

**Returns:**

- Base64-encoded ciphertext string
- Error if key size is invalid or encryption fails

**Example:**

```go
key := make([]byte, 32) // Generate a proper 256-bit key
_, err := io.ReadFull(rand.Reader, key)
if err != nil {
    log.Fatal(err)
}

ciphertext, err := cryptoutil.Encrypt([]byte("secret message"), key)
if err != nil {
    log.Fatal(err)
}
fmt.Println(ciphertext) // base64-encoded string
```

---

### Decrypt

```go
func Decrypt(ciphertextB64 string, key []byte) ([]byte, error)
```

Decrypts base64-encoded ciphertext using AES-GCM with the given key. Returns `ErrDecryptFailed` if authentication fails (indicating tampered data or wrong key).

**Parameters:**

- `ciphertextB64` — Base64-encoded ciphertext from `Encrypt`
- `key` — The encryption key (must match the key used for encryption)

**Returns:**

- Decrypted plaintext as byte slice
- Error if decryption or authentication fails

**Example:**

```go
plaintext, err := cryptoutil.Decrypt(ciphertext, key)
if err != nil {
    if errors.Is(err, cryptoutil.ErrDecryptFailed) {
        log.Fatal("Wrong key or tampered data")
    }
    log.Fatal(err)
}
fmt.Println(string(plaintext)) // "secret message"
```

---

### EncryptString

```go
func EncryptString(plaintext string, key []byte) (string, error)
```

Convenience function that encrypts a string and returns base64-encoded ciphertext. Equivalent to `Encrypt([]byte(plaintext), key)`.

**Example:**

```go
ciphertext, err := cryptoutil.EncryptString("Hello, World!", key)
if err != nil {
    log.Fatal(err)
}
```

---

### DecryptString

```go
func DecryptString(ciphertextB64 string, key []byte) (string, error)
```

Convenience function that decrypts base64-encoded ciphertext and returns the plaintext string. Equivalent to `string(Decrypt(ciphertextB64, key))`.

**Example:**

```go
plaintext, err := cryptoutil.DecryptString(ciphertext, key)
if err != nil {
    log.Fatal(err)
}
fmt.Println(plaintext) // "Hello, World!"
```

---

### GenerateKey

```go
func GenerateKey(size int) ([]byte, error)
```

Generates a cryptographically secure random key of the specified size using `crypto/rand`. Size must be 16, 24, or 32 bytes.

**Parameters:**

- `size` — Key size in bytes (16 for AES-128, 24 for AES-192, 32 for AES-256)

**Returns:**

- Random key as byte slice
- Error if size is invalid or random generation fails

**Example:**

```go
key, err := cryptoutil.GenerateKey(32) // AES-256
if err != nil {
    log.Fatal(err)
}
// key is 32 random bytes — store this securely!
```

---

### DeriveKey

```go
func DeriveKey(password, salt string) []byte
```

Derives a 32-byte key from a password and salt using SHA-256. This is a simple key derivation suitable for non-security-critical applications. For password storage or high-security applications, use `argon2` or `bcrypt` instead.

**Parameters:**

- `password` — The password string
- `salt` — A unique salt string

**Returns:**

- 32-byte key derived from password + salt

**Example:**

```go
key := cryptoutil.DeriveKey("my-password", "my-salt")
ciphertext, err := cryptoutil.Encrypt(data, key)
```

---

### GenerateNonce

```go
func GenerateNonce(size int) ([]byte, error)
```

Generates a random nonce of the specified size using `crypto/rand`. For AES-GCM, the standard nonce size is 12 bytes.

**Parameters:**

- `size` — Nonce size in bytes (must be positive)

**Returns:**

- Random nonce as byte slice
- Error if size is not positive or random generation fails

**Example:**

```go
nonce, err := cryptoutil.GenerateNonce(12) // Standard AES-GCM nonce size
if err != nil {
    log.Fatal(err)
}
```

---

### Hash

```go
func Hash(data []byte) []byte
```

Returns the SHA-256 hash of data as a 32-byte slice.

**Example:**

```go
hash := cryptoutil.Hash([]byte("Hello, World!"))
// hash is 32 bytes
```

---

### HashString

```go
func HashString(s string) string
```

Returns the SHA-256 hash of a string as a 64-character hex-encoded string.

**Example:**

```go
hash := cryptoutil.HashString("Hello, World!")
// hash is a 64-char lowercase hex string
```

---

### CompareHash

```go
func CompareHash(data, expectedHash []byte) bool
```

Compares a hash with a computed hash of data in **constant time**. This prevents timing attacks by ensuring the comparison takes the same amount of time regardless of how many bytes match.

**Parameters:**

- `data` — The original data
- `expectedHash` — The expected hash to compare against

**Returns:**

- `true` if hashes match, `false` otherwise

**Example:**

```go
hash := cryptoutil.Hash([]byte("secret"))
if cryptoutil.CompareHash([]byte("secret"), hash) {
    fmt.Println("Match!")
}
```

---

### RandomBytes

```go
func RandomBytes(n int) ([]byte, error)
```

Returns `n` cryptographically secure random bytes using `crypto/rand`.

**Parameters:**

- `n` — Number of bytes to generate (must be non-negative)

**Returns:**

- Random bytes as slice
- Error if `n` is negative or random generation fails

**Example:**

```go
bytes, err := cryptoutil.RandomBytes(32)
if err != nil {
    log.Fatal(err)
}
```

---

## Usage Examples

### Basic Encryption and Decryption

```go
package main

import (
    "fmt"
    "log"

    "github.com/alessiosavi/GoGPUtils/cryptoutil"
)

func main() {
    // Generate a secure 256-bit key
    key, err := cryptoutil.GenerateKey(32)
    if err != nil {
        log.Fatal(err)
    }

    // Encrypt a message
    plaintext := "Hello, secure world!"
    ciphertext, err := cryptoutil.EncryptString(plaintext, key)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Encrypted:", ciphertext)

    // Decrypt the message
    decrypted, err := cryptoutil.DecryptString(ciphertext, key)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Decrypted:", decrypted)
}
```

### Password-Based Encryption

```go
package main

import (
    "fmt"
    "log"

    "github.com/alessiosavi/GoGPUtils/cryptoutil"
)

func main() {
    // Derive a key from password and salt
    key := cryptoutil.DeriveKey("my-secure-password", "unique-salt-123")

    // Encrypt
    ciphertext, err := cryptoutil.EncryptString("Sensitive data", key)
    if err != nil {
        log.Fatal(err)
    }

    // Decrypt (same password and salt required)
    decrypted, err := cryptoutil.DecryptString(ciphertext, key)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(decrypted)
}
```

### Working with Byte Slices

```go
package main

import (
    "fmt"
    "log"

    "github.com/alessiosavi/GoGPUtils/cryptoutil"
)

func main() {
    key, _ := cryptoutil.GenerateKey(32)

    // Encrypt raw bytes
    data := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
    ciphertext, err := cryptoutil.Encrypt(data, key)
    if err != nil {
        log.Fatal(err)
    }

    // Decrypt back to bytes
    decrypted, err := cryptoutil.Decrypt(ciphertext, key)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%x\n", decrypted) // 0102030405
}
```

### Hashing and Verification

```go
package main

import (
    "fmt"

    "github.com/alessiosavi/GoGPUtils/cryptoutil"
)

func main() {
    data := []byte("important data")

    // Compute hash
    hash := cryptoutil.Hash(data)
    fmt.Printf("Hash: %x\n", hash)

    // Verify in constant time (safe from timing attacks)
    if cryptoutil.CompareHash(data, hash) {
        fmt.Println("Data integrity verified")
    }
}
```

---

## Security Best Practices

### Key Management

This package does **not** handle key management. It is the caller's responsibility to:

1. **Generate keys securely** — Always use `GenerateKey()` or `crypto/rand` instead of `math/rand`
2. **Store keys securely** — Use secret management systems (AWS Secrets Manager, HashiCorp Vault, environment variables)
3. **Never hardcode keys** — Keys in source code are exposed in version control
4. **Rotate keys periodically** — Regular key rotation limits exposure
5. **Use 32-byte keys** — Prefer AES-256 (32 bytes) over AES-128 (16 bytes) for sensitive data

### Nonce Handling

- Nonces are **generated automatically** by `Encrypt()` using `crypto/rand`
- **Never reuse a key with the same nonce** — this breaks AES-GCM security
- The nonce is **prepended to the ciphertext** and extracted during decryption

### What Not to Use This For

| Use Case              | Recommended Alternative                |
| --------------------- | -------------------------------------- |
| Password storage      | `bcrypt`, `argon2`, or `scrypt`        |
| Digital signatures    | RSA or ECDSA                           |
| Key exchange          | ECDH or RSA key exchange               |
| Large file encryption | Stream encryption with chunked AES-GCM |

### Important Warnings

{: .warning }

> **Do not use `DeriveKey` for password hashing.** `DeriveKey` uses a single SHA-256 iteration which is not suitable for password storage. Use `bcrypt`, `argon2`, or `scrypt` for password hashing.

{: .warning }

> **Keep keys secret.** If an attacker obtains the key, they can decrypt all data encrypted with that key. Treat keys as sensitive as the data they protect.

{: .warning }

> **Verify errors.** Always check for `ErrDecryptFailed` — it indicates either the wrong key was used or the ciphertext was tampered with.

---

## Testing

Run the cryptoutil tests:

```bash
go test ./cryptoutil/...
```

Run with race detector:

```bash
go test -race ./cryptoutil/...
```

Run benchmarks:

```bash
go test -bench=. ./cryptoutil/...
```

### Test Coverage

The test suite covers:

- Round-trip encryption/decryption
- String encryption/decryption
- Different outputs for same plaintext (nonce uniqueness)
- Wrong key detection
- Tampered ciphertext detection
- Invalid key sizes
- Key generation
- Key derivation
- Hashing and comparison
- Random byte generation
- Empty and large plaintext handling
