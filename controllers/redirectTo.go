package controllers

import (
	"net/http"
	"html/template"
	"../utils"
	"../settings"
)

func redirectTo(to string, w http.ResponseWriter, r *http.Request) {
    switch to {
        case "404": page404(w, r)
        case "archive": archivePage(w, r)
    }
}

func page404(w http.ResponseWriter, r *http.Request) {
	if !settings.User.Auth {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

	var data = struct {
		Auth bool
		Login string
		Mode string
	} {
		Auth: true,
		Login: settings.User.Login,
		Mode: settings.CurrentMode(),
	}

	t, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "page404.html")
	utils.CheckError(err)
    t.Execute(w, data)
}
