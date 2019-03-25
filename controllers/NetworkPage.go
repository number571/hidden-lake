package controllers

import (
    "strings"
    "net/http"
    "html/template"
    "../utils"
    "../connect"
    "../settings"
)

func NetworkPage(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
       r.ParseForm()

        if _, ok := r.Form["add"]; ok {
            connect.Connect(strings.Split(r.FormValue("addr"), " "))

        } else if _, ok := r.Form["delete"]; ok {
            connect.Disconnect(strings.Split(r.FormValue("addr"), " "))
        }
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "index.html", settings.PATH_VIEWS + "network.html")
    utils.CheckError(err)
    tmpl.Execute(w, settings.User)
}
