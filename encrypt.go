package main

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
)

// AES-GCM should be used because the operation is an authenticated encryption
// algorithm designed to provide both data authenticity (integrity) as well as
// confidentiality.

// Merged into Golang in https://go-review.googlesource.com/#/c/18803/

func encrypt(plaintext []byte, key string) (string, error) {
	nonce := []byte("eequohn8aeCh") // key is random each encrypt, dont really need an extra nonce each time
	bkey := []byte(key)

	block, err := aes.NewCipher(bkey)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

func decrypt(text string, key string) (string, error) {
	nonce := []byte("eequohn8aeCh") // key is random each encrypt, dont really need an extra nonce each time
	bkey := []byte(key)

	ciphertext, _ := hex.DecodeString(text)

	block, err := aes.NewCipher(bkey)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
