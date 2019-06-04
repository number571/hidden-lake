package connect

import (
    "net"
    "encoding/hex"
    "../models"
    "../settings"
)

// Send connect-package to node.
func ConnectArchCS(conn net.Conn, check bool) {
    if conn == nil { return }
    var new_pack = models.PackageTCP {
        From: models.From {
            Hash: settings.User.Hash.P2P,
        },
        Head: models.Head {
            Title: settings.HEAD_CONNECT,
            Mode: settings.MODE_READ_LOCAL,
        },
    }

    if !check {
        new_pack.Head.Mode = settings.MODE_READ
        new_pack.Body = hex.EncodeToString([]byte(settings.User.Public.Data.P2P))
    }

    createRedirectArchCSPackage(&new_pack)
    if sendPackageByArchCS(conn, new_pack) == settings.EXIT_SUCCESS {
        go readConnection(conn)
    }
}
