package controllers

import (
	"net/http"
	"html/template"
	"../utils"
	"../settings"
)

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/profile/" {
		redirectTo("404", w, r)
		return
	}

	if r.Method == "POST" {
		r.ParseForm()

		if _, ok := r.Form["set_info"]; ok {
			settings.User.Info = r.FormValue("info")
		}
	}

	tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "index.html", settings.PATH_VIEWS + "profile.html")
    utils.CheckError(err)
    tmpl.Execute(w, settings.User)
}
