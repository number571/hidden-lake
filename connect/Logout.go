package connect

import (
	"../settings"
)

func Logout() {
	settings.Mutex.Lock()
	settings.DataBase.Close()
	if settings.ServerListenTCP != nil {
		settings.ServerListenTCP.Close()
	}
	settings.Node.Address.P2P = nil
	settings.GoroutinesIsRun = false
	settings.User.Auth = false
	settings.User.Login = ""
	settings.User.Password = []byte{}
	settings.Mutex.Unlock()
}
