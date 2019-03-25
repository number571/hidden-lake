package connect

import (
    "../utils"
    "../settings"
)

func DeleteLocalMessages(slice []string) {
    settings.Mutex.Lock()
    for _, user := range slice {
        _, err := settings.DataBase.Exec("DELETE FROM Local" + user)
        utils.CheckError(err)
    }
    settings.Mutex.Unlock()
}