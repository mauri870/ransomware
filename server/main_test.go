package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServerHandler(t *testing.T) {
	ts := httptest.NewServer(server())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Error(err)
	}

	if err != nil {
		t.Error(err)
	}

	if res.StatusCode != 404 {
		t.Error(err)
	}
}
