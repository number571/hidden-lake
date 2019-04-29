package controllers

import (
    "fmt"
    "strings"
    "net/http"
    "html/template"
    "../utils"
    "../models"
    "../crypto"
    "../connect"
    "../settings"
)

func networkChatGlobal(w http.ResponseWriter, r *http.Request, list_of_status []models.ConnStatus) {
    if !settings.User.Auth {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    go func() {
        settings.Messages.NewDataExistGlobal <- true
    }()

    if r.Method == "POST" {
        r.ParseForm()
        
        if _, ok := r.Form["delete_message"]; ok {
            connect.DeleteGlobalMessages()

        } else if _, ok := r.Form["send_message"]; ok {
            var message = strings.TrimSpace(r.FormValue("text"))
            if message != "" {
                for username := range settings.User.NodeAddress {
                    var new_pack = settings.PackageTCP {
                        From: models.From {
                        Name: settings.User.Hash,
                        },
                        To: username,
                        Head: models.Head {
                            Header: settings.HEAD_MESSAGE,
                            Mode: settings.MODE_GLOBAL,
                        },
                        Body: message,
                    }
                    connect.CreateRedirectPackage(&new_pack)
                    connect.SendInitRedirectPackage(new_pack)
                }
                settings.Mutex.Lock()
                _, err := settings.DataBase.Exec(
                    "INSERT INTO GlobalMessages (User, Body) VALUES ($1, $2)",
                    settings.User.Hash,
                    crypto.Encrypt(
                        settings.User.Password,
                        fmt.Sprintf("[%s]: %s\n", settings.User.Login, message),
                    ),
                )
                settings.Messages.CurrentIdGlobal++
                settings.Mutex.Unlock()
                utils.CheckError(err)
            }
        }
    }

    rows, err := settings.DataBase.Query("SELECT Body FROM GlobalMessages ORDER BY Id")
    utils.CheckError(err)

    var messages []string
    var message string

    for rows.Next() {
        rows.Scan(&message)
        messages = append(messages, crypto.Decrypt(settings.User.Password, message))
    }

    var data = struct {
        Auth bool
        Login string
        Messages []string
        Connections []models.ConnStatus
    } {
        Auth: true,
        Login: settings.User.Login,
        Messages: messages,
        Connections: list_of_status,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_chat.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
