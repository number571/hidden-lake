package connect

import (
    "../settings"
)

func nullNode(username string) {
    settings.Mutex.Lock()
    delete(settings.User.NodeLogin, username)
    delete(settings.User.NodeConnection, username)
    delete(settings.User.NodePublicKey, username)
    delete(settings.User.NodeSessionKey, username)
    delete(settings.User.NodeAddress, username)
    delete(settings.Messages.NewDataExistLocal, username)
    delete(settings.Messages.CurrentIdLocal, username)
    settings.Mutex.Unlock()
}
