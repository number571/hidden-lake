package connect

import (
    "../settings"
)

func DeleteLocalMessages(slice []string) {
    settings.Mutex.Lock()
    for _, check := range slice {
        for _, username := range settings.User.Connections {
            if username == check {
                settings.User.LocalMessages[username] = []string{}
                break
            }
        }
    }
    settings.Mutex.Unlock()
}