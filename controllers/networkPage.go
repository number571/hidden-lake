package controllers

import (
    "strings"
    "net/http"
    "html/template"
    "../utils"
    "../connect"
    "../settings"
)

func networkPage(w http.ResponseWriter, r *http.Request) {
    if !settings.User.Auth {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    
    if r.Method == "POST" {
       r.ParseForm()

        if _, ok := r.Form["add"]; ok {
            connect.Connect(strings.Split(r.FormValue("addr"), " "))

        } else if _, ok := r.Form["delete"]; ok {
            connect.Disconnect(strings.Split(r.FormValue("addr"), " "))
        }
    }

    var data = struct {
        Auth bool
        Login string
        Connections []string
    } {
        Auth: true,
        Login: settings.User.Login,
        Connections: settings.User.Connections,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
