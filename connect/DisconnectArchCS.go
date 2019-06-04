package connect

import (
	"../models"
	"../settings"
)

func DisconnectArchCS() {
	settings.Mutex.Lock()
	if settings.Node.ConnServer.Addr != nil {
		settings.Node.ConnServer.Addr.Close()
	}
    settings.Node.ConnServer.Hash = ""
    settings.User.Mode = models.C_S_mode
    settings.Mutex.Unlock()
}
