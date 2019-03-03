package conn

import (
    "net"
    "encoding/hex"
    "encoding/json"
    "../utils"
    "../crypto"
    "../models"
    "../settings"
)

func encrypt(to, data string) string {
    result, _ := crypto.EncryptAES(
        []byte(data),
        settings.User.NodeSessionKey[to],
    )
    return hex.EncodeToString(result)
}

func SendEncryptedPackage(pack settings.Package) {
    if settings.User.NodeConnection[pack.To] == 1 {
        var to = pack.To
        pack = settings.Package {
            From: models.From {
                Address: pack.From.Address,
                Name: encrypt(pack.To, pack.From.Name),
            },
            To: encrypt(pack.To, pack.To),
            Head: models.Head {
                Header: encrypt(pack.To, pack.Header),
                Mode: encrypt(pack.To, pack.Mode),
            }, 
            Body: encrypt(pack.To, pack.Body),
        }
        SendPackage(to, pack)
    } else {
        nullNode(pack.To)
    }
}

func SendPackage(to string, pack settings.Package) {
    conn, err := net.Dial(settings.PROTOCOL_TCP, to)
    if err != nil {
        nullNode(to)
        return
    }

    if settings.User.NodeConnection[to] == 0 {
        settings.Mutex.Lock()
        settings.User.NodeConnection[to] = -1
        settings.Mutex.Unlock()
    }
            
    data, err := json.Marshal(pack)
    utils.CheckError(err)

    conn.Write(data)
    conn.Close()
}
