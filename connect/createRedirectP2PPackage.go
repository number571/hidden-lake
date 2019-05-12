package connect

import (
    "encoding/hex"
    "../utils"
    "../crypto"
    "../models"
    "../settings"
)

func createRedirectP2PPackage(pack *models.PackageTCP) {
    encrypted_hashname, err := crypto.EncryptRSA(
        []byte(settings.User.Hash.P2P),
        settings.Node.PublicKey[pack.To],
    )
    utils.CheckError(err)
    *pack = models.PackageTCP {
        From: models.From {
            Name: pack.From.Name,
            Address: onionOverlay(pack.To, settings.QUAN_OF_ROUTING_NODES),
        },
        Head: models.Head {
            Title: settings.HEAD_REDIRECT,
            Mode: hex.EncodeToString(encrypted_hashname) + settings.SEPARATOR +
                crypto.Encrypt(
                    settings.Node.SessionKey.P2P[pack.To], 
                    pack.Head.Title + settings.SEPARATOR + pack.Head.Mode,
                ),
        },
        Body: crypto.Encrypt(settings.Node.SessionKey.P2P[pack.To], pack.Body),
    }
}

func onionOverlay(to string, quan uint8) string {
    var (
        list []string
        result string
    )

    if settings.DYNAMIC_ROUTING {
        list = append(list, to)
    }

    for node := range settings.Node.Address.P2P {
        if node == to { continue }
        if quan == 0 { break }
        list = append(list, node)
        quan--
    }

    utils.Shuffle(list)

    if !settings.DYNAMIC_ROUTING {
        list = append(list, to)
    }

    for i := len(list)-1; i > 0; i-- {
        var session_key = crypto.SessionKey(settings.ROUTING_KEY_BYTES)
        encrypted_key, err := crypto.EncryptRSA(
            session_key,
            settings.Node.PublicKey[list[i-1]],
        )
        utils.CheckError(err)

        result = hex.EncodeToString(encrypted_key) + settings.SEPARATOR + 
            crypto.Encrypt(session_key, list[i] + settings.SEPARATOR + result)
    }

    result = list[0] + settings.SEPARATOR + result
    return result
}
