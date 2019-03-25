package connect

import (
	"../utils"
	"../settings"
)

func DeleteGlobalMessages() {
    settings.Mutex.Lock()
    _, err := settings.DataBase.Exec("DELETE FROM GlobalMessages")
    settings.Mutex.Unlock()
    utils.CheckError(err)
}
