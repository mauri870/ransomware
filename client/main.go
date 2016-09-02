package client

import (
	"net/http"
	"net/url"
	"strings"
)

const (
	SERVER_URL = "http://localhost:8080"
)

// Call the server
func CallServer(method string, endpoint string, data url.Values) (*http.Response, error) {
	req, err := http.NewRequest(method, SERVER_URL+endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return new(http.Response), err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return new(http.Response), err
	}

	return res, nil
}
