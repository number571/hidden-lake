package views

import (
	"net/http"
	"html/template"
	"../utils"
	"../settings"
)

func page404(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(settings.PATH_TEMPLATES + "base.html", settings.PATH_TEMPLATES + "404_page.html")
	utils.CheckWarning(err)
    t.Execute(w, nil)
}
