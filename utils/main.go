package utils

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"os/user"
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

// Return a list containing the letters used by the current drives
func GetDrives() (letters []string) {
	for _, drive := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		_, err := os.Open(string(drive) + ":\\")
		if err == nil {
			letters = append(letters, string(drive)+":\\")
		}
	}
	return
}

// Return an *os.User instance representing the current user
func GetCurrentUser() *user.User {
	usr, err := user.Current()
	if err != nil {
		return &user.User{}
	}
	return usr
}
