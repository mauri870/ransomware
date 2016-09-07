package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFound(t *testing.T) {
	ts := httptest.NewServer(server())
	defer ts.Close()
	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Error(err)
	}

	if ApiResponseNotFound+"\n" != string(body) {
		t.Errorf("Expected %s got %s", ApiResponseNotFound, string(body))
	}

	if res.StatusCode != 404 {
		t.Errorf("Expected %s got %s", 404, res.StatusCode)
	}
}
