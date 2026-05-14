package main

import _ "embed"
import "net/http"

//go:embed select.html
var frontendHTML []byte

func serveFrontend() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(frontendHTML)
	})
}