package connect

import (
    "strings"
    "../models"
    "../settings"
)

// Send created redirect-package to P2P nodes.
func sendInitRedirectP2PPackage(pack models.PackageTCP) {
    var addresses = strings.Split(pack.To.Address, settings.SEPARATOR)
    var new_pack = models.PackageTCP {
        From: models.From {
            Hash: settings.User.Hash.P2P,
            Address: pack.From.Address,
        },
        To: models.To {
            Hash: addresses[0],
            Address: strings.Join(addresses[1:], settings.SEPARATOR),
        },
        Head: models.Head {
            Title: pack.Head.Title,
            Mode: pack.Head.Mode,
        },
        Body: pack.Body,
    }
    SendEncryptedPackage(new_pack, models.P2P_mode)
}
