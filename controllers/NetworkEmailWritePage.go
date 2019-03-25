package controllers

import (
    "time"
    "net/http"
    "html/template"
    "../utils"
    "../models"
    "../settings"
    "../connect"
)

func NetworkEmailWritePage(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/network/email/write/" {
        redirectTo("404", w, r)
        return
    }

    var data = dataConnections{
        Connections: settings.User.Connections,
        Error: 0,
    }

    if r.Method == "POST" {
        var node = r.FormValue("node")
        if node == "none" {
            data = dataConnections{
                Connections: settings.User.Connections,
                Error: 1,
            }
        } else {
            var new_pack = settings.Package {
                From: models.From {
                    Name: settings.User.Name,
                },
                To: node,
                Head: models.Head {
                    Header: settings.HEAD_EMAIL,
                    Mode: settings.MODE_SAVE,
                }, 
                Body: 
                    r.FormValue("title") + settings.SEPARATOR +
                    r.FormValue("body") + settings.SEPARATOR +
                    time.Now().Format(time.RFC850),
            }
            connect.SendEncryptedPackage(new_pack)

            data = dataConnections{
                Connections: settings.User.Connections,
                Error: -1,
            }
        }
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "index.html", settings.PATH_VIEWS + "network_email_write.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
