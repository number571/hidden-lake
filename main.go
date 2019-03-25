package main

import (
	"os"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"./connect"
	"./settings"
)

import (
	"net/http"
	"./utils"
	"./controllers"
)

func main() {
	settings.CreateDatabase("database.db")

	db, err := sql.Open("sqlite3", "database.db")
	utils.CheckError(err)
	defer db.Close()

	settings.DataBase = db
	settings.Initialization(os.Args)

	fmt.Println("Server is listening...")

	go connect.ClientTCP()
	go connect.ServerTCP()

	http.Handle("/static/", http.StripPrefix(
		"/static/", 
		controllers.HandleFileServer(http.Dir(settings.PATH_STATIC))),
	)

	http.HandleFunc("/", controllers.IndexPage)
	http.HandleFunc("/about/", controllers.AboutPage)
	http.HandleFunc("/profile/", controllers.ProfilePage)
	http.HandleFunc("/archive/", controllers.ArchivePage)
	http.HandleFunc("/network/", controllers.NetworkPage)
	http.HandleFunc("/network/chat/", controllers.NetworkChatPage)
	http.HandleFunc("/network/chat/global/", controllers.NetworkChatGlobalPage)
	http.HandleFunc("/network/email/", controllers.NetworkEmailPage)
	http.HandleFunc("/network/email/read/", controllers.NetworkEmailReadPage)
	http.HandleFunc("/network/email/write/", controllers.NetworkEmailWritePage)
	http.HandleFunc("/network/archive/", controllers.NetworkArchivePage)
	http.HandleFunc("/network/profile/", controllers.NetworkProfilePage)

	utils.CheckError(http.ListenAndServe(settings.IPV4_HTTP + settings.PORT_HTTP, nil))
}
