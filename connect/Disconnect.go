package connect

import (
	"../utils"
    "../settings"
)

// Disconnect from P2P node.
func Disconnect(username string) {
	nullNode(username)
	settings.Mutex.Lock()
	_, err := settings.DataBase.Exec(`
DELETE FROM Local` + username + ` WHERE Mode = 'P2P';
DELETE FROM Connections WHERE User = '` + username + `';
`)
	settings.Mutex.Unlock()
	utils.CheckError(err)
}

// Disconnect from F2F node.
func DisconnectF2F(username string) {
	if _, ok := settings.User.NodeAddressF2F[username]; !ok { return }
	settings.Mutex.Lock()
	delete(settings.User.NodeAddressF2F, username)
	delete(settings.User.NodeSessionKeyF2F, username)
	delete(settings.Messages.NewDataExistLocal, username)
	delete(settings.Messages.CurrentIdLocal, username)
	_, err := settings.DataBase.Exec(`
DELETE FROM Local` + username + ` WHERE Mode = 'F2F';
DELETE FROM ConnectionsF2F WHERE User = '` + username + `';
`)
	settings.Mutex.Unlock()
	utils.CheckError(err)
}
