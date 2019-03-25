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


func NetworkProfilePage(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/network/profile/" {
        redirectTo("404", w, r)
        return
    }

    var result = strings.TrimPrefix(r.URL.Path, "/network/profile/")

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

    var new_pack = settings.Package {
        From: models.From {
            Name: settings.User.Name,
        },
        To: settings.User.TempConnect,
        Head: models.Head {
            Header: settings.HEAD_PROFILE,
            Mode: settings.MODE_READ,
        },
    }

    settings.Mutex.Lock()
    settings.User.TempProfile = []string{}
    settings.Mutex.Unlock()

    connect.SendEncryptedPackage(new_pack)
    time.Sleep(time.Second * settings.TIME_SLEEP)

    if len(settings.User.TempProfile) != 3 {
        redirectTo("404", w, r)
        return
    }

    var data = profileInfo {
        Name: settings.User.TempProfile[0],
        Info: settings.User.TempProfile[1],
        Connections: strings.Split(settings.User.TempProfile[2], settings.SEPARATOR_ADDRESS),
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "index.html", settings.PATH_VIEWS + "network_profile.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
