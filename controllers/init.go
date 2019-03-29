package controllers

import (
	"os"
	"net/http"
	"html/template"
	"../utils"
	"../models"
	"../settings"
)

type dataConnections struct {
    Connections []string
    Error int
}

type dataEmail struct {
    Emails []models.Email
}

type dataMessages struct {
	Connections []string
    Messages []string
    TempConnect string
}

type profileInfo struct {
    Name string
    Info string
    Connections []string
}

type fileArchive struct {
    TempConnect string
    Files []string
}

func redirectTo(to string, w http.ResponseWriter, r *http.Request) {
    switch to {
        case "404": page404(w, r)
        case "archive": ArchivePage(w, r)
    }
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
	t, err := template.ParseFiles(settings.PATH_VIEWS + "index.html", settings.PATH_VIEWS + "page404.html")
	utils.CheckError(err)
    t.Execute(w, nil)
}
