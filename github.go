package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

// writeSecret writes the value to a secret name in the configured repository
func writeSecret(name string, value string) error {
	var jsonStr = []byte(`{"key_id":"` + config.pubkey.Key_id + `","encrypted_value":"` + value + `"}`)

	// TODO This has to be changed: github_repo and secret_name
	req, err := http.NewRequest("PUT", "https://api.github.com/repos/otto-ec/dv_opal-permission/actions/secrets/AWS_TEST_SECRET2", bytes.NewBuffer(jsonStr))
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
	return err
}
