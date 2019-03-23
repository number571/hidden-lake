package controllers

import (
	"os"
	"net/http"
	"html/template"
	"../utils"
	"../settings"
)

type profileInfo struct {
    Name string
    Info string
    Connections []string
}

type fileArchive struct {
    TempConnect string
    Files []string
}

func HandleFileServer(fs http.FileSystem) http.Handler {
	return http.HandlerFunc( func(w http.ResponseWriter, r *http.Request) {
		if _, err := fs.Open(r.URL.Path); os.IsNotExist(err) {
			redirectTo("404", w, r)
			return
		}
		http.FileServer(fs).ServeHTTP(w, r)
	})
}

func page404(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "page404.html")
	utils.CheckWarning(err)
    t.Execute(w, nil)
}
