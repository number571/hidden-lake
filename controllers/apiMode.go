package controllers

import (
    "net/http"
    "encoding/json"
    "../models"
    "../settings"
)

func apiMode(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    settings.Mutex.Lock()
    if settings.User.Mode == models.P2P_mode {
        settings.User.Mode = models.F2F_mode
    } else if settings.User.Mode == models.F2F_mode {
        settings.User.Mode = models.P2P_mode
    }
    settings.Mutex.Unlock()

    var mode = settings.CurrentMode()
    var data = models.PackageHTTP {
        Exists: true,
        Head: settings.HEAD_MESSAGE,
        Body: mode,
    }
    json.NewEncoder(w).Encode(data)
}
