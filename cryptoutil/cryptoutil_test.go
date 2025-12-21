package cryptoutil

import (
	"bytes"
	"errors"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key, err := GenerateKey(32)
	if err != nil {
		t.Fatalf("GenerateKey() error: %v", err)
	}

	plaintext := []byte("Hello, World! This is a secret message.")

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() error: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt() error: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Errorf("Decrypted = %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptDecryptString(t *testing.T) {
	key, _ := GenerateKey(32)
	plaintext := "Secret message"

	ciphertext, err := EncryptString(plaintext, key)
	if err != nil {
		t.Fatalf("EncryptString() error: %v", err)
	}

	decrypted, err := DecryptString(ciphertext, key)
	if err != nil {
		t.Fatalf("DecryptString() error: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("DecryptString() = %q, want %q", decrypted, plaintext)
	}
}

func TestEncryptDifferentOutputs(t *testing.T) {
	key, _ := GenerateKey(32)
	plaintext := []byte("Same message")

	// Encrypt same message twice
	cipher1, _ := Encrypt(plaintext, key)
	cipher2, _ := Encrypt(plaintext, key)

	// Should produce different ciphertexts due to random nonce
	if cipher1 == cipher2 {
		t.Error("Same plaintext should produce different ciphertexts")
	}

	// But both should decrypt to the same message
	plain1, _ := Decrypt(cipher1, key)
	plain2, _ := Decrypt(cipher2, key)

	if !bytes.Equal(plain1, plain2) || !bytes.Equal(plain1, plaintext) {
		t.Error("Both ciphertexts should decrypt to the same plaintext")
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key1, _ := GenerateKey(32)
	key2, _ := GenerateKey(32)

	ciphertext, _ := Encrypt([]byte("secret"), key1)

	_, err := Decrypt(ciphertext, key2)
	if !errors.Is(err, ErrDecryptFailed) {
		t.Errorf("Decrypt with wrong key error = %v, want ErrDecryptFailed", err)
	}
}

func TestDecryptTampered(t *testing.T) {
	key, _ := GenerateKey(32)
	ciphertext, _ := Encrypt([]byte("secret"), key)

	// Tamper with the ciphertext
	tampered := ciphertext[:len(ciphertext)-1] + "X"

	_, err := Decrypt(tampered, key)
	if err == nil {
		t.Error("Decrypt should fail for tampered ciphertext")
	}
}

func TestInvalidKeySize(t *testing.T) {
	invalidKeys := []int{0, 1, 15, 17, 23, 25, 31, 33, 64}

	for _, size := range invalidKeys {
		key := make([]byte, size)

		_, err := Encrypt([]byte("test"), key)
		if !errors.Is(err, ErrInvalidKeySize) {
			t.Errorf("Encrypt with %d-byte key error = %v, want ErrInvalidKeySize", size, err)
		}
	}

	validKeys := []int{16, 24, 32}
	for _, size := range validKeys {
		key := make([]byte, size)

		_, err := Encrypt([]byte("test"), key)
		if err != nil {
			t.Errorf("Encrypt with %d-byte key should succeed, got error: %v", size, err)
		}
	}
}

func TestGenerateKey(t *testing.T) {
	sizes := []int{16, 24, 32}

	for _, size := range sizes {
		key, err := GenerateKey(size)
		if err != nil {
			t.Errorf("GenerateKey(%d) error: %v", size, err)
		}

		if len(key) != size {
			t.Errorf("GenerateKey(%d) returned %d bytes", size, len(key))
		}
	}

	// Invalid sizes
	_, err := GenerateKey(17)
	if !errors.Is(err, ErrInvalidKeySize) {
		t.Errorf("GenerateKey(17) error = %v, want ErrInvalidKeySize", err)
	}
}

func TestDeriveKey(t *testing.T) {
	key := DeriveKey("password", "salt")

	if len(key) != 32 {
		t.Errorf("DeriveKey() returned %d bytes, want 32", len(key))
	}

	// Same input should produce same output
	key2 := DeriveKey("password", "salt")
	if !bytes.Equal(key, key2) {
		t.Error("DeriveKey() should be deterministic")
	}

	// Different input should produce different output
	key3 := DeriveKey("password", "different-salt")
	if bytes.Equal(key, key3) {
		t.Error("Different salt should produce different key")
	}

	// Key should work for encryption
	ciphertext, err := EncryptString("test", key)
	if err != nil {
		t.Errorf("Encrypt with derived key error: %v", err)
	}

	plaintext, err := DecryptString(ciphertext, key)
	if err != nil || plaintext != "test" {
		t.Error("Decryption with derived key failed")
	}
}

func TestHash(t *testing.T) {
	data := []byte("Hello, World!")
	hash1 := Hash(data)
	hash2 := Hash(data)

	if len(hash1) != 32 {
		t.Errorf("Hash() returned %d bytes, want 32", len(hash1))
	}

	if !bytes.Equal(hash1, hash2) {
		t.Error("Hash() should be deterministic")
	}

	// Different data should produce different hash
	hash3 := Hash([]byte("Different data"))
	if bytes.Equal(hash1, hash3) {
		t.Error("Different data should produce different hash")
	}
}

func TestHashString(t *testing.T) {
	hash := HashString("Hello, World!")

	if len(hash) != 64 { // 32 bytes * 2 hex chars
		t.Errorf("HashString() returned %d chars, want 64", len(hash))
	}

	// Should be lowercase hex
	for _, c := range hash {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
			t.Errorf("HashString() contains invalid hex char: %c", c)
		}
	}
}

func TestCompareHash(t *testing.T) {
	data := []byte("test data")
	hash := Hash(data)

	if !CompareHash(data, hash) {
		t.Error("CompareHash should return true for matching data")
	}

	if CompareHash([]byte("different data"), hash) {
		t.Error("CompareHash should return false for different data")
	}

	if CompareHash(data, []byte("wrong hash")) {
		t.Error("CompareHash should return false for wrong hash")
	}
}

func TestRandomBytes(t *testing.T) {
	sizes := []int{0, 1, 16, 32, 100}

	for _, size := range sizes {
		data, err := RandomBytes(size)
		if err != nil {
			t.Errorf("RandomBytes(%d) error: %v", size, err)
		}

		if len(data) != size {
			t.Errorf("RandomBytes(%d) returned %d bytes", size, len(data))
		}
	}

	// Two calls should produce different output
	data1, _ := RandomBytes(32)

	data2, _ := RandomBytes(32)

	if bytes.Equal(data1, data2) {
		t.Error("RandomBytes should produce different output on each call")
	}

	// Negative size should error
	_, err := RandomBytes(-1)
	if err == nil {
		t.Error("RandomBytes(-1) should error")
	}
}

func TestGenerateNonce(t *testing.T) {
	nonce, err := GenerateNonce(12)
	if err != nil {
		t.Fatalf("GenerateNonce() error: %v", err)
	}

	if len(nonce) != 12 {
		t.Errorf("GenerateNonce(12) returned %d bytes", len(nonce))
	}

	_, err = GenerateNonce(0)
	if err == nil {
		t.Error("GenerateNonce(0) should error")
	}
}

func TestEmptyPlaintext(t *testing.T) {
	key, _ := GenerateKey(32)

	ciphertext, err := Encrypt([]byte{}, key)
	if err != nil {
		t.Fatalf("Encrypt empty error: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}

	if len(decrypted) != 0 {
		t.Errorf("Decrypted empty = %q, want empty", decrypted)
	}
}

func TestLargePlaintext(t *testing.T) {
	key, _ := GenerateKey(32)

	// 1MB of data
	plaintext := make([]byte, 1024*1024)
	for i := range plaintext {
		plaintext[i] = byte(i % 256)
	}

	ciphertext, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt large error: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("Decrypt large error: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Error("Large plaintext mismatch")
	}
}

// ============================================================================
// Benchmarks
// ============================================================================

func BenchmarkEncrypt(b *testing.B) {
	key, _ := GenerateKey(32)
	plaintext := make([]byte, 1024)

	for b.Loop() {
		Encrypt(plaintext, key)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	key, _ := GenerateKey(32)
	plaintext := make([]byte, 1024)
	ciphertext, _ := Encrypt(plaintext, key)

	for b.Loop() {
		Decrypt(ciphertext, key)
	}
}

func BenchmarkHash(b *testing.B) {
	data := make([]byte, 1024)

	for b.Loop() {
		Hash(data)
	}
}

func BenchmarkGenerateKey(b *testing.B) {
	for b.Loop() {
		GenerateKey(32)
	}
}
