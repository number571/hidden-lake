package connect

import (
    "strings"
    "../models"
    "../settings"
)

// Send redirect-package to friends.
func sendRedirectF2FPackage(pack models.PackageTCP) {
    var (
        addresses = strings.Split(pack.To.Address, settings.SEPARATOR)
        to_addresses []string
    )
    for username := range settings.Node.Address.F2F {
        var flag bool
        for _, address := range addresses {
            if address == username {
                flag = true
                break
            }
        }
        if !flag {
            addresses = append(addresses, username)
            to_addresses = append(to_addresses, username)
        }
    }
    for _, address := range to_addresses {
        var new_pack = models.PackageTCP {
            From: models.From {
                Hash: settings.User.Hash.F2F,
                Address: pack.From.Address,
            },
            To: models.To {
                Hash: address,
                Address: strings.Join(addresses, settings.SEPARATOR),
            },
            Head: models.Head {
                Title: pack.Head.Title,
                Mode: pack.Head.Mode,
            },
            Body: pack.Body,
        }
        SendEncryptedPackage(new_pack, models.F2F_mode)
    }
}
