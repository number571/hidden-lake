package connect

import (
    "../utils"
    "../settings"
)

func DeleteLocalMessages(slice []string) {
    settings.Mutex.Lock()
    for _, user := range slice {
        for _, check := range settings.User.Connections {
            if check == user { 
                _, err := settings.DataBase.Exec("DELETE FROM Local" + user)
                utils.CheckError(err)
            }
        }
    }
    settings.Mutex.Unlock()
}
