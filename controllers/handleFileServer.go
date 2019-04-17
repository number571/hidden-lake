package controllers

import (
	"os"
	"net/http"
)

func handleFileServer(fs http.FileSystem) http.Handler {
	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		if _, err := fs.Open(r.URL.Path); os.IsNotExist(err) {
			redirectTo("404", w, r)
			return
		}
		http.FileServer(fs).ServeHTTP(w, r)
	})
}
