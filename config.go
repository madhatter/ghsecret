package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const API_USER string = "FKT-dv-jenkins"
const API_KEY_PATH string = "/dv/common/github.apikey"

func init() {
	log.SetOutput(os.Stdout)
}

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
	RandomOverride io.Reader
	debug          bool
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
		debug:          false,
	}
}

// TODO Maybe put the pubkey stuff in it's own file?
type PubKey struct {
	Raw    [32]byte
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
func (config *Config) fetchPublicKey() {
	urlString := "https://api.github.com/repos/otto-ec/" + config.github_repo +
		"/actions/secrets/public-key"
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		log.Errorln(err)
		os.Exit(127)
	}
	req.SetBasicAuth(config.github_user, config.github_apikey)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorln(err)
		os.Exit(127)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, io_err := io.ReadAll(resp.Body)
		if io_err != nil {
			log.Fatal(io_err)
		}
		bodyString := string(bodyBytes)
		json.Unmarshal([]byte(bodyString), config.pubkey)
		log.Debugln("Public key: " + config.pubkey.Key)
	} else {
		log.Errorln("Repository not found or accessible. Status code " +
			strconv.Itoa(resp.StatusCode))
		os.Exit(127)
	}

	decoded, err := base64.StdEncoding.DecodeString(config.pubkey.Key)
	if err != nil {
		log.Error("failed to decode public key: %w", err)
	}

	copy(config.pubkey.Raw[:], decoded[0:32])
}

// fetchGithubAPIKey gets the credentials needed to access Github from the
// parameter store in AWS
func (config *Config) fetchGithubAPIKey() {
	path := API_KEY_PATH
	config.github_apikey = getSecret(&path)
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
	flag.StringVar(&config.github_user, "github_user", "", "Github user name. Default: "+API_USER+" (Optional)")
	flag.StringVar(&config.github_apikey, "github_apikey", "", "Github API key. (Optional)")
	flag.StringVar(&config.github_repo, "github_repo", "", "Github repository where the secrets will be added. (Required)")
	flag.StringVar(&config.text, "text", "", "Text to encrypt. Local mode! Will not be stored to Github secrets yet. (Optional)")
	flag.BoolVar(&config.debug, "debug", false, "Enable debug logging. (Optional)")
	flag.Parse()

	if config.debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

// validate checks if all necessary arguments were given on the command line
func (config *Config) validate() error {
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
