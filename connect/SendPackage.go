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

func SendPackage(pack settings.PackageTCP, is_f2f bool) {
    if is_f2f {
        sendEncryptedPackage(pack, true)
    } else {
        createRedirectPackage(&pack)
        sendInitRedirectPackage(pack)
    }
}

func CreateRedirectF2FPackage(pack *settings.PackageTCP, to string) {
    var (
        addresses = make([]string, len(settings.User.NodeAddressF2F) + 1)
        session_message = string(crypto.SessionKey(8))
        index uint
    )
    for username := range settings.User.NodeAddressF2F {
        addresses[index] = username
        index++
    }
    addresses[index] = settings.User.Hash
    *pack = settings.PackageTCP {
        From: models.From {
            Name: settings.User.Hash,
            Login: pack.From.Name,
            Address: strings.Join(addresses, settings.SEPARATOR),
        },
        Head: models.Head {
            Header: settings.HEAD_REDIRECT,
            Mode:   to + settings.SEPARATOR + pack.Head.Header + 
                    settings.SEPARATOR + pack.Head.Mode + settings.SEPARATOR +
                    session_message,
        },
        Body: pack.Body,
    }
}

func sendRedirectF2FPackage(pack settings.PackageTCP) {
    var (
        addresses = strings.Split(pack.From.Address, settings.SEPARATOR)
        to_addresses []string
    )
    for username := range settings.User.NodeAddressF2F {
        var flag bool
        for _, address := range addresses {
            if address == username {
                flag = true
                break
            }
        }
        if !flag {
            addresses = append(addresses, username)
            to_addresses = append(to_addresses, username)
        }
    }
    for _, address := range to_addresses {
        var new_pack = settings.PackageTCP {
            From: models.From {
                Name: settings.User.Hash,
                Login: pack.From.Name,
                Address: strings.Join(addresses, settings.SEPARATOR),
            },
            To: address,
            Head: models.Head {
                Header: pack.Head.Header,
                Mode: pack.Head.Mode,
            },
            Body: pack.Body,
        }
        sendEncryptedPackage(new_pack, true)
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

    for node := range settings.User.NodeAddress {
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
            settings.User.NodePublicKey[list[i-1]],
        )
        utils.CheckError(err)

        result = hex.EncodeToString(encrypted_key) + settings.SEPARATOR + 
            crypto.Encrypt(session_key, list[i] + settings.SEPARATOR + result)
    }

    result = list[0] + settings.SEPARATOR + result
    return result
}

func createRedirectPackage(pack *settings.PackageTCP) {
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

func sendInitRedirectPackage(pack settings.PackageTCP) {
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
    sendEncryptedPackage(new_pack, false)
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

    sendEncryptedPackage(new_pack, false)
}

func sendEncryptedPackage(pack settings.PackageTCP, is_f2f bool) int8 {
    if !settings.User.ModeF2F && settings.User.NodeConnection[pack.To] != 1 {
        nullNode(pack.To)
        return settings.EXIT_FAILED
    }

    var session_key []byte 
    if is_f2f {
        if _, ok := settings.User.NodeSessionKeyF2F[pack.To]; !ok {
            return settings.EXIT_FAILED
        }
        session_key = settings.User.NodeSessionKeyF2F[pack.To]
    } else {
        session_key = settings.User.NodeSessionKey[pack.To]
    }

    var new_pack = settings.PackageTCP {
        From: models.From {
            Name: pack.From.Name,
            Login: crypto.Encrypt(session_key, pack.From.Login),
            Address: crypto.Encrypt(session_key, pack.From.Address),
        },
        To: crypto.Encrypt(session_key, pack.To),
        Head: models.Head {
            Header: crypto.Encrypt(session_key, pack.Header),
            Mode: crypto.Encrypt(session_key, pack.Mode),
        },
        Body: crypto.Encrypt(session_key, pack.Body),
    }

    return sendNodePackage(pack.To, new_pack, is_f2f)
}

func sendNodePackage(to string, pack settings.PackageTCP, is_f2f bool) int8 {
    var address string
    if is_f2f {
        address = settings.User.NodeAddressF2F[to]
    } else {
        address = settings.User.NodeAddress[to]
    }

    conn, err := net.Dial(settings.PROTOCOL_TCP, address)
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
