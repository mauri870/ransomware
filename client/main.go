package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	SERVER_URL = "http://localhost:8080"
)

var (
	ErrorBadJson = errors.New("Bad response content from server, expecting valid JSON")
)

// Make a request to server
func CallServer(endpoint string) (*http.Response, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s", SERVER_URL, endpoint), nil)
	req.Header.Set("User-Agent", "Ransomware/1.0")
	res, err := client.Do(req)
	if err != nil {
		return new(http.Response), err
	}

	return res, nil
}

// Parse a Json response from server
func ParseServerJson(body []byte) (map[string]string, error) {
	var content map[string]string
	if err := json.Unmarshal(body, &content); err != nil {
		return content, ErrorBadJson
	}

	return content, nil
}
