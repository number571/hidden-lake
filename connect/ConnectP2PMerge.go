package connect

import (
    "encoding/hex"
    "../models"
    "../settings"
)

// Send connect-package to node with merge.
func ConnectP2PMerge(addr string) {
    var address = settings.User.IPv4 + settings.User.Port

    if addr == address {
        return
    }

    for _, node_addr := range settings.Node.Address.P2P {
        if addr == node_addr {
            return
        }
    }

    var new_pack = models.PackageTCP {
        From: models.From {
            Address: address,
            Hash: settings.User.Hash.P2P,
        },
        Head: models.Head {
            Title: settings.HEAD_CONNECT,
            Mode: settings.MODE_READ_GLOBAL,
        },
        Body: hex.EncodeToString([]byte(settings.User.Public.Data.P2P)),
    }

    sendPackageByAddr(addr, new_pack)
}
