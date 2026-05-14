package main

import (
	"embed"
	"path"
)

//go:embed web
var webFS embed.FS

// readEmbeddedWeb returns a file from the embedded web/ directory (e.g. "index.html").
func readEmbeddedWeb(name string) ([]byte, error) {
	return webFS.ReadFile(path.Join("web", name))
}
