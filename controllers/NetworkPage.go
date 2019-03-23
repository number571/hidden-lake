package controllers

import (
    "fmt"
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

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "profile_node.html")
    utils.CheckWarning(err)
    tmpl.Execute(w, data)
}

func NetworkArchivePage(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/network/archive/" {
        redirectTo("404", w, r)
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

    var new_pack = settings.Package {
        From: models.From {
            Name: settings.User.Name,
        },
        To: settings.User.TempConnect,
        Head: models.Head {
            Header: settings.HEAD_ARCHIVE,
            Mode: settings.MODE_READ_LIST,
        },
    }

    connect.SendEncryptedPackage(new_pack)
    time.Sleep(time.Second * settings.TIME_SLEEP)

    var data = fileArchive{
        TempConnect: settings.User.TempConnect,
    }

    for _, file := range settings.User.TempArchive {
        if file != "" {
            data.Files = append(data.Files, file)
        }
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "archive_node.html")
    utils.CheckWarning(err)
    tmpl.Execute(w, data)
}

func NetworkGlobalPage(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/network/global/" {
        redirectTo("404", w, r)
        return
    }

    if r.Method == "POST" {
        switch r.FormValue("act") {
            case "Del_messages":
                connect.DeleteGlobalMessages()

            case "Send_message":
                message := r.FormValue("text")
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
                    settings.User.GlobalMessages = append(
                        settings.User.GlobalMessages,
                        fmt.Sprintf("[%s]: %s\n", settings.User.Name, message),
                    )
                    settings.Mutex.Unlock()
                }
            }
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "global_chat.html")
    utils.CheckWarning(err)
    tmpl.Execute(w, settings.User)
}

func NetworkPage(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/network/" {
        redirectTo("network", w, r)
        return
    }

    if r.Method == "POST" {
        switch r.FormValue("act") {
            case "Delete": 
                connect.Disconnect([]string{settings.User.TempConnect})

            case "Del_messages":
                connect.DeleteLocalMessages([]string{settings.User.TempConnect})

            case "Send_message":
                message := r.FormValue("text")
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
                    settings.User.LocalMessages[settings.User.TempConnect] = append(
                        settings.User.LocalMessages[settings.User.TempConnect],
                        fmt.Sprintf("[%s]: %s\n", settings.User.Name, message),
                    )
                    settings.Mutex.Unlock()
                }
        }
    }

    var result = strings.TrimPrefix(r.URL.Path, "/network/")

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

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "chat.html")
    utils.CheckWarning(err)
    tmpl.Execute(w, settings.User)
}

func networkSettingPage(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        switch r.FormValue("act") {
            case "Add_connection": 
                connect.Connect(strings.Split(r.FormValue("addrs"), " "))

            case "Del_connection": 
                connect.Disconnect(strings.Split(r.FormValue("addrs"), " "))
        }
    }

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "network.html")
    utils.CheckWarning(err)
    tmpl.Execute(w, settings.User)
}
