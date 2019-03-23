package connect

import (
    "fmt"
    "../models"
    "../settings"
)

func Disconnect(slice []string) {
    for _, addr := range slice {
        fmt.Println("|", addr)
        var new_pack = settings.Package {
            From: models.From {
                Address: settings.User.IPv4 + settings.User.Port,
                Name: settings.User.Name,
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