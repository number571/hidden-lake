package connect

import (
    "strings"
    "../models"
    "../settings"
)

func sendRedirectF2FPackage(pack models.PackageTCP) {
    var (
        addresses = strings.Split(pack.From.Address, settings.SEPARATOR)
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
                Name: settings.User.Hash.F2F,
                Login: pack.From.Name,
                Address: strings.Join(addresses, settings.SEPARATOR),
            },
            To: address,
            Head: models.Head {
                Title: pack.Head.Title,
                Mode: pack.Head.Mode,
            },
            Body: pack.Body,
        }
        sendEncryptedPackage(new_pack, settings.F2F_mode)
    }
}
