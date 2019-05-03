package controllers

import (
    "strings"
    "net/http"
    "html/template"
    "../utils"
    "../crypto"
    "../connect"
    "../settings"
)

func settingsPage(w http.ResponseWriter, r *http.Request) {
	if !settings.User.Auth {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    if r.Method == "POST" {
        var connects = r.FormValue("conn")

        settings.Mutex.Lock()
    	settings.User.IPv4 = r.FormValue("ipv4")
    	settings.User.Port = ":" + r.FormValue("port")
        settings.User.DefaultConnections = strings.Split(connects, "\r\n")  
        settings.Mutex.Unlock()

        utils.WriteFile(settings.FILE_CONNECTS, crypto.Encrypt(settings.User.Password, connects))
        utils.WriteFile(settings.FILE_SETTINGS, crypto.Encrypt(
            settings.User.Password, 
            settings.User.IPv4 + settings.User.Port,
        ))
        
    	if !settings.GoroutinesIsRun {
            settings.Mutex.Lock()
            settings.GoroutinesIsRun = true
            settings.Mutex.Unlock()
			go connect.ServerTCP()
            go connect.FindConnects(10)
    	}
    }

    var port = strings.TrimPrefix(settings.User.Port, ":")

    var data = struct{
        IPv4 string
        Port string
        Conn string
        PublicKey string
        Auth bool
        Hash string
        Login string
        ModeF2F bool
    } {
        IPv4: settings.User.IPv4,
        Port: port,
        Conn: crypto.Decrypt(settings.User.Password, utils.ReadFile(settings.FILE_CONNECTS)),
        PublicKey: settings.User.PublicData,
        Auth: true,
        Hash: settings.User.Hash,
        Login: settings.User.Login,
        ModeF2F: settings.User.ModeF2F,
    }

	tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "settings.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
