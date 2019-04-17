package connect

import (
    "../models"
    "../settings"
)

func Disconnect(slice []string) {
    for _, addr := range slice {
        var new_pack = settings.PackageTCP {
            From: models.From {
                Address: settings.User.IPv4 + settings.User.Port,
                Name: settings.User.Hash,
            },
            To: addr,
            Head: models.Head {
                Header: settings.HEAD_WARNING,
                Mode: settings.MODE_SAVE,
            },
        }
        SendEncryptedPackage(new_pack)
        nullNode(addr)
    }
}
