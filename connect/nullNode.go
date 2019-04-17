package connect

import (
    "../utils"
    "../settings"
)

func nullNode(username string) {
    settings.Mutex.Lock()
    settings.User.NodeAddress[username] = ""
    settings.User.NodeLogin[username] = ""
    settings.User.NodeConnection[username] = 0
    settings.User.NodePublicKey[username]  = nil
    settings.User.NodeSessionKey[username] = nil
    settings.User.Connections = utils.RemoveByElem(
        settings.User.Connections,
        username,
    )
    delete(settings.Messages.NewDataExistLocal, username)
    delete(settings.Messages.CurrentIdLocal, username)
    settings.Mutex.Unlock()
}
