package connect

import (
    "time"
    "encoding/hex"
    "../models"
    "../settings"
)

func connectP2P(slice []string, check bool) {
next:
    for _, addr := range slice {
        var address = settings.User.IPv4 + settings.User.Port

        if addr == address {
            continue
        }

        for _, node_addr := range settings.Node.Address.P2P {
            if addr == node_addr {
                continue next
            }
        }

        var new_pack = models.PackageTCP {
            From: models.From {
                Address: address,
                Name: settings.User.Hash.P2P,
            },
            Head: models.Head {
                Title: settings.HEAD_CONNECT,
                Mode: settings.MODE_READ_LOCAL,
            },
        }

        if !check {
            new_pack.From.Login = settings.User.Login
            new_pack.Head.Mode = settings.MODE_READ
            new_pack.Body = hex.EncodeToString([]byte(settings.User.Public.Data.P2P))
        }

        sendAddrPackage(addr, new_pack)
        time.Sleep(time.Millisecond * 500)
    }
}
