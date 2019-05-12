package connect

import (
    "../crypto"
    "../models"
    "../settings"
)

func sendEncryptedPackage(pack models.PackageTCP, mode settings.ModeNet) int8 {
    if !settings.User.ModeF2F && settings.Node.Connection[pack.To] != 1 {
        nullNode(pack.To)
        return settings.EXIT_FAILED
    }

    var session_key []byte 
    switch mode {
        case settings.P2P_mode: session_key = settings.Node.SessionKey.P2P[pack.To]
        case settings.F2F_mode:
            if _, ok := settings.Node.SessionKey.F2F[pack.To]; !ok {
                return settings.EXIT_FAILED
            }
            session_key = settings.Node.SessionKey.F2F[pack.To]
    }

    var new_pack = models.PackageTCP {
        From: models.From {
            Name: pack.From.Name,
            Login: crypto.Encrypt(session_key, pack.From.Login),
            Address: crypto.Encrypt(session_key, pack.From.Address),
        },
        To: crypto.Encrypt(session_key, pack.To),
        Head: models.Head {
            Title: crypto.Encrypt(session_key, pack.Head.Title),
            Mode: crypto.Encrypt(session_key, pack.Head.Mode),
        },
        Body: crypto.Encrypt(session_key, pack.Body),
    }

    return sendNodePackage(pack.To, new_pack, mode)
}
