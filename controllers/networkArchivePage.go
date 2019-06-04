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
            Mode string
            Connections []string
        } {
            Auth: true,
            Login: settings.User.Login,
            Mode: settings.CurrentMode(),
            Connections: connects,
        }

        tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_archive.html")
        utils.CheckError(err)
        tmpl.Execute(w, data)
        return
    }

    var result = strings.TrimPrefix(r.URL.Path, "/network/archive/")

    var user_is_not_exist = true
    if _, ok := node_address[result]; ok {
        user_is_not_exist = false
    }

    if user_is_not_exist {
        redirectTo("404", w, r)
        return
    }

    var new_pack = models.PackageTCP {
        From: models.From {
            Hash: settings.CurrentHash(),
        },
        To: models.To {
            Hash: result,
        },
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
        Mode models.ModeNet
        Files []string
        TempConnect string
    } {
        Auth: true,
        Login: settings.User.Login,
        Mode: settings.User.Mode,
        TempConnect: result,
        Files: settings.User.TempArchive,
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network_archive_X.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
