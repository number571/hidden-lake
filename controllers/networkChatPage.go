package controllers

import (
    "fmt"
    "strings"
    "net/http"
    "database/sql"
    "html/template"
    "../utils"
    "../models"
    "../crypto"
    "../connect"
    "../settings"
)

func networkChatPage(w http.ResponseWriter, r *http.Request) {
    if !settings.User.Auth {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
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

    if r.URL.Path == "/network/chat/" {
        networkChatGlobal(w, r, list_of_users)
        return
    }

    var result = strings.Split(strings.TrimPrefix(r.URL.Path, "/network/chat/"), "/")[0]

    if r.Method == "POST" {
        r.ParseForm()

        if _, ok := r.Form["delete_message"]; ok {
            settings.DeleteLocalMessages([]string{result})

        } else if _, ok := r.Form["send_message"]; ok {
            var message = strings.TrimSpace(r.FormValue("text"))
            if message != "" {
                var hashname = settings.CurrentHash()
                var mode = settings.CurrentMode()
                var new_pack = models.PackageTCP {
                    From: models.From {
                        Hash: hashname,
                    },
                    To: models.To {
                        Hash: result,
                    },
                    Head: models.Head {
                        Title: settings.HEAD_MESSAGE,
                        Mode: settings.MODE_LOCAL,
                    },
                    Body: message,
                }

                connect.SendPackage(new_pack, settings.User.Mode)
                
                // if settings.User.Mode == models.F2F_mode {
                //     connect.CreateRedirectF2FPackage(&new_pack, settings.User.TempConnect)
                // }

                // if settings.User.Mode == models.F2F_mode {
                //     for username := range settings.Node.Address.F2F {
                //         new_pack.To.Hash = username
                //         connect.SendPackage(new_pack, settings.CurrentModeNet())
                //     }
                // } else {
                //     connect.SendPackage(new_pack, settings.CurrentModeNet())
                // }

                // var from = settings.User.Login
                // if settings.User.Mode == models.F2F_mode { from = hashname }

                settings.Mutex.Lock()
                _, err := settings.DataBase.Exec(
                    "INSERT INTO Local" + result + " (User, Mode, Body) VALUES ($1, $2, $3)",
                    settings.CurrentHash(),
                    mode,
                    crypto.Encrypt(
                        settings.User.Password,
                        fmt.Sprintf("(%s)[%s]: %s\n", mode, hashname, message),
                    ),
                )
                settings.Mutex.Unlock()
                utils.CheckError(err)
            }
        }
    }

    var user_is_not_exist = true
    for _, value := range list_of_users {
        if value == result {
            user_is_not_exist = false
        }
    }

    if user_is_not_exist {
        redirectTo("404", w, r)
        return
    }

    settings.Mutex.Lock()
    settings.Messages.CurrentIdLocal[result] = 0
    settings.Mutex.Unlock()

    var status bool
    if settings.User.Mode == models.P2P_mode {
        _, status = settings.Node.Address.P2P[result]
    }
    go func() {
        settings.Messages.NewDataExistLocal[result] <- true
    }()

    rows, err = settings.DataBase.Query(
        "SELECT Body FROM Local" + result + " WHERE Mode = $1 ORDER BY Id",
        settings.CurrentMode(),
    )
    utils.CheckError(err)

    var (
        messages []string
        message string
    )
    for rows.Next() {
        rows.Scan(&message)
        messages = append(messages, crypto.Decrypt(settings.User.Password, message))
    }
    rows.Close()

    var data = struct {
        Messages []string
        Connections []string
        Status bool
        Auth bool
        Login string
        Mode string
        TempConnect string
    } {
        Messages: messages,
        Connections: list_of_users,
        Status: status,
        Auth: true,
        Login: settings.User.Login,
        Mode: settings.CurrentMode(),
        TempConnect: result,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_chat_X.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
