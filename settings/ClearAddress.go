package settings

import (
	"../utils"
	"../models"
)

func ClearAddress() {
	Mutex.Lock()
    User.IPv4 = ""
    User.Port = ""
    User.Mode = models.C_S_mode
    _, err := DataBase.Exec("DELETE FROM Address")
    Mutex.Unlock()
    utils.CheckError(err)
}
