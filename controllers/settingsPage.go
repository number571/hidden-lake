package controllers

import (
    "strings"
    "net/http"
    "html/template"
    "../utils"
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
        
        settings.SaveDefaultConnections(settings.User.DefaultConnections)
        settings.SaveAddress(settings.User.IPv4, settings.User.Port)
        
    	if !settings.GoroutinesIsRun {
            settings.Mutex.Lock()
            settings.GoroutinesIsRun = true
            settings.Mutex.Unlock()
			go connect.ServerTCP()
            go connect.CheckConnects()
    	}
    }

    var public_key string
    if settings.User.ModeF2F {
        public_key = settings.User.Public.Data.F2F
    } else {
        public_key = settings.User.Public.Data.P2P
    }

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
        Port: strings.TrimPrefix(settings.User.Port, ":"),
        Conn: strings.Join(settings.User.DefaultConnections, "\r\n"),
        PublicKey: public_key,
        Auth: true,
        Hash: settings.CurrentHash(),
        Login: settings.User.Login,
        ModeF2F: settings.User.ModeF2F,
    }

	tmpl, err := template.ParseFiles(settings.PATH_VIEWS + "base.html", settings.PATH_VIEWS + "settings.html")
    utils.CheckError(err)
    tmpl.Execute(w, data)
}
