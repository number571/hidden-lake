package connect

import (
    "net"
    "encoding/json"
    "../utils"
    "../models"
    "../crypto"
    "../settings"
)

func SendEncryptedPackage(pack settings.PackageTCP) int8 {
    if settings.User.NodeConnection[pack.To] != 1 {
        nullNode(pack.To)
        return settings.EXIT_FAILED
    }

    var new_pack = settings.PackageTCP {
        From: models.From {
            Name: pack.From.Name,
        },
        To: crypto.Encrypt(settings.User.NodeSessionKey[pack.To], pack.To),
        Head: models.Head {
            Header: crypto.Encrypt(settings.User.NodeSessionKey[pack.To], pack.Header),
            Mode: crypto.Encrypt(settings.User.NodeSessionKey[pack.To], pack.Mode),
        },
        Body: crypto.Encrypt(settings.User.NodeSessionKey[pack.To], pack.Body),
    }

    return sendNodePackage(pack.To, new_pack)
}

func sendNodePackage(to string, pack settings.PackageTCP) int8 {
    conn, err := net.Dial(settings.PROTOCOL_TCP, settings.User.NodeAddress[to])
    if err != nil {
        nullNode(to)
        return settings.EXIT_FAILED
    }

    data, err := json.Marshal(pack)
    utils.CheckError(err)

    conn.Write(data)
    conn.Close()

    return settings.EXIT_SUCCESS
}

func sendAddrPackage(to string, pack settings.PackageTCP) int8 {
    conn, err := net.Dial(settings.PROTOCOL_TCP, to)
    if err != nil {
        return settings.EXIT_FAILED
    }

    data, err := json.Marshal(pack)
    utils.CheckError(err)

    conn.Write(data)
    conn.Close()

    return settings.EXIT_SUCCESS
}
