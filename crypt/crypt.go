package crypt

import (
	"crypto/aes"
	"encoding/hex"
)

func EncryptAES(key []byte, plaintext string) string {
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	out := make([]byte, len(plaintext))
	c.Encrypt(out, []byte(plaintext))
	return hex.EncodeToString(out)
}

func DecryptAES(key []byte, ct string) string {
	ciphertext, _ := hex.DecodeString(ct)
	c, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	return s
}
