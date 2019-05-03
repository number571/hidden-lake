package controllers

import (
    "os"
    "time"
    "net/http"
    "io/ioutil"
    "html/template"
    "../utils"
    "../models"
    "../connect"
    "../settings"
)

func archivePage(w http.ResponseWriter, r *http.Request) {
    if !settings.User.Auth {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    if r.Method == "POST" {
        r.ParseForm()
        var filename = r.FormValue("filename")

        if _, ok := r.Form["delete"]; ok {
            err := os.Remove(settings.PATH_ARCHIVE + filename)
            utils.CheckWarning(err)

        } else if _, ok := r.Form["copy"]; ok {
            utils.WriteFile(
                settings.PATH_ARCHIVE + "copy_" + filename, 
                utils.ReadFile(settings.PATH_ARCHIVE + filename),
            )

        } else if _, ok := r.Form["download"]; ok {
            var new_pack = settings.PackageTCP {
                From: models.From {
                    Name: settings.User.Hash,
                },
                To: settings.User.TempConnect,
                Head: models.Head {
                    Header: settings.HEAD_ARCHIVE,
                    Mode: settings.MODE_READ_FILE,
                },
                Body: filename,
            }
            connect.SendPackage(new_pack, settings.User.ModeF2F)
            time.Sleep(time.Second * settings.TIME_SLEEP) // FIX
        }
    }

    files, err := ioutil.ReadDir(settings.PATH_ARCHIVE)
    utils.CheckError(err)

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
    }

    for _, file := range files {
        data.Files = append(data.Files, file.Name())
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "archive.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
