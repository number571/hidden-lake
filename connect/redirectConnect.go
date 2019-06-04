package connect

import (
    "strings"
    "../models"
    "../settings"
)

// Redirect package in merge mode.
func redirectConnect(connected_nodes map[string]string, connects []string) {
    var new_pack = models.PackageTCP {
        From: models.From {
            Hash: settings.User.Hash.P2P,
        },
        Head: models.Head {
            Title: settings.HEAD_CONNECT,
            Mode: settings.MODE_GLOBAL,
        },
        Body: strings.Join(connects, settings.SEPARATOR),
    }

    for node := range connected_nodes {
        new_pack.To.Hash = node
        SendEncryptedPackage(new_pack, models.P2P_mode)
    }
}
