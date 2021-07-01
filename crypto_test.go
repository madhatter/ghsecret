package main

import (
	b64 "encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypto_DecryptedMatchesEncryptedText(t *testing.T) {
	plaintext := "Fancy Secret"
	plainbytes := []byte(plaintext)
	//pubkey := "+YM3PNgG3jET4XyWRuxpc8p2frjgI0D/OULKqNZ2cBM="
	pubkey := "Ju5t4R4nd0mK3yTh4t5L0Ng3N0uGh2W0rKF0rS0D1um="
	key, _ := b64.StdEncoding.DecodeString(pubkey)

	crypted, _ := Encrypt(plainbytes, key)
	decrypted, _ := Decrypt(crypted, key)

	assert.Equal(t, plainbytes, decrypted)
}
