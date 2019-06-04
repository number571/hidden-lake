package connect

import (
    "../crypto"
    "../models"
    "../settings"
)

// Send package with encrypt data by session key.
func SendEncryptedPackage(pack models.PackageTCP, mode models.ModeNet) int8 {
    // settings.User.Mode != models.F2F_mode && 
    if settings.Node.ConnectionMode[pack.To.Hash] != models.CONN {
        nullNode(pack.To.Hash)
        return settings.EXIT_FAILED
    }

    var session_key []byte 
    switch mode {
        case models.P2P_mode, models.C_S_mode:
            if _, ok := settings.Node.SessionKey.P2P[pack.To.Hash]; !ok {
                return settings.EXIT_FAILED
            }
            session_key = settings.Node.SessionKey.P2P[pack.To.Hash]
        case models.F2F_mode:
            if _, ok := settings.Node.SessionKey.F2F[pack.To.Hash]; !ok {
                return settings.EXIT_FAILED
            }
            session_key = settings.Node.SessionKey.F2F[pack.To.Hash]
    }

    var new_pack = models.PackageTCP {
        From: models.From {
            Hash: pack.From.Hash,
            Address: crypto.Encrypt(session_key, pack.From.Address),
        },
        To: models.To {
            Hash: crypto.Encrypt(session_key, pack.To.Hash),
            Address: crypto.Encrypt(session_key, pack.To.Address),
        },
        Head: models.Head {
            Title: crypto.Encrypt(session_key, pack.Head.Title),
            Mode: crypto.Encrypt(session_key, pack.Head.Mode),
        },
        Body: crypto.Encrypt(session_key, pack.Body),
    }

    if mode == models.C_S_mode {
        if settings.User.Mode == models.C_S_mode {
            return sendPackageByArchCS(settings.Node.ConnServer.Addr, new_pack)
        }
        return sendPackageByArchCS(settings.Node.Address.C_S[pack.To.Hash], new_pack)
    } 
    return sendPackageByNode(pack.To.Hash, new_pack, mode)
}
