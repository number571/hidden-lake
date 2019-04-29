package controllers

import (
    "strings"
    "net/http"
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
        rows, err := settings.DataBase.Query("SELECT User, Login FROM Connections")
        utils.CheckError(err)

        var list_of_status []models.ConnStatus
        var temp_list models.ConnStatus
        for rows.Next() {
            rows.Scan(&temp_list.User, &temp_list.Login)
            temp_list.Login = crypto.Decrypt(settings.User.Password, temp_list.Login)
            list_of_status = append(list_of_status, temp_list)
        }
        rows.Close()

        for index, value := range list_of_status {
            if _, ok := settings.User.NodeAddress[value.User]; ok {
                list_of_status[index].Status = true
            }
        }

        var data = struct {
            Auth bool
            Login string
            Connections []models.ConnStatus
        } {
            Auth: true,
            Login: settings.User.Login,
            Connections: list_of_status,
        }

        tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_profile.html")
        utils.CheckError(err)
        tmpl.Execute(w, data)
        return
    }

    var result = strings.TrimPrefix(r.URL.Path, "/network/profile/")
    var row = settings.DataBase.QueryRow("SELECT User, Login, PublicKey FROM Connections WHERE User = $1", result)

    var user_hash, login, public_data string
    row.Scan(&user_hash, &login, &public_data)

    var status = false
    if _, ok := settings.User.NodeAddress[user_hash]; ok {
        status = true
    }

    var data = struct {
        UserHash string
        UserLogin string
        PublicKey string
        Status bool
        Auth bool
        Login string
    } {
        UserHash: user_hash,
        UserLogin: crypto.Decrypt(settings.User.Password, login),
        PublicKey: public_data,
        Status: status,
        Auth: true,
        Login: settings.User.Login,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_profile_X.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
