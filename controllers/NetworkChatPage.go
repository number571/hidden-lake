package controllers

import (
    "fmt"
    "strings"
    "net/http"
    "html/template"
    "../utils"
    "../models"
    "../connect"
    "../settings"
)

func NetworkChatPage(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/network/chat/" {
        var data = dataMessages {
            Connections: settings.User.Connections,
        }

        tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "index.html", settings.PATH_VIEWS + "network_chat.html")
        utils.CheckError(err)
        tmpl.Execute(w, data)
        return
    }

    if r.Method == "POST" {
        r.ParseForm()

        if _, ok := r.Form["delete_message"]; ok {
            connect.DeleteLocalMessages([]string{settings.User.TempConnect})

        } else if _, ok := r.Form["send_message"]; ok {
            var message = strings.TrimSpace(r.FormValue("text"))
            if message != "" {
                var new_pack = settings.Package {
                    From: models.From {
                        Name: settings.User.Name,
                    },
                    To: settings.User.TempConnect,
                    Head: models.Head {
                        Header: settings.HEAD_MESSAGE,
                        Mode: settings.MODE_LOCAL,
                    },
                    Body: message,
                }
                connect.SendEncryptedPackage(new_pack)
                settings.Mutex.Lock()
                _, err := settings.DataBase.Exec(
                    "INSERT INTO Local" + settings.User.TempConnect + " (User, Body) VALUES ($1, $2)",
                    settings.User.Name,
                    fmt.Sprintf("[%s]: %s\n", settings.User.Name, message),
                )
                settings.Mutex.Unlock()
                utils.CheckError(err)
            }
        }
    }
    
    var result = strings.TrimPrefix(r.URL.Path, "/network/chat/")

    settings.Mutex.Lock()
    settings.User.TempConnect = ""
    settings.Mutex.Unlock()

    for _, username := range settings.User.Connections {
        if username == result {
            settings.Mutex.Lock()
            settings.User.TempConnect = username
            settings.Mutex.Unlock()
            break
        }
    }

    if settings.User.TempConnect == "" {
        redirectTo("404", w, r)
        return
    }

    rows, err := settings.DataBase.Query("SELECT Body FROM Local" + settings.User.TempConnect + " ORDER BY Id")
    utils.CheckError(err)

    var messages []string
    var message string

    for rows.Next() {
        rows.Scan(&message)
        messages = append(messages, message)
    }

    rows.Close()

    var data = dataMessages {
        Messages: messages,
        TempConnect: settings.User.TempConnect,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "index.html", settings.PATH_VIEWS + "network_chat_X.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
