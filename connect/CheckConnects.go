package connect

import (
    "time"
    "../models"
    "../settings"
)

const (
    SECONDS_OF_WAITING = 15
    SECONDS_OF_TIMER = 10
)

func CheckConnects() {
    for {
        if !settings.User.Auth { break }
        for username := range settings.Node.Address.P2P {
            go checkUser(username, SECONDS_OF_TIMER)
            time.Sleep(time.Second * 1)
        }
        time.Sleep(time.Second * SECONDS_OF_WAITING)
    }
}

func checkUser(username string, timer uint) {
    var new_pack = models.PackageTCP {
        From: models.From {
            Name: settings.User.Hash.P2P,
        },
        Head: models.Head {
            Title: settings.HEAD_CONNECT,
            Mode: settings.MODE_READ_CHECK,
        },
    }

    new_pack.To = username
    sendEncryptedPackage(new_pack, settings.P2P_mode)
    __check_connection[username] = false
    
check_again:
    if !__check_connection[username] && timer > 0 {
        time.Sleep(time.Second * 1)
        timer--
        goto check_again
    }
    if !__check_connection[username] {
        nullNode(username)
    }
}
