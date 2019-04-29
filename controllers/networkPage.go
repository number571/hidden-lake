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
            connect.Connect(strings.Split(r.FormValue("addr"), " "), false)
        }
    }

    var (
        connects = make([]string, len(settings.User.NodeAddress))
        index uint32
    )
    for username := range settings.User.NodeAddress {
        connects[index] = username
        index++
    }

    var data = struct {
        Auth bool
        Login string
        Connections []string
    } {
        Auth: true,
        Login: settings.User.Login,
        Connections: connects,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
