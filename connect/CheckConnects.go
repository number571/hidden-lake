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

// Check node connection in real time.
func CheckConnects() {
    for {
        if !settings.User.Auth { break }
        for name, conn := range settings.Node.Address.C_S {
            if conn == nil { 
                delete(settings.Node.Address.C_S, name)
                delete(settings.Node.SessionKey.P2P, name)
            }
        }
        for username := range settings.Node.Address.P2P {
            go checkUser(username, SECONDS_OF_TIMER)
            time.Sleep(time.Second * 1)
        }
        time.Sleep(time.Second * SECONDS_OF_WAITING)
    }
}

// Send test-package to node.
func checkUser(username string, timer uint) {
    var new_pack = models.PackageTCP {
        From: models.From {
            Hash: settings.User.Hash.P2P,
        },
        Head: models.Head {
            Title: settings.HEAD_CONNECT,
            Mode: settings.MODE_READ_CHECK,
        },
    }

    new_pack.To.Hash = username
    SendEncryptedPackage(new_pack, models.P2P_mode)
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
