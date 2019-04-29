package connect

import (
    "net"
    "strings"
    "encoding/hex"
    "encoding/json"
    "../utils"
    "../models"
    "../crypto"
    "../settings"
)

func onionOverlay(to string, quan uint8) string {
    var (
        list []string
        result string
    )

    for node := range settings.User.NodeAddress {
        if node == to { continue }
        if quan == 0 { break }
        list = append(list, node)
        quan--
    }

    utils.Shuffle(list)
    list = append(list, to)

    for i := len(list)-1; i > 0; i-- {
        var session_key = crypto.SessionKey(settings.ROUTING_KEY_BYTES)
        encrypted_key, err := crypto.EncryptRSA(
            session_key,
            settings.User.NodePublicKey[list[i-1]],
        )
        utils.CheckError(err)

        result = hex.EncodeToString(encrypted_key) + settings.SEPARATOR + 
            crypto.Encrypt(session_key, list[i] + settings.SEPARATOR + result)
    }

    result = list[0] + settings.SEPARATOR + result
    return result
}

func CreateRedirectPackage(pack *settings.PackageTCP) {
    encrypted_hashname, err := crypto.EncryptRSA(
        []byte(settings.User.Hash),
        settings.User.NodePublicKey[pack.To],
    )
    utils.CheckError(err)
    *pack = settings.PackageTCP {
        From: models.From {
            Name: pack.From.Name,
            Address: onionOverlay(pack.To, settings.QUAN_OF_ROUTING_NODES),
        },
        Head: models.Head {
            Header: settings.HEAD_REDIRECT,
            Mode: hex.EncodeToString(encrypted_hashname) + settings.SEPARATOR +
                crypto.Encrypt(
                    settings.User.NodeSessionKey[pack.To], 
                    pack.Head.Header + settings.SEPARATOR + pack.Head.Mode,
                ),
        },
        Body: crypto.Encrypt(settings.User.NodeSessionKey[pack.To], pack.Body),
    }
}

func SendInitRedirectPackage(pack settings.PackageTCP) {
    var addresses = strings.Split(pack.From.Address, settings.SEPARATOR)
    var new_pack = settings.PackageTCP {
        From: models.From {
            Name: settings.User.Hash,
            Address: strings.Join(addresses[1:], settings.SEPARATOR),
        },
        To: addresses[0],
        Head: models.Head {
            Header: pack.Head.Header,
            Mode: pack.Head.Mode,
        },
        Body: pack.Body,
    }
    sendEncryptedPackage(new_pack)
}

func sendRedirectPackage(pack settings.PackageTCP) {
    var data = strings.Split(pack.From.Address, settings.SEPARATOR)

    decoded, err := hex.DecodeString(data[0])
    utils.CheckError(err)

    session_key, err := crypto.DecryptRSA(
        []byte(decoded),
        settings.User.PrivateKey,
    )
    utils.CheckError(err)

    var result = crypto.Decrypt(session_key, data[1])
    var addresses = strings.Split(result, settings.SEPARATOR)

    var new_pack = settings.PackageTCP {
        From: models.From {
            Name: settings.User.Hash,
            Address: strings.Join(addresses[1:], settings.SEPARATOR),
        },
        To: addresses[0],
        Head: models.Head {
            Header: pack.Head.Header,
            Mode: pack.Head.Mode,
        },
        Body: pack.Body,
    }

    sendEncryptedPackage(new_pack)
}

func sendEncryptedPackage(pack settings.PackageTCP) int8 {
    if settings.User.NodeConnection[pack.To] != 1 {
        nullNode(pack.To)
        return settings.EXIT_FAILED
    }

    var new_pack = settings.PackageTCP {
        From: models.From {
            Name: pack.From.Name,
            Address: crypto.Encrypt(settings.User.NodeSessionKey[pack.To], pack.From.Address),
        },
        To: crypto.Encrypt(settings.User.NodeSessionKey[pack.To], pack.To),
        Head: models.Head {
            Header: crypto.Encrypt(settings.User.NodeSessionKey[pack.To], pack.Header),
            Mode: crypto.Encrypt(settings.User.NodeSessionKey[pack.To], pack.Mode),
        },
        Body: crypto.Encrypt(settings.User.NodeSessionKey[pack.To], pack.Body),
    }

    return sendNodePackage(pack.To, new_pack)
}

func sendNodePackage(to string, pack settings.PackageTCP) int8 {
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

func sendAddrPackage(to string, pack settings.PackageTCP) int8 {
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
