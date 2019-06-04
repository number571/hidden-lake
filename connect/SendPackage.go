package connect

import (
    "../models"
    "../settings"
)

// Interface for send package.
func SendPackage(pack models.PackageTCP, mode models.ModeNet) {
    switch mode {
        case models.P2P_mode: 
            // If P2P node is exists, then create and send redirect package.
            if _, ok := settings.Node.Address.P2P[pack.To.Hash]; ok {
                createRedirectP2PPackage(&pack)
                sendInitRedirectP2PPackage(pack)
                return
            } else if _, ok := settings.Node.Address.C_S[pack.To.Hash]; ok {
                SendEncryptedPackage(pack, models.C_S_mode)
                return
            }
            // If P2P node is not specified, then send package 
            // to all nodes and clients.
            if pack.To.Hash == "" {
                for username := range settings.Node.Address.P2P {
                    pack.To.Hash = username
                    SendEncryptedPackage(pack, models.P2P_mode)
                }
                for username := range settings.Node.Address.C_S {
                    // If is a redirect-package by client, then pass.
                    if username == pack.From.Address { continue }
                    pack.To.Hash = username
                    SendEncryptedPackage(pack, models.C_S_mode)
                }
            }
        case models.F2F_mode:
            // If P2P node is exists, then send package.
            if _, ok := settings.Node.Address.F2F[pack.To.Hash]; ok {
                SendEncryptedPackage(pack, models.F2F_mode)
                return
            }
            // If P2P node is not specified, then create and 
            // send redirect package to all F2F nodes.
            createRedirectF2FPackage(&pack, pack.To.Hash)
            for username := range settings.Node.Address.F2F {
                pack.To.Hash = username
                SendEncryptedPackage(pack, models.F2F_mode)
            }
        case models.C_S_mode:
            // If sender is a client and he have connection with a server, then
            // he create and send redirect-package throw server.
            if settings.User.Mode == models.C_S_mode {
                createRedirectArchCSPackage(&pack)
            // Send package to all clients by P2P node.
            } else if pack.To.Hash == "" {
                for username := range settings.Node.Address.C_S {
                    pack.To.Hash = username
                    SendEncryptedPackage(pack, models.C_S_mode)
                }
                return
            }
            // Just send a package to one client (from server) 
            // or to server (from client).
            SendEncryptedPackage(pack, models.C_S_mode)
    }
}
