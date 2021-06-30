package main

import (
	"crypto/rand"

	"github.com/jamesruan/sodium"
)

const NONCE_LEN int = 24

func randomNonce() ([]byte, error) {
	b := make([]byte, NONCE_LEN)
	_, err := rand.Read(b)
	return b, err
}

func Encrypt(plaintext, k []byte) ([]byte, error) {
	nonce, err := randomNonce()
	if err != nil {
		return nil, err
	}
	cyphertext := sodium.Bytes(plaintext).SecretBox(
		sodium.SecretBoxNonce{nonce},
		sodium.SecretBoxKey{k})
	return append(nonce, cyphertext...), nil
}

func Decrypt(cyphertext, k []byte) ([]byte, error) {
	nonce := sodium.SecretBoxNonce{cyphertext[:NONCE_LEN]}
	enc := sodium.Bytes(cyphertext[NONCE_LEN:])
	return enc.SecretBoxOpen(nonce, sodium.SecretBoxKey{k})
}
