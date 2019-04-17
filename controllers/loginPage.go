package controllers

import (
    "net/http"
    "html/template"
    "../utils"
    "../connect"
    "../settings"
)

func loginPage(w http.ResponseWriter, r *http.Request) {
    var data = struct {
        Code int8
        Result string
        Auth bool
        Hash string
        Login string
    }{}

    var code int8 = 0

    if r.Method == "POST" && !settings.User.Auth {
        var login = r.FormValue("login")
        var password = r.FormValue("password")

        code = settings.Authorization(login, password)

        switch code {
            case  1: data.Result = "Length of login = 0 byte"
            case  2: data.Result = "Length of login > 64 bytes"
            case  3: data.Result = "File password.hash undefined"
            case  4: data.Result = "Wrong login or password"
            default: 
                if !settings.GoroutinesIsRun && settings.User.Port != "" {
                    settings.Mutex.Lock()
                    settings.GoroutinesIsRun = true
                    settings.Mutex.Unlock()
                    go connect.ServerTCP()
                    go connect.FindConnects(10)
                }
                http.Redirect(w, r, "/", http.StatusSeeOther)
        }
    }

    data.Code = code
    data.Auth = settings.User.Auth
    data.Hash = settings.User.Hash
    data.Login = settings.User.Login

    tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "login.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
