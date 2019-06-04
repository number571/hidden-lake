package connect

import (
    // "../crypto"
	"../models"
	"../settings"
)

func createRedirectArchCSPackage(pack *models.PackageTCP) {
    // var body = pack.Body
    // if pack.To.Hash != "" { 
    //     body = crypto.Encrypt(settings.Node.SessionKey.P2P[pack.To.Hash], pack.Body)
    // } 
    // else if settings.Node.ConnServer.Hash != "" {
    //     body = crypto.Encrypt(settings.Node.SessionKey.P2P[settings.Node.ConnServer.Hash], pack.Body)
    // }

	*pack = models.PackageTCP {
		From: models.From {
			Address: settings.User.Hash.P2P,
            Hash: settings.User.Hash.P2P,
        },
        To: models.To {
        	Address: pack.To.Hash,
            Hash: settings.Node.ConnServer.Hash,
        },
        Head: models.Head {
            Title: settings.HEAD_REDIRECT,
            Mode: pack.Head.Title + settings.SEPARATOR + pack.Head.Mode,
        }, 
        Body: pack.Body,
	}
}
