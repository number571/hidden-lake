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

func networkChatPage(w http.ResponseWriter, r *http.Request) {
    if !settings.User.Auth {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

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
    
    if r.URL.Path == "/network/chat/" {
        networkChatGlobal(w, r, list_of_status)
        return
    }

    if r.Method == "POST" {
        r.ParseForm()

        if _, ok := r.Form["delete_message"]; ok {
            connect.DeleteLocalMessages([]string{settings.User.TempConnect})

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
                // connect.SendEncryptedPackage(new_pack)
                connect.CreateRedirectPackage(&new_pack)
                connect.SendRedirectPackage(new_pack)
                settings.Mutex.Lock()
                _, err := settings.DataBase.Exec(
                    "INSERT INTO Local" + settings.User.TempConnect + " (User, Body) VALUES ($1, $2)",
                    settings.User.Hash,
                    crypto.Encrypt(
                        settings.User.Password,
                        fmt.Sprintf("[%s]: %s\n", settings.User.Hash, message),
                    ),
                )
                settings.Messages.CurrentIdLocal[settings.User.TempConnect]++
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
    for _, value := range list_of_status {
        if value.User == result {
            settings.Mutex.Lock()
            settings.User.TempConnect = value.User
            user_login = value.Login
            settings.Mutex.Unlock()
            break
        }
    }

    if settings.User.TempConnect == "" {
        redirectTo("404", w, r)
        return
    }

    var status = false
    if _, ok := settings.User.NodeAddress[settings.User.TempConnect]; ok {
        status = true
    }

    go func() {
        settings.Messages.NewDataExistLocal[settings.User.TempConnect] <- true
    }()

    rows, err = settings.DataBase.Query("SELECT Body FROM Local" + settings.User.TempConnect + " ORDER BY Id")
    utils.CheckError(err)

    var messages []string
    var message string
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
        TempConnect string
    } {
        Messages: messages,
        Connections: list_of_status,
        Status: status,
        UserLogin: user_login,
        Auth: true,
        Login: settings.User.Login,
        TempConnect: settings.User.TempConnect,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_chat_X.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
