package views

import (
	"net/http"
	"html/template"
	"../utils"
	"../settings"
)

type fileArchive struct {
    TempConnect string
    Files []string
}

func IndexPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		redirectTo("404", w, r)
		return
	}

	tmpl, err := template.ParseFiles(settings.PATH_TEMPLATES + "base.html")
    utils.CheckWarning(err)
    tmpl.Execute(w, nil)
}
