package controllers

import (
    "time"
    "strings"
    "net/http"
    "html/template"
    "../utils"
    "../models"
    "../connect"
    "../settings"
)

func networkArchivePage(w http.ResponseWriter, r *http.Request) {
    if !settings.User.Auth {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }
    
    if r.URL.Path == "/network/archive/" {
        var data = struct {
            Auth bool
            Login string
            Connections []string
        } {
            Auth: true,
            Login: settings.User.Login,
            Connections: settings.User.Connections,
        }

        tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_archive.html")
        utils.CheckError(err)
        tmpl.Execute(w, data)
        return
    }

    var result = strings.TrimPrefix(r.URL.Path, "/network/archive/")

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

    var new_pack = settings.PackageTCP {
        From: models.From {
            Name: settings.User.Hash,
        },
        To: settings.User.TempConnect,
        Head: models.Head {
            Header: settings.HEAD_ARCHIVE,
            Mode: settings.MODE_READ_LIST,
        },
    }

    connect.SendEncryptedPackage(new_pack)
    time.Sleep(time.Second * settings.TIME_SLEEP)

    var data = struct {
        Auth bool
        Login string
        Files []string
        TempConnect string
    } {
        Auth: true,
        Login: settings.User.Login,
        TempConnect: settings.User.TempConnect,
    }

    for _, file := range settings.User.TempArchive {
        if file != "" {
            data.Files = append(data.Files, file)
        }
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_archive_X.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
