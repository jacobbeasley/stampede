package templates

import (
	"embed"
	"io/fs"
	"os"
)

//go:embed * */*
var files embed.FS

func FS() fs.FS {
	if os.Getenv("GO_ENV") == "production" || os.Getenv("GO_ENV") == "test" {
		return files
	}
	if _, err := os.Stat("templates"); err == nil {
		return os.DirFS("templates")
	}
	if _, err := os.Stat("../templates"); err == nil {
		return os.DirFS("../templates")
	}
	return files
}
