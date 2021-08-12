package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// writeSecret writes the value to a secret name in the configured repository
func writeSecret(name string, value string) error {
	var urlStr = "https://api.github.com/repos/otto-ec/" + config.github_repo + "/actions/secrets/" + name
	var jsonStr = []byte(`{"key_id":"` + config.pubkey.Key_id + `","encrypted_value":"` + value + `"}`)

	req, err := http.NewRequest("PUT", urlStr, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println(err)
		os.Exit(127)
	}
	req.SetBasicAuth(config.github_user, config.github_apikey)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(127)
	}
	defer resp.Body.Close()
	log.Infof("Testing something.\n")
	return err
}
