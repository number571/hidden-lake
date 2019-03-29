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

func NetworkChatGlobalPage(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/network/chat/global/" {
        redirectTo("404", w, r)
        return
    }

    if r.Method == "POST" {
        r.ParseForm()
        
        if _, ok := r.Form["delete_message"]; ok {
            connect.DeleteGlobalMessages()

        } else if _, ok := r.Form["send_message"]; ok {
            var message = strings.TrimSpace(r.FormValue("text"))
            if message != "" {
                for _, username := range settings.User.Connections {
                    var new_pack = settings.Package {
                        From: models.From {
                        Name: settings.User.Name,
                        },
                        To: username,
                        Head: models.Head {
                            Header: settings.HEAD_MESSAGE,
                            Mode: settings.MODE_GLOBAL,
                        },
                        Body: message,
                    }
                    connect.SendEncryptedPackage(new_pack)
                }
                settings.Mutex.Lock()
                _, err := settings.DataBase.Exec(
                    "INSERT INTO GlobalMessages (User, Body) VALUES ($1, $2)",
                    settings.User.Name,
                    fmt.Sprintf("[%s]: %s\n", settings.User.Name, message),
                )
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
        messages = append(messages, message)
    }

    var data = dataMessages {
        Messages: messages,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "index.html", settings.PATH_VIEWS + "network_chat_global.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
