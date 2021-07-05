package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var config *Config

func handleText(key []byte) {
	if config.decrypt == true {
		b, _ := b64.StdEncoding.DecodeString(config.text)
		decrypted, _ := Decrypt(b, []byte(key))
		fmt.Println(string(decrypted))
	} else {
		plainbytes := []byte(config.text)
		cyphercyper, _ := Encrypt(plainbytes, []byte(key))
		fmt.Println(b64.StdEncoding.EncodeToString(cyphercyper))
	}
}

func handleParameterstoreValues(key []byte) {
	var m map[string]string
	m = make(map[string]string)

	fmt.Println("Getting credentials from Parameter Store.")
	secretString := getSecret(&config.aws_parameter)
	json.Unmarshal([]byte(secretString), &m)
	config.storeAWSCredentials(m["AWS_ACCESS_KEY_ID"], m["AWS_ACCESS_KEY_SECRET"])

	for k, v := range m {
		plainbytes := []byte(v)
		cyphercyper, _ := Encrypt(plainbytes, []byte(key))
		fmt.Println("Updating Github secret " + k)
		if err := writeSecret(k, b64.StdEncoding.EncodeToString(cyphercyper)); err != nil {
			panic(err)
		}
	}
}

func main() {
	config = NewConfig()
	config.parseCLIArgs()

	if err := config.validate(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(1)
	}

	createAWSClient(&config.aws_profile, false)

	// get the default API key from parameter store if it wasn't given on cli
	if config.github_apikey == "" {
		config.fetchGithubAPIKey()
	}

	// public key fetched from github-repo
	config.FetchPublicKey()
	key, _ := b64.StdEncoding.DecodeString(config.pubkey.Key)

	// convert string from command line
	if config.text != "" {
		handleText(key)
	} else {
		handleParameterstoreValues(key)
	}
}
