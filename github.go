package main

import "net/http"

// writeSecret writes the value to a secret name in the configured repository
func writeSecret(name string, value string) error {
	// curl -u arvidwarnecke0042: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 -H "Accept: application/vnd.github.v3+json" https://api.github.com/repos/otto-ec/dv_opal-permission/actions/secrets/public-key
	// ==

	req, err := http.NewRequest("GET", "https://api.github.com/repos/otto-ec/dv_opal-permission/actions/secrets/public-key", nil)
	if err != nil {
		// handle err
	}
	req.SetBasicAuth("", "")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	return err
}
