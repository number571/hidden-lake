package connect

import (
    "strings"
    "../models"
    "../settings"
)

func sendInitRedirectP2PPackage(pack models.PackageTCP) {
    var addresses = strings.Split(pack.From.Address, settings.SEPARATOR)
    var new_pack = models.PackageTCP {
        From: models.From {
            Name: settings.User.Hash.P2P,
            Address: strings.Join(addresses[1:], settings.SEPARATOR),
        },
        To: addresses[0],
        Head: models.Head {
            Title: pack.Head.Title,
            Mode: pack.Head.Mode,
        },
        Body: pack.Body,
    }
    sendEncryptedPackage(new_pack, settings.P2P_mode)
}
