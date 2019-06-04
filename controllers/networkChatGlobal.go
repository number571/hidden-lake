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

func networkChatGlobal(w http.ResponseWriter, r *http.Request, list_of_users []string) {
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
                // node_address = settings.CurrentNodeAddress()
            )

            var new_pack = models.PackageTCP {
                From: models.From {
                    Hash: settings.CurrentHash(),
                },
                Head: models.Head {
                    Title: settings.HEAD_MESSAGE,
                    Mode: settings.MODE_GLOBAL,
                },
                Body: message,
            }

            if len(splited) >= 3 {
                switch splited[0] {
                    case settings.TERM_SEND: 
                        new_pack.To.Hash = splited[1]
                        new_pack.Head.Mode = settings.MODE_LOCAL
                        new_pack.Body = strings.Join(splited[2:], " ")
                        message = "(" + new_pack.To.Hash + ")" + new_pack.Body
                }
            }

            connect.SendPackage(new_pack, settings.User.Mode)

            // if settings.User.Mode == models.F2F_mode {
            //     connect.CreateRedirectF2FPackage(&new_pack, new_pack.To.Hash)
            // }

            // for username := range node_address {
            //     new_pack.To.Hash = username
            //     connect.SendPackage(new_pack, settings.CurrentModeNet())
            // }

            // var from = settings.User.Login
            // if settings.User.Mode == models.F2F_mode { from = settings.CurrentHash() }

            settings.Mutex.Lock()
            _, err := settings.DataBase.Exec(
                "INSERT INTO GlobalMessages (User, Mode, Body) VALUES ($1, $2, $3)",
                settings.CurrentHash(),
                mode,
                crypto.Encrypt(
                    settings.User.Password,
                    fmt.Sprintf("(%s)[%s]: %s\n", settings.CurrentMode(), settings.CurrentHash(), message),
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
        Mode string
        Messages []string
        Connections []string
    } {
        Auth: true,
        Login: settings.User.Login,
        Mode: settings.CurrentMode(),
        Messages: messages,
        Connections: list_of_users,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_chat.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
