package connect

import (
	"../utils"
	"../settings"
)

// Delete from database node and messages.
func DisconnectF2F(username string) {
	if _, ok := settings.Node.Address.F2F[username]; !ok { return }
	settings.Mutex.Lock()
	delete(settings.Node.ConnectionMode, username)
	delete(settings.Node.Address.F2F, username)
	delete(settings.Node.SessionKey.F2F, username)
	delete(settings.Messages.NewDataExistLocal, username)
	delete(settings.Messages.CurrentIdLocal, username)
	_, err := settings.DataBase.Exec(`
DELETE FROM Local` + username + ` WHERE Mode = 'F2F';
DELETE FROM ConnectionsF2F WHERE User = '` + username + `';
`)
	settings.Mutex.Unlock()
	utils.CheckError(err)
}
