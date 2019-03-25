package controllers

import (
	"net/http"
	"html/template"
	"../utils"
	"../settings"
)

func AboutPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/about/" {
		redirectTo("404", w, r)
		return
	}

	tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "index.html", settings.PATH_VIEWS + "about.html")
    utils.CheckError(err)
    tmpl.Execute(w, nil)
}
