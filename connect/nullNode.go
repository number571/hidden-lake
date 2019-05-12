package connect

import (
    "../settings"
)

func nullNode(username string) {
    settings.Mutex.Lock()
    delete(settings.Node.Login, username)
    delete(settings.Node.Connection, username)
    delete(settings.Node.PublicKey, username)
    delete(settings.Node.SessionKey.P2P, username)
    delete(settings.Node.Address.P2P, username)
    delete(settings.Messages.NewDataExistLocal, username)
    delete(settings.Messages.CurrentIdLocal, username)
    settings.Mutex.Unlock()
}
