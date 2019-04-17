package connect

import (
    "time"
    "../models"
    "../settings"
)

func FindConnects(seconds time.Duration) {
    for {
        if !settings.User.Auth {
            break 
        }
        Connect(settings.User.DefaultConnections)
        for _, username := range settings.User.Connections {
            var new_pack = settings.PackageTCP {
                From: models.From {
                    Address: settings.User.IPv4 + settings.User.Port,
                    Name: settings.User.Hash,
                },
                To: username,
                Head: models.Head {
                    Header: settings.HEAD_CONNECT,
                    Mode: settings.MODE_READ_LIST,
                },
            }
            SendEncryptedPackage(new_pack)
        }
        time.Sleep(seconds * time.Second)
    }
}
