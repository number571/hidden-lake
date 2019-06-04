package connect

import (
    "encoding/hex"
    "../utils"
    "../crypto"
    "../models"
    "../settings"
)

// Set password for node by username.
func ConnectF2F(username, address, password string) {
    var (
        session_key = crypto.HashSum([]byte(password))
        encrypted_address = crypto.Encrypt(settings.User.Password, address)
        encrypted_session_key = crypto.Encrypt(settings.User.Password, hex.EncodeToString(session_key))
    )
    settings.Mutex.Lock()
    settings.Node.ConnectionMode[username] = models.CONN
    settings.Messages.NewDataExistLocal[username] = make(chan bool)
    settings.Node.Address.F2F[username] = address
    settings.Node.SessionKey.F2F[username] = session_key
    _, err := settings.DataBase.Exec(`
CREATE TABLE IF NOT EXISTS Local` + username + ` (
Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
User VARCHAR(64),
Mode VARCHAR(3),
Body TEXT
);

INSERT INTO ConnectionsF2F(User, Address, SessionKey) 
SELECT '` + username + `', '` + encrypted_address + `', '` + encrypted_session_key + `'
WHERE NOT EXISTS(SELECT 1 FROM ConnectionsF2F WHERE User = '` + username + `');

UPDATE ConnectionsF2F 
SET Address = '` + encrypted_address + `', SessionKey = '` + encrypted_session_key + `'
WHERE User = '` + username + `';
`)
    settings.Mutex.Unlock()
    utils.CheckError(err)
}
