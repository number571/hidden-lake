package main

import (
	"os"
	"fmt"
	"./conn"
	"./settings"
)

import (
	"net/http"
	"./views"
	"./utils"
)

func main() {
	settings.Initialization(os.Args)

	fmt.Println("Server is listening...")

	go conn.ClientTCP()
	go conn.ServerTCP()

	http.Handle("/static/", http.StripPrefix(
		"/static/", 
		views.HandleFileServer(http.Dir(settings.PATH_STATIC))),
	)

	http.HandleFunc("/", views.IndexPage)
	http.HandleFunc("/profile", views.ProfilePage)
	http.HandleFunc("/archive", views.ArchivePage)
	http.HandleFunc("/network/", views.NetworkPage)
	http.HandleFunc("/network/global", views.NetworkGlobalPage)
	http.HandleFunc("/network/archive/", views.NetworkArchivePage)
	http.HandleFunc("/network/profile/", views.NetworkProfilePage)

	utils.CheckError(http.ListenAndServe(settings.IPV4_HTTP + settings.PORT_HTTP, nil))
}
