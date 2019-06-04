package controllers

import (
    "strings"
    "net/http"
    "database/sql"
    "html/template"
    "../utils"
    "../crypto"
    "../models"
    "../settings"
)

func networkProfilePage(w http.ResponseWriter, r *http.Request) {
    if !settings.User.Auth {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    if r.URL.Path == "/network/profile/" {
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

        tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_profile.html")
        utils.CheckError(err)
        tmpl.Execute(w, data)
        return
    }

    var (
        result = strings.TrimPrefix(r.URL.Path, "/network/profile/")
        row *sql.Row
        address, public_data string
    )

    if settings.User.Mode == models.F2F_mode {
        row = settings.DataBase.QueryRow("SELECT Address FROM ConnectionsF2F WHERE User = $1", result)
        row.Scan(&address)
        if address == "" { 
            redirectTo("404", w, r) 
            return
        }
        address = crypto.Decrypt(settings.User.Password, address)

    } else {
        row = settings.DataBase.QueryRow("SELECT PublicKey FROM Connections WHERE User = $1", result)
        row.Scan(&public_data)
        if public_data == "" { 
            redirectTo("404", w, r) 
            return
        }
    }

    var data = struct {
        UserHash string
        PublicKey string
        Address string
        Auth bool
        Login string
        Mode string
    } {
        UserHash: result,
        PublicKey: public_data,
        Address: address,
        Auth: true,
        Login: settings.User.Login,
        Mode: settings.CurrentMode(),
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_profile_X.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
