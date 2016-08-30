package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// Generate a random alphanumeric string with the given size
func GenerateRandomANString(size int) (string, error) {
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(key)[:size], nil
}

// Check if a value exists on slice
func StringInSlice(search string, slice []string) bool {
	for _, v := range slice {
		if v == search {
			return true
		}
	}
	return false
}
