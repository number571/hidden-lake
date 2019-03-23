package connect

import (
    "net"
    "fmt"
    "time"
    "encoding/hex"
    "encoding/json"
    "../utils"
    "../crypto"
    "../models"
    "../settings"
)

func encrypt(session_key []byte, data string) string {
    result, _ := crypto.EncryptAES(
        []byte(data),
        session_key,
    )
    return hex.EncodeToString(result)
}

func decrypt(session_key []byte, data string) string {
    decoded, _ := hex.DecodeString(data)
    result, _ := crypto.DecryptAES(
        decoded,
        session_key,
    )
    return string(result)
}

func nullNode(username string) {
    settings.Mutex.Lock()
    settings.User.NodeAddress[username] = ""
    settings.User.NodeConnection[username] = 0
    settings.User.NodePublicKey[username]  = nil
    settings.User.NodeSessionKey[username] = nil
    settings.User.Connections = utils.RemoveByElem(
        settings.User.Connections,
        username,
    )
    settings.Mutex.Unlock()
}

func findConnects(seconds time.Duration) {
    for {
        for _, username := range settings.User.Connections {
            var new_pack = settings.Package {
                From: models.From {
                    Address: settings.User.IPv4 + settings.User.Port,
                    Name: settings.User.Name,
                },
                To: username,
                Head: models.Head {
                    Header: settings.HEAD_CONNECT,
                    Mode: settings.MODE_READ_LIST,
                },
            }
            SendEncryptedPackage(new_pack)
        }
        time.Sleep(seconds * time.Second)
    }
}

func printGlobalHistory() {
    for _, mes := range settings.User.GlobalMessages {
        fmt.Println("|", mes)
    }
}

func printLocalHistory(slice []string) {
    for _, check := range slice {
        for _, username := range settings.User.Connections {
            if username == check {
                fmt.Printf("| %s:\n", username)
                for _, mes := range settings.User.LocalMessages[username] {
                    fmt.Println("|", mes)
                }
                break
            }
        }
    }
}

func sendNodePackage(to string, pack settings.Package) int8 {
    conn, err := net.Dial(settings.PROTOCOL_TCP, settings.User.NodeAddress[to])
    if err != nil {
        nullNode(to)
        return settings.EXIT_FAILED
    }

    data, err := json.Marshal(pack)
    utils.CheckError(err)

    conn.Write(data)
    conn.Close()

    return settings.EXIT_SUCCESS
}

func sendAddrPackage(to string, pack settings.Package) int8 {
    conn, err := net.Dial(settings.PROTOCOL_TCP, to)
    if err != nil {
        return settings.EXIT_FAILED
    }

    data, err := json.Marshal(pack)
    utils.CheckError(err)

    conn.Write(data)
    conn.Close()

    return settings.EXIT_SUCCESS
}
