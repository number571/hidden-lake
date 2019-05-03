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

    var mode = settings.CurrentMode()

    if r.Method == "POST" {
        r.ParseForm()
        
        if _, ok := r.Form["delete_message"]; ok {
            settings.DeleteGlobalMessages()

        } else if _, ok := r.Form["send_message"]; ok {
            var message = strings.TrimSpace(r.FormValue("text"))
            if message == "" { goto render_page }

            var (
                splited = strings.Split(message, " ")
                node_address = settings.CurrentNodeAddress()
            )

            var new_pack = settings.PackageTCP {
                From: models.From {
                    Name: settings.User.Hash,
                },
                Head: models.Head {
                    Header: settings.HEAD_MESSAGE,
                    Mode: settings.MODE_GLOBAL,
                },
                Body: message,
            }

            if len(splited) >= 3 {
                switch splited[0] {
                    case settings.TERM_SEND: 
                        new_pack.To = splited[1]
                        new_pack.Head.Mode = settings.MODE_LOCAL
                        new_pack.Body = strings.Join(splited[2:], " ")
                        message = "(" + new_pack.To + ")" + new_pack.Body
                }
            }

            if settings.User.ModeF2F {
                connect.CreateRedirectF2FPackage(&new_pack, new_pack.To)
            }

            for username := range node_address {
                new_pack.To = username
                connect.SendPackage(new_pack, settings.User.ModeF2F)
            }

            var from = settings.User.Login
            if settings.User.ModeF2F { from = settings.User.Hash }

            settings.Mutex.Lock()
            _, err := settings.DataBase.Exec(
                "INSERT INTO GlobalMessages (User, Mode, Body) VALUES ($1, $2, $3)",
                settings.User.Hash,
                mode,
                crypto.Encrypt(
                    settings.User.Password,
                    fmt.Sprintf("[%s]: %s\n", from, message),
                ),
            )
            settings.Messages.CurrentIdGlobal++
            settings.Mutex.Unlock()
            utils.CheckError(err)
        }
    }

render_page:
    settings.Mutex.Lock()
    settings.Messages.CurrentIdGlobal = 0
    settings.Mutex.Unlock()

    rows, err := settings.DataBase.Query(
        "SELECT Body FROM GlobalMessages WHERE Mode = $1 ORDER BY Id",
        mode,
    )
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
        ModeF2F bool
        Messages []string
        Connections []models.ConnStatus
    } {
        Auth: true,
        Login: settings.User.Login,
        ModeF2F: settings.User.ModeF2F,
        Messages: messages,
        Connections: list_of_status,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_chat.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
