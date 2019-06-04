package connect

import (
	"../models"
	"../settings"
)

// Delete and stop all localdata by current usernode.
func Logout() {
	settings.Mutex.Lock()
	settings.DataBase.Close()
	if settings.ServerListenTCP != nil {
		settings.ServerListenTCP.Close()
	}
	settings.User.Mode = models.C_S_mode
	settings.Node.ConnServer.Hash = ""
	settings.Node.ConnServer.Addr = nil
	settings.Node.Address.P2P = nil
	settings.GoroutinesIsRun = false
	settings.User.Auth = false
	settings.User.Login = ""
	settings.User.Password = []byte{}
	settings.Mutex.Unlock()
}
