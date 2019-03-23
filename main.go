package main

import (
	"os"
	"fmt"
	"./connect"
	"./settings"
)

import (
	"net/http"
	"./utils"
	"./controllers"
)

func main() {
	settings.Initialization(os.Args)

	fmt.Println("Server is listening...")

	go connect.ClientTCP()
	go connect.ServerTCP()

	http.Handle("/static/", http.StripPrefix(
		"/static/", 
		controllers.HandleFileServer(http.Dir(settings.PATH_STATIC))),
	)

	http.HandleFunc("/", controllers.IndexPage)
	http.HandleFunc("/profile/", controllers.ProfilePage)
	http.HandleFunc("/archive/", controllers.ArchivePage)
	http.HandleFunc("/network/", controllers.NetworkPage)
	http.HandleFunc("/network/global/", controllers.NetworkGlobalPage)
	http.HandleFunc("/network/archive/", controllers.NetworkArchivePage)
	http.HandleFunc("/network/profile/", controllers.NetworkProfilePage)

	utils.CheckError(http.ListenAndServe(settings.IPV4_HTTP + settings.PORT_HTTP, nil))
}
