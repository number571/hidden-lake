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

    var node_address = settings.CurrentNodeAddress()
    if !settings.User.ModeF2F {
        for index, value := range list_of_status {
            _, list_of_status[index].Status = node_address[value.User]
        }
    }
    
    if r.URL.Path == "/network/chat/" {
        networkChatGlobal(w, r, list_of_status)
        return
    }

    if r.Method == "POST" {
        r.ParseForm()

        if _, ok := r.Form["delete_message"]; ok {
            settings.DeleteLocalMessages([]string{settings.User.TempConnect})

        } else if _, ok := r.Form["send_message"]; ok {
            var message = strings.TrimSpace(r.FormValue("text"))
            if message != "" {
                var new_pack = settings.PackageTCP {
                    From: models.From {
                        Name: settings.User.Hash,
                    },
                    To: settings.User.TempConnect,
                    Head: models.Head {
                        Header: settings.HEAD_MESSAGE,
                        Mode: settings.MODE_LOCAL,
                    },
                    Body: message,
                }
                var mode = settings.CurrentMode()
                if settings.User.ModeF2F {
                    connect.CreateRedirectF2FPackage(&new_pack, settings.User.TempConnect)
                }
                if settings.User.ModeF2F {
                    for username := range settings.User.NodeAddressF2F {
                        new_pack.To = username
                        connect.SendPackage(new_pack, true)
                    }
                } else {
                    connect.SendPackage(new_pack, false)
                }

                var from = settings.User.Login
                if settings.User.ModeF2F { from = settings.User.Hash }

                settings.Mutex.Lock()
                _, err := settings.DataBase.Exec(
                    "INSERT INTO Local" + settings.User.TempConnect + " (User, Mode, Body) VALUES ($1, $2, $3)",
                    settings.User.Hash,
                    mode,
                    crypto.Encrypt(
                        settings.User.Password,
                        fmt.Sprintf("[%s]: %s\n", from, message),
                    ),
                )
                settings.Mutex.Unlock()
                utils.CheckError(err)
            }
        }
    }
    
    var result = strings.Split(strings.TrimPrefix(r.URL.Path, "/network/chat/"), "/")[0]

    settings.Mutex.Lock()
    settings.User.TempConnect = ""
    settings.Mutex.Unlock()

    var user_login string

    // if _, ok := node_address[result]; ok {
    //     settings.Mutex.Lock()
    //     settings.User.TempConnect = result
    //     settings.Mutex.Unlock()
    //     if !settings.User.ModeF2F {
    //         user_login = settings.User.NodeLogin[result]
    //     }
    // }

    for _, value := range list_of_status {
        if value.User == result {
            settings.Mutex.Lock()
            settings.User.TempConnect = result
            settings.Mutex.Unlock()
            if !settings.User.ModeF2F {
                user_login = settings.User.NodeLogin[result]
            }
        }
    }

    if settings.User.TempConnect == "" {
        redirectTo("404", w, r)
        return
    }

    settings.Mutex.Lock()
    settings.Messages.CurrentIdLocal[settings.User.TempConnect] = 0
    settings.Mutex.Unlock()

    var mode = settings.CurrentMode()
    _, status := settings.User.NodeAddress[settings.User.TempConnect]
    go func() {
        settings.Messages.NewDataExistLocal[settings.User.TempConnect] <- true
    }()

    rows, err = settings.DataBase.Query(
        "SELECT Body FROM Local" + settings.User.TempConnect + " WHERE Mode = $1 ORDER BY Id",
        mode,
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
        Connections []models.ConnStatus
        Status bool
        UserLogin string
        Auth bool
        Login string
        ModeF2F bool
        TempConnect string
    } {
        Messages: messages,
        Connections: list_of_status,
        Status: status,
        UserLogin: user_login,
        Auth: true,
        Login: settings.User.Login,
        ModeF2F: settings.User.ModeF2F,
        TempConnect: settings.User.TempConnect,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_chat_X.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
