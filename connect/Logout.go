package connect

import (
	"../settings"
)

func Logout() {
	var (
		connects = make([]string, len(settings.User.NodeAddress))
		index uint32
	)
	for username := range settings.User.NodeAddress {
		connects[index] = username
		index++
	}
	settings.Mutex.Lock()
	settings.DataBase.Close()
	if settings.ServerListenTCP != nil {
		settings.ServerListenTCP.Close()
	}
	settings.User.NodeAddress = nil
	settings.GoroutinesIsRun = false
	settings.User.Auth = false
	settings.User.Login = ""
	settings.User.Password = []byte{}
	settings.Mutex.Unlock()
}
