package main

import (
	"os"
	"fmt"
	"./clients"
	"./connect"
	"./settings"
	"./controllers"
)

func main() {
	var gui_interface = settings.Initialization(os.Args)
	controllers.RouterHTTP()

	if settings.User.Port != "" {
		settings.GoroutinesIsRun = true
		go connect.ServerTCP()
		go connect.CheckConnects()
	}

	if gui_interface {
		go clients.ClientHTTP()
	}

	fmt.Println("[Server is listening]")
	clients.ClientTCP()
}
