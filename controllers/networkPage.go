package controllers

import (
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
            if settings.User.ModeF2F {
                var name, addr, pasw string
                name = r.FormValue("name")
                addr = r.FormValue("addr")
                pasw = r.FormValue("pasw")
                if name != "" && addr != "" && pasw != "" {
                    connect.ConnectF2F(name, addr, pasw)
                }
            } else {
                connect.ConnectP2PMerge(r.FormValue("addr"))
            }
        }

        if _, ok := r.Form["disconnect"]; ok {
            if settings.User.ModeF2F {
                connect.DisconnectF2F(r.FormValue("name"))
            } else {
                connect.DisconnectP2P(r.FormValue("name"))
            }
        }
    }

    var (
        node_address = settings.CurrentNodeAddress()
        connects = settings.MakeConnects(node_address)
    )

    var data = struct {
        Auth bool
        Login string
        ModeF2F bool
        Connections []string
    } {
        Auth: true,
        Login: settings.User.Login,
        ModeF2F: settings.User.ModeF2F,
        Connections: connects,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
