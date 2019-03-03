package views

import (
	"net/http"
	"html/template"
	"../utils"
	"../settings"
)

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/profile" {
		redirectTo("404", w, r)
		return
	}

	if r.Method == "POST" {
		switch r.FormValue("act") {
			case "Set_Info": 
				settings.User.Info = r.FormValue("info")
		}
	}

	tmpl, err := template.ParseFiles(settings.PATH_TEMPLATES + "base.html", settings.PATH_TEMPLATES + "profile.html")
    utils.CheckWarning(err)
    tmpl.Execute(w, settings.User)
}
