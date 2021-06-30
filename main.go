package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	var m map[string]string
	m = make(map[string]string)

	config := NewConfig()
	config.parseCLIArgs()

	if err := config.validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(1)
	}

	createAWSClient(&config.aws_profile, false)

	secretString := getSecret(&config.aws_parameter)
	json.Unmarshal([]byte(secretString), &m)
	config.storeAWSCredentials(m["AWS_ACCESS_KEY_ID"], m["AWS_ACCESS_KEY_SECRET"])

	// get the default API key from parameter store if it wasn't given on cli
	if config.github_apikey == "" {
		config.fetchGithubAPIKey()
	}

	// convert string from command line
	plaintext := []byte(config.text)

	// public key fetched from github-repo
	config.FetchPublicKey()

	key, _ := b64.StdEncoding.DecodeString(config.pubkey.Key)

	cyphercyper, _ := Encrypt(plaintext, []byte(key))
	fmt.Println(b64.StdEncoding.EncodeToString(cyphercyper))

	decrypted, _ := Decrypt(cyphercyper, []byte(key))
	fmt.Println(string(decrypted))
}
