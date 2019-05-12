package connect

import (
	"../utils"
    "../settings"
)

func DisconnectP2P(username string) {
	nullNode(username)
	settings.Mutex.Lock()
	_, err := settings.DataBase.Exec(`
DELETE FROM Local` + username + ` WHERE Mode = 'P2P';
DELETE FROM Connections WHERE User = '` + username + `';
`)
	settings.Mutex.Unlock()
	utils.CheckError(err)
}
