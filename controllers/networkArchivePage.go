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

    var node_address = settings.CurrentNodeAddress()
    
    if r.URL.Path == "/network/archive/" {
        var connects = settings.MakeConnects(node_address)

        var data = struct {
            Auth bool
            Login string
            ModeF2F bool
            Connections []string
        } {
            Auth: true,
            Login: settings.User.Login,
            ModeF2F: settings.User.ModeF2F,
            Connections: connects,
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

    if _, ok := node_address[result]; ok {
        settings.Mutex.Lock()
        settings.User.TempConnect = result
        settings.Mutex.Unlock()
    }

    if settings.User.TempConnect == "" {
        redirectTo("404", w, r)
        return
    }

    var new_pack = models.PackageTCP {
        From: models.From {
            Name: settings.CurrentHash(),
        },
        To: settings.User.TempConnect,
        Head: models.Head {
            Title: settings.HEAD_ARCHIVE,
            Mode: settings.MODE_READ_LIST,
        },
    }

    settings.Mutex.Lock()
    settings.User.TempArchive = nil
    settings.Mutex.Unlock()

    connect.SendPackage(new_pack, settings.CurrentModeNet())

    var seconds = 0
check_again:
    if settings.User.TempArchive == nil && seconds < 10 {
        time.Sleep(time.Second * 1)
        seconds++
        goto check_again
    }

    var data = struct {
        Auth bool
        Login string
        ModeF2F bool
        Files []string
        TempConnect string
    } {
        Auth: true,
        Login: settings.User.Login,
        ModeF2F: settings.User.ModeF2F,
        TempConnect: settings.User.TempConnect,
        Files: settings.User.TempArchive,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_archive_X.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
