package views

import (
    "fmt"
    "time"
    "strings"
    "net/http"
    "html/template"
    "../conn"
    "../utils"
    "../models"
    "../settings"
)

type ProfileInfo struct {
    Name string
    IPv4 string
    Port string
    Info string
}

func NetworkProfilePage(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path == "/network/profile/" {
        redirectTo("404", w, r)
        return
    }

    var result = strings.TrimPrefix(r.URL.Path, "/network/profile/")

    settings.Mutex.Lock()
    settings.User.TempConnect = ""
    settings.Mutex.Unlock()

    for _, addr := range settings.User.Connections {
        if addr == result {
            settings.Mutex.Lock()
            settings.User.TempConnect = addr
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
            Address: settings.User.IPv4 + settings.User.Port,
            Name: settings.User.Name,
        },
        To: settings.User.TempConnect,
        Head: models.Head {
            Header: settings.HEAD_PROFILE,
            Mode: settings.MODE_GET,
        },
    }

    settings.Mutex.Lock()
    settings.User.TempProfile = []string{}
    settings.Mutex.Unlock()

    conn.SendEncryptedPackage(new_pack)
    time.Sleep(time.Second * settings.TIME_SLEEP)

    if len(settings.User.TempProfile) != 4 {
        redirectTo("404", w, r)
        return
    }

    var data = ProfileInfo {
        Name: settings.User.TempProfile[0],
        IPv4: settings.User.TempProfile[1],
        Port: settings.User.TempProfile[2],
        Info: settings.User.TempProfile[3],
    }

    tmpl, err := template.ParseFiles(settings.PATH_TEMPLATES + "base.html", settings.PATH_TEMPLATES + "profile_node.html")
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

    for _, addr := range settings.User.Connections {
        if addr == result {
            settings.Mutex.Lock()
            settings.User.TempConnect = addr
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
            Address: settings.User.IPv4 + settings.User.Port,
            Name: settings.User.Name,
        },
        To: settings.User.TempConnect,
        Head: models.Head {
            Header: settings.HEAD_ARCHIVE,
            Mode: settings.MODE_GET_LIST,
        },
    }

    conn.SendEncryptedPackage(new_pack)
    time.Sleep(time.Second * settings.TIME_SLEEP)

    var data = FileArchive{
        TempConnect: settings.User.TempConnect,
    }
    for _, file := range settings.User.TempArchive {
        if file != "" {
            data.Files = append(data.Files, file)
        }
    }

    tmpl, err := template.ParseFiles(settings.PATH_TEMPLATES + "base.html", settings.PATH_TEMPLATES + "archive_node.html")
    utils.CheckWarning(err)
    tmpl.Execute(w, data)
}

func NetworkGlobalPage(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        switch r.FormValue("act") {
            case "Del_messages":
                conn.DeleteGlobalMessages()

            case "Send_message":
                message := r.FormValue("text")
                if message != "" {
                    for _, addr := range settings.User.Connections {
                        var new_pack = settings.Package {
                            From: models.From {
                                Address: settings.User.IPv4 + settings.User.Port,
                                Name: settings.User.Name,
                            },
                            To: addr,
                            Head: models.Head {
                                Header: settings.HEAD_MESSAGE,
                                Mode: settings.MODE_GLOBAL,
                            },
                            Body: message,
                        }
                        conn.SendEncryptedPackage(new_pack)
                    }
                    settings.Mutex.Lock()
                    settings.User.GlobalMessages = append(
                        settings.User.GlobalMessages,
                        fmt.Sprintf("[%s]: %s\n", 
                            settings.User.IPv4 + settings.User.Port + "/" + 
                            settings.User.Name, message),
                    )
                    settings.Mutex.Unlock()
                }
            }
    }

    tmpl, err := template.ParseFiles(settings.PATH_TEMPLATES + "base.html", settings.PATH_TEMPLATES + "global_chat.html")
    utils.CheckWarning(err)
    tmpl.Execute(w, settings.User)
}

func NetworkPage(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        switch r.FormValue("act") {
            case "Delete": 
                conn.Disconnect([]string{settings.User.TempConnect})

            case "Del_messages":
                conn.DeleteLocalMessages([]string{settings.User.TempConnect})

            case "Send_message":
                message := r.FormValue("text")
                if message != "" {
                    var new_pack = settings.Package {
                        From: models.From {
                            Address: settings.User.IPv4 + settings.User.Port,
                            Name: settings.User.Name,
                        },
                        To: settings.User.TempConnect,
                        Head: models.Head {
                            Header: settings.HEAD_MESSAGE,
                            Mode: settings.MODE_LOCAL,
                        },
                        Body: message,
                    }
                    conn.SendEncryptedPackage(new_pack)
                    settings.Mutex.Lock()
                    settings.User.LocalMessages[settings.User.TempConnect] = append(
                        settings.User.LocalMessages[settings.User.TempConnect],
                        fmt.Sprintf("[%s/%s]: %s\n", 
                            settings.User.IPv4 + settings.User.Port, 
                            settings.User.Name, message),
                    )
                    settings.Mutex.Unlock()
                }
        }
    }

    if r.URL.Path == "/network/" {
        redirectTo("network", w, r)
        return
    }

    var result = strings.TrimPrefix(r.URL.Path, "/network/")

    settings.Mutex.Lock()
    settings.User.TempConnect = ""
    settings.Mutex.Unlock()

    for _, addr := range settings.User.Connections {
        if addr == result {
            settings.Mutex.Lock()
            settings.User.TempConnect = addr
            settings.Mutex.Unlock()
            break
        }
    }

    if settings.User.TempConnect == "" {
        redirectTo("404", w, r)
        return
    }

    tmpl, err := template.ParseFiles(settings.PATH_TEMPLATES + "base.html", settings.PATH_TEMPLATES + "chat.html")
    utils.CheckWarning(err)
    tmpl.Execute(w, settings.User)
}

func networkSettingPage(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        switch r.FormValue("act") {
            case "Add_connection": 
                conn.Connect(strings.Split(r.FormValue("addrs"), " "))

            case "Del_connection": 
                conn.Disconnect(strings.Split(r.FormValue("addrs"), " "))
        }
    }

    tmpl, err := template.ParseFiles(settings.PATH_TEMPLATES + "base.html", settings.PATH_TEMPLATES + "network.html")
    utils.CheckWarning(err)
    tmpl.Execute(w, settings.User)
}
