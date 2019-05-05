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

    var node_address = settings.CurrentNodeAddress()
    
    if r.URL.Path == "/network/profile/" {
        var (
            rows *sql.Rows
            err error
        )

        if settings.User.ModeF2F {
            rows, err = settings.DataBase.Query("SELECT User FROM ConnectionsF2F")
        } else {
            rows, err = settings.DataBase.Query("SELECT User, Login FROM Connections")
        }
        utils.CheckError(err)

        var list_of_status []models.ConnStatus
        var temp_list models.ConnStatus

        if settings.User.ModeF2F {
            for rows.Next() {
                rows.Scan(&temp_list.User)
                list_of_status = append(list_of_status, temp_list)
            }
        } else {
            for rows.Next() {
                rows.Scan(&temp_list.User, &temp_list.Login)
                temp_list.Login = crypto.Decrypt(settings.User.Password, temp_list.Login)
                list_of_status = append(list_of_status, temp_list)
            }
        }
        
        rows.Close()

        if settings.User.ModeF2F {
            for index, value := range list_of_status {
                if _, ok := node_address[value.User]; ok {
                    list_of_status[index].Status = true
                }
            }
        }
        
        var data = struct {
            Auth bool
            Login string
            ModeF2F bool
            Connections []models.ConnStatus
        } {
            Auth: true,
            Login: settings.User.Login,
            ModeF2F: settings.User.ModeF2F,
            Connections: list_of_status,
        }

        tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_profile.html")
        utils.CheckError(err)
        tmpl.Execute(w, data)
        return
    }

    var (
        result = strings.TrimPrefix(r.URL.Path, "/network/profile/")
        row *sql.Row
        address, login, public_data string
    )

    if settings.User.ModeF2F {
        row = settings.DataBase.QueryRow("SELECT Address FROM ConnectionsF2F WHERE User = $1", result)
        row.Scan(&address)
        if address == "" { 
            redirectTo("404", w, r) 
            return
        }
        address = crypto.Decrypt(settings.User.Password, address)
    } else {
        row = settings.DataBase.QueryRow("SELECT Login, PublicKey FROM Connections WHERE User = $1", result)
        row.Scan(&login, &public_data)
        if login == "" { 
            redirectTo("404", w, r) 
            return
        }
        login = crypto.Decrypt(settings.User.Password, login)
    }

    _, status := node_address[result]

    var data = struct {
        UserHash string
        UserLogin string
        PublicKey string
        Address string
        Status bool
        Auth bool
        Login string
        ModeF2F bool
    } {
        UserHash: result,
        UserLogin: login,
        PublicKey: public_data,
        Address: address,
        Status: status,
        Auth: true,
        Login: settings.User.Login,
        ModeF2F: settings.User.ModeF2F,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_profile_X.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
