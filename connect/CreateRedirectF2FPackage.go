package connect

import (
    "strings"
    "../crypto"
    "../models"
    "../settings"
)

func CreateRedirectF2FPackage(pack *models.PackageTCP, to string) {
    var (
        connects = append(
            settings.MakeConnects(settings.Node.Address.F2F), 
            settings.User.Hash.F2F,
        )
        session_message = string(crypto.SessionKey(8))
    )
    *pack = models.PackageTCP {
        From: models.From {
            Name: settings.User.Hash.F2F,
            Login: pack.From.Name,
            Address: strings.Join(connects, settings.SEPARATOR),
        },
        Head: models.Head {
            Title: settings.HEAD_REDIRECT,
            Mode:   to + settings.SEPARATOR + pack.Head.Title + 
                    settings.SEPARATOR + pack.Head.Mode + settings.SEPARATOR +
                    session_message,
        },
        Body: pack.Body,
    }
}
