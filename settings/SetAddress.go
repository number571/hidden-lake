package settings

import (
	"../models"
)

func SetAddress(ipv4, port string) {
	ClearAddress()
	Mutex.Lock()
    User.IPv4 = ipv4
    User.Port = ":" + port
    if NeedF2FMode {
    	User.Mode = models.F2F_mode
    } else {
    	User.Mode = models.P2P_mode
    }
    Mutex.Unlock()
    SaveAddress(ipv4, port)
}
