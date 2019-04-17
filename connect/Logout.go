package connect

import (
	"../settings"
)

func Logout() {
	Disconnect(settings.User.Connections)
	settings.Mutex.Lock()
	settings.DataBase.Close()
	if settings.ServerListenTCP != nil {
		settings.ServerListenTCP.Close()
	}
	settings.User.Connections = []string{}
	settings.GoroutinesIsRun = false
	settings.User.Auth = false
	settings.User.Login = ""
	settings.User.Password = []byte{}
	settings.Mutex.Unlock()
}
