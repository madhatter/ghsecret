package main

import (
	"crypto/rand"
	b64 "encoding/base64"
	"flag"
	"fmt"
	"os"

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

func main() {
	config := NewConfig()
	config.parseCLIArgs()

	if err := config.validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(1)
	}

	// convert string from command line
	plaintext := []byte(config.text)

	// public key fetched from github-repo
	config.pubkey.FetchPublicKey(config.github_repo)

	//pkey := "+YM3PNgG3jET4XyWRuxpc8p2frjgI0D/OULKqNZ2cBM="
	//key, _ := b64.StdEncoding.DecodeString(pkey)

	key, _ := b64.StdEncoding.DecodeString(config.pubkey.key)

	cyphercyper, _ := Encrypt(plaintext, []byte(key))
	fmt.Println(b64.StdEncoding.EncodeToString(cyphercyper))

	decrypted, _ := Decrypt(cyphercyper, []byte(key))
	fmt.Println(string(decrypted))
}
