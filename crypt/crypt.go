package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/hex"
	"io"
)

func Encrypt(data []byte, passphrase string) (string, error) {
	h, err := HashMD5(passphrase)
	if err != nil {
		return "", err
	}
	block, _ := aes.NewCipher([]byte(h))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return b64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(data, passphrase string) (string, error) {
	h, err := HashMD5(passphrase)
	if err != nil {
		return "", err
	}
	key := []byte(h)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	raw, err := b64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := raw[:nonceSize], raw[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func HashMD5(key string) (string, error) {
	hasher := md5.New()
	if _, err := hasher.Write([]byte(key)); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
