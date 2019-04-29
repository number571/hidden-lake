package connect

import (
    "time"
    "encoding/hex"
    "../models"
    "../settings"
)

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
