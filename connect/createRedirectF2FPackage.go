package connect

import (
    "strings"
    "../crypto"
    "../models"
    "../settings"
)

// Initialize redirect-package for F2F.
func createRedirectF2FPackage(pack *models.PackageTCP, to string) {
    var (
        connects = append(
            settings.MakeConnects(settings.Node.SessionKey.F2F), 
            settings.User.Hash.F2F,
        )
        session_message = string(crypto.SessionKey(8))
    )
    *pack = models.PackageTCP {
        From: models.From {
            Hash: settings.User.Hash.F2F,
            Address: pack.From.Address,
        },
        To: models.To {
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
