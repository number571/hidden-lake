package controllers

import (
	"net/http"
	"html/template"
	"../utils"
	"../settings"
)

func indexPage(w http.ResponseWriter, r *http.Request) {
	if !settings.User.Auth {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
	
	if r.URL.Path != "/" {
		redirectTo("404", w, r)
		return
	}

	var data = struct {
		Auth bool
		Login string
	} {
		Auth: true,
		Login: settings.User.Login,
	}

	tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "index.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
