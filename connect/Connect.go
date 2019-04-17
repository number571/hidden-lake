package connect

import (
    "time"
    "encoding/hex"
    "../models"
    "../settings"
)

func Connect(slice []string) {
    next:
    for _, addr := range slice {
        var address = settings.User.IPv4 + settings.User.Port

        if addr == address {
            continue
        }

        for _, username := range settings.User.Connections {
            if addr == settings.User.NodeAddress[username] {
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
                Mode: settings.MODE_READ,
            },
            Body: hex.EncodeToString([]byte(settings.User.PublicData)),
        }

        sendAddrPackage(addr, new_pack)
        time.Sleep(time.Millisecond * 500)
    }
}
