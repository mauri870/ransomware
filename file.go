package main

import (
	"os"
)

// Type File wrap an os.FileInfo
type File struct {
	os.FileInfo
	Extension string
	Path      string
}
