package controllers

import (
    "time"
    "strings"
    "net/http"
    "encoding/json"
    "../utils"
    "../crypto"
    "../settings"
)

func apiChatGlobal(w http.ResponseWriter, r *http.Request) {
    var result = strings.TrimPrefix(r.URL.Path, "/api/chat/global/")
    w.Header().Set("Content-Type", "application/json")

    switch result {
        case "update": 
            updateChatGlobal(w)
        default: 
            json.NewEncoder(w).Encode(settings.PackageHTTP{Head: settings.HEAD_WARNING})
    } 
}

func updateChatGlobal(w http.ResponseWriter) {
    timeout := make (chan bool)
    go func() {
        time.Sleep(time.Second * 30)
        timeout <- true
    }()
 
    select {
        case <-timeout:
            json.NewEncoder(w).Encode(settings.PackageHTTP{Exists:false})
            return

        case <-settings.Messages.NewDataExistGlobal:
            var mode = settings.CurrentMode()
            rows, err := settings.DataBase.Query(
                "SELECT Body FROM GlobalMessages WHERE Mode = $1 ORDER BY Id DESC LIMIT $2",
                mode,
                settings.Messages.CurrentIdGlobal,
            )
            utils.CheckError(err)

            var messages, message string
            for rows.Next() {
                rows.Scan(&message)
                messages += crypto.Decrypt(settings.User.Password, message)
            }

            var data = settings.PackageHTTP {
                Exists: true,
                Head: settings.HEAD_MESSAGE,
                Body: messages,
            }

            settings.Mutex.Lock()
            settings.Messages.CurrentIdGlobal = 0
            settings.Mutex.Unlock()

            json.NewEncoder(w).Encode(data)
            return
    }
}
