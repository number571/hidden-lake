package connect

import (
    "../settings"
)

// Delete all localdata in P2P node.
func nullNode(username string) {
    settings.Mutex.Lock()
    delete(settings.Node.ConnectionMode, username)
    delete(settings.Node.PublicKey, username)
    delete(settings.Node.SessionKey.P2P, username)
    delete(settings.Node.Address.P2P, username)
    delete(settings.Node.Address.C_S, username)
    delete(settings.Messages.NewDataExistLocal, username)
    delete(settings.Messages.CurrentIdLocal, username)
    settings.Mutex.Unlock()
}
