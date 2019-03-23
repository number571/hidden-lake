package connect

import (
	"../settings"
)

func DeleteGlobalMessages() {
    settings.Mutex.Lock()
    settings.User.GlobalMessages = []string{}
    settings.Mutex.Unlock()
}
