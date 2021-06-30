package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

const API_USER string = "FKT-dv-jenkins"
const API_KEY_PATH string = "/dv/common/github.apikey"

type Config struct {
	aws_profile    string
	aws_parameter  string
	aws_key_id     string
	aws_key_secret string
	github_user    string
	github_apikey  string
	github_repo    string
	text           string
	pubkey         *PubKey
	decrypt        bool
}

func NewConfig() *Config {
	pkey := NewPubKey()

	return &Config{
		aws_profile:    "",
		aws_parameter:  "/dv/github-aws-credential-json",
		aws_key_id:     "",
		aws_key_secret: "",
		github_user:    API_USER,
		github_apikey:  "",
		github_repo:    "",
		text:           "",
		pubkey:         pkey,
		decrypt:        false,
	}
}

// TODO Maybe put the pubkey stuff in it's own file?
type PubKey struct {
	Key    string
	Key_id string
}

func NewPubKey() *PubKey {
	return &PubKey{
		Key:    "",
		Key_id: "",
	}
}

// fetchPublicKey gets the public key and key_id for a given Github repository
func (config *Config) FetchPublicKey() {
	urlString := "https://api.github.com/repos/otto-ec/" + config.github_repo + "/actions/secrets/public-key"
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(127)
	}
	req.SetBasicAuth(config.github_user, config.github_apikey)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(127)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		json.Unmarshal([]byte(bodyString), config.pubkey)
	} else {
		fmt.Println("Repository not found or accessible. Status code " + strconv.Itoa(resp.StatusCode))
		os.Exit(127)
	}
}

// fetchGithubAPIKey gets the credentials needed to access Github from the parameter store in AWS
func (config *Config) fetchGithubAPIKey() {
	path := API_KEY_PATH
	config.github_apikey = getSecret(&path)
	fmt.Println("Api: " + config.github_apikey)
}

// storeAWSCredentials stores the AWS credentials to the config
func (config *Config) storeAWSCredentials(keyId string, keySecret string) {
	config.aws_key_id = keyId
	config.aws_key_secret = keySecret
}

// parseCLIArgs parses all command line arguments
func (config *Config) parseCLIArgs() {
	// read parameters from the command line
	flag.StringVar(&config.aws_profile, "aws_profile", "", "AWS profile. (Required)")
	flag.StringVar(&config.github_user, "github_user", "", "Github user name. (Required)")
	flag.StringVar(&config.github_apikey, "github_apikey", "", "Github API key. (Required)")
	flag.StringVar(&config.github_repo, "github_repo", "", "Github repository where the secrets will be added. (Required)")
	flag.BoolVar(&config.decrypt, "decrypt", false, "Decrypt given cypher text. Default is to encrypt from parameter store data.")
	flag.StringVar(&config.text, "text", "", "Text to either encrypt or decrypt. (Optional)")
	flag.Parse()
}

// validate checks if all necessary arguments were given on the command line
func (config *Config) validate() error {
	// TODO If aws_profile given then it might be okay to fetch github_user and apikey from parameter store instead of cli
	if err := checkStringFlagNotEmpty("github_user", config.github_user); err != nil {
		return err
	}

	if err := checkStringFlagNotEmpty("github_apikey", config.github_apikey); err != nil {
		return err
	}

	if err := checkStringFlagNotEmpty("github_repo", config.github_repo); err != nil {
		return err
	}

	if err := checkStringFlagNotEmpty("aws_profile", config.aws_profile); err != nil {
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
