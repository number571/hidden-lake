package controllers

import (
    "net/http"
    "database/sql"
    "html/template"
    "../utils"
    "../models"
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
            if settings.User.Mode == models.F2F_mode {
                var name, addr, pasw string
                name = r.FormValue("name")
                addr = r.FormValue("addr")
                pasw = r.FormValue("pasw")
                if name != "" && addr != "" && pasw != "" {
                    connect.ConnectF2F(name, addr, pasw)
                }
            } else if settings.User.Mode == models.P2P_mode {
                connect.ConnectP2PMerge(r.FormValue("addr"))
            } else {
                if settings.Node.ConnServer.Addr != nil {
                    connect.DisconnectArchCS()
                }
                settings.Node.ConnServer.Addr = connect.GetConnection(r.FormValue("addr"))
                connect.ConnectArchCS(settings.Node.ConnServer.Addr, true)
            }
        }

        if _, ok := r.Form["disconnect"]; ok {
            if settings.User.Mode == models.F2F_mode {
                connect.DisconnectF2F(r.FormValue("name"))
            } else {
                connect.DisconnectP2P(r.FormValue("name"))
            } 
        }
    }

    var (
        rows *sql.Rows
        err error
    )

    if settings.User.Mode == models.F2F_mode {
        rows, err = settings.DataBase.Query("SELECT User FROM ConnectionsF2F")
    } else {
        rows, err = settings.DataBase.Query("SELECT User FROM Connections")
    }
    utils.CheckError(err)

    var list_of_users []string
    var username string

    for rows.Next() {
        rows.Scan(&username)
        list_of_users = append(list_of_users, username)
    }
    
    rows.Close()

    var data = struct {
        Auth bool
        Login string
        Mode string
        Connections []string
    } {
        Auth: true,
        Login: settings.User.Login,
        Mode: settings.CurrentMode(),
        Connections: list_of_users,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
