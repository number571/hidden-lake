package controllers

import (
    "net/http"
)

func redirectTo(to string, w http.ResponseWriter, r *http.Request) {
    switch to {
        case "404": page404(w, r)
        case "archive": ArchivePage(w, r)
        case "network": networkSettingPage(w, r)
    }
}
