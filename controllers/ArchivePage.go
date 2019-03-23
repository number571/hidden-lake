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

func ArchivePage(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/archive/" {
        redirectTo("404", w, r)
        return
    }

    if r.Method == "POST" {
        var filename = r.FormValue("filename")
        switch r.FormValue("act") {
            case "Download": 
                var new_pack = settings.Package {
                    From: models.From {
                        Name: settings.User.Name,
                    },
                    To: settings.User.TempConnect,
                    Head: models.Head {
                        Header: settings.HEAD_ARCHIVE,
                        Mode: settings.MODE_READ_FILE,
                    },
                    Body: filename,
                }
                connect.SendEncryptedPackage(new_pack)
                time.Sleep(time.Second * settings.TIME_SLEEP)

            case "Delete": 
                err := os.Remove(settings.PATH_ARCHIVE + filename)
                utils.CheckWarning(err)

            case "Copy":
                utils.WriteFile(settings.PATH_ARCHIVE + "copy_" + filename, 
                    utils.ReadFile(settings.PATH_ARCHIVE + filename),
                )
        }
    }

    files, err := ioutil.ReadDir(settings.PATH_ARCHIVE)
    utils.CheckError(err)

    var data = fileArchive {
        TempConnect: settings.User.TempConnect,
    }

    for _, file := range files {
        data.Files = append(data.Files, file.Name())
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "archive.html")
    utils.CheckWarning(err)
    tmpl.Execute(w, data)
}
