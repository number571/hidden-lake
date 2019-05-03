package connect

import (
    "time"
    "encoding/hex"
    "../utils"
    "../models"
    "../crypto"
    "../settings"
)

// Connect to P2P nodes.
func Connect(slice []string, not_check bool) {
    next:
    for _, addr := range slice {
        var address = settings.User.IPv4 + settings.User.Port

        if addr == address {
            continue
        }

        for _, node_addr := range settings.User.NodeAddress {
            if addr == node_addr {
                continue next
            }
        }

        var new_pack = settings.PackageTCP {
            From: models.From {
                Address: address,
                Login: settings.User.Login,
                Name: settings.User.Hash,
            },
            Head: models.Head {
                Header: settings.HEAD_CONNECT,
                Mode: settings.MODE_GLOBAL,
            },
        }

        if not_check {
            new_pack.Head.Mode = settings.MODE_READ
            new_pack.Body = hex.EncodeToString([]byte(settings.User.PublicData))
        }

        sendAddrPackage(addr, new_pack)
        time.Sleep(time.Millisecond * 500)
    }
}

// Connect to F2F node.
func ConnectF2F(username, address, password string) {
    var (
        session_key = crypto.HashSum([]byte(password))
        encrypted_address = crypto.Encrypt(settings.User.Password, address)
        encrypted_session_key = crypto.Encrypt(settings.User.Password, hex.EncodeToString(session_key))
    )
    settings.Mutex.Lock()
    settings.Messages.NewDataExistLocal[username] = make(chan bool)
    settings.User.NodeAddressF2F[username] = address
    settings.User.NodeSessionKeyF2F[username] = session_key
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
