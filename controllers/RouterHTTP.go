package controllers

import (
	"net/http"
	"../settings"
)

func RouterHTTP() {
	http.Handle("/static/", http.StripPrefix(
		"/static/", 
		handleFileServer(http.Dir(settings.PATH_STATIC))),
	)

	http.HandleFunc("/", indexPage)
	http.HandleFunc("/about", aboutPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/logout", logoutPage)
	http.HandleFunc("/archive", archivePage)
	http.HandleFunc("/settings", settingsPage)
	http.HandleFunc("/network/", networkPage)
	http.HandleFunc("/network/chat/", networkChatPage)
	http.HandleFunc("/network/email/", networkEmailPage)
	http.HandleFunc("/network/email/read/", networkEmailReadPage)
	http.HandleFunc("/network/email/write", networkEmailWritePage)
	http.HandleFunc("/network/archive/", networkArchivePage)
	http.HandleFunc("/network/profile/", networkProfilePage)

	http.HandleFunc("/api/mode", apiMode)
	http.HandleFunc("/api/chat/", apiChatLocal)
	http.HandleFunc("/api/chat/global/", apiChatGlobal)
}
