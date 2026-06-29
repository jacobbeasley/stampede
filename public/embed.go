package public

import (
	"embed"
	"io/fs"
	"os"
)

//go:embed all:*
var files embed.FS

func FS() fs.FS {
	if os.Getenv("GO_ENV") == "production" {
		return files
	}
	return os.DirFS("public")
}
