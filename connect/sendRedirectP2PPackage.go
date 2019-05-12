package connect

import (
    "strings"
    "encoding/hex"
    "../utils"
    "../crypto"
    "../models"
    "../settings"
)

func sendRedirectP2PPackage(pack models.PackageTCP) {
    var data = strings.Split(pack.From.Address, settings.SEPARATOR)

    decoded, err := hex.DecodeString(data[0])
    utils.CheckError(err)

    session_key, err := crypto.DecryptRSA(
        []byte(decoded),
        settings.User.Private.Key.P2P,
    )
    utils.CheckError(err)

    var result = crypto.Decrypt(session_key, data[1])
    var addresses = strings.Split(result, settings.SEPARATOR)

    var new_pack = models.PackageTCP {
        From: models.From {
            Name: settings.User.Hash.P2P,
            Address: strings.Join(addresses[1:], settings.SEPARATOR),
        },
        To: addresses[0],
        Head: models.Head {
            Title: pack.Head.Title,
            Mode: pack.Head.Mode,
        },
        Body: pack.Body,
    }

    sendEncryptedPackage(new_pack, settings.P2P_mode)
}
