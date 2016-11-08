package client

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/http2"
)

const (
	SERVER_URL = "https://localhost:8080"
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
