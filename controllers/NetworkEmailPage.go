package controllers

import (
	"net/http"
	"html/template"
	"../utils"
	"../settings"
)

func NetworkEmailPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/network/email/" {
		redirectTo("404", w, r)
		return
	}

	tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "index.html", settings.PATH_VIEWS + "network_email.html")
    utils.CheckError(err)
    tmpl.Execute(w, nil)
}
