package connect

import (
    "../utils"
    "../settings"
)

func DeleteLocalMessages(slice []string) {
    settings.Mutex.Lock()
    for _, user := range slice {
        if _, ok := settings.User.NodeAddress[user]; ok {
            _, err := settings.DataBase.Exec("DELETE FROM Local" + user)
            utils.CheckError(err)
        }
    }
    settings.Mutex.Unlock()
}
