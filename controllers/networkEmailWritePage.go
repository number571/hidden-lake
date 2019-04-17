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

func networkEmailWritePage(w http.ResponseWriter, r *http.Request) {
    if !settings.User.Auth {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    
    if r.URL.Path != "/network/email/write" {
        redirectTo("404", w, r)
        return
    }

    var err_page int8

    if r.Method == "POST" {
        var node = r.FormValue("node")
        if node == "none" {
            err_page = 1
        } else {
            var new_pack = settings.PackageTCP {
                From: models.From {
                    Name: settings.User.Hash,
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

            err_page = -1
        }
    }

    var data = struct {
        Auth bool
        Login string
        Connections []string
        Error int8
    } {
        Auth: true,
        Login: settings.User.Login,
        Connections: settings.User.Connections,
        Error: err_page,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_email_write.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
