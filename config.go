package main

import (
	"flag"
	"fmt"
)

type Config struct {
	github_user   string
	github_apikey string
	github_repo   string
	text          string
	pubkey        *PubKey
	decrypt       bool
}

func NewConfig() *Config {
	pkey := NewPubKey()

	return &Config{
		github_user:   "",
		github_apikey: "",
		github_repo:   "",
		text:          "",
		pubkey:        pkey,
		decrypt:       false,
	}
}

// TODO Maybe put the pubkey stuff in it's own file?
type PubKey struct {
	key_id string
	key    string
}

func NewPubKey() *PubKey {
	return &PubKey{
		key_id: "",
		key:    "",
	}
}

// fetchPublicKey gets the public key and key_id for a given Github repository
func (pubkey *PubKey) FetchPublicKey(repository string) {
	// TODO This has to fetch the key for real
	pubkey.key_id = "568250167242549743"
	pubkey.key = "+YM3PNgG3jET4XyWRuxpc8p2frjgI0D/OULKqNZ2cBM="
}

func (config *Config) parseCLIArgs() {
	// read parameters from the command line
	flag.StringVar(&config.github_user, "github_user", "", "Github user name. (Required)")
	flag.StringVar(&config.github_apikey, "github_apikey", "", "Github API key. (Required)")
	flag.StringVar(&config.github_repo, "github_repo", "", "Github repository where the secrets will be added. (Required)")
	flag.BoolVar(&config.decrypt, "decrypt", false, "Decrypt given cypher text. Default is to encrypt from parameter store data.")
	flag.StringVar(&config.text, "text", "", "Text to either encrypt or decrypt. (Required)")
	flag.Parse()
}

func (config *Config) validate() error {
	if err := checkStringFlagNotEmpty("github_user", config.github_user); err != nil {
		return err
	}

	if err := checkStringFlagNotEmpty("github_apikey", config.github_apikey); err != nil {
		return err
	}

	if err := checkStringFlagNotEmpty("github_repo", config.github_repo); err != nil {
		return err
	}

	if err := checkStringFlagNotEmpty("text", config.text); err != nil {
		return err
	}

	return nil
}

func checkStringFlagNotEmpty(name string, flag string) error {
	if flag == "" {
		return fmt.Errorf("Missing mandatory parameter: %s", name)
	}
	return nil
}
