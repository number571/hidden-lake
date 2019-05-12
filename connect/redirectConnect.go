package connect

import (
    "strings"
    "../models"
    "../settings"
)

func redirectConnect(connected_nodes map[string]string, connects []string) {
    var new_pack = models.PackageTCP {
        From: models.From {
            Name: settings.User.Hash.P2P,
        },
        Head: models.Head {
            Title: settings.HEAD_CONNECT,
            Mode: settings.MODE_GLOBAL,
        },
        Body: strings.Join(connects, settings.SEPARATOR),
    }

    for node := range connected_nodes {
        new_pack.To = node
        sendEncryptedPackage(new_pack, settings.P2P_mode)
    }
}
