package settings

import (
	"../utils"
)

func DeleteGlobalMessages() {
    Mutex.Lock()
    _, err := DataBase.Exec("DELETE FROM GlobalMessages")
    Mutex.Unlock()
    utils.CheckError(err)
}
