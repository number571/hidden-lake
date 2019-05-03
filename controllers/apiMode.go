package controllers

import (
    "net/http"
    "encoding/json"
    "../settings"
)

func apiMode(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    settings.Mutex.Lock()
    settings.User.ModeF2F = !settings.User.ModeF2F
    settings.Mutex.Unlock()

    var mode = settings.CurrentMode()
    var data = settings.PackageHTTP {
        Exists: true,
        Head: settings.HEAD_MESSAGE,
        Body: mode,
    }
    json.NewEncoder(w).Encode(data)
}
