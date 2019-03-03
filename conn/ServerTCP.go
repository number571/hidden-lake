package conn

import (
    "fmt"
    "net"
    "strings"
    "io/ioutil"
    "encoding/hex"
    "encoding/json"
    "../utils"
    "../models"
    "../crypto"
    "../settings"
    "../encoding"
)

func decrypt(from, data string) string {
    decoded, _ := hex.DecodeString(data)
    result, _ := crypto.DecryptAES(
        decoded,
        settings.User.NodeSessionKey[from],
    )
    return string(result)
}

func nullNode(addr string) {
    settings.Mutex.Lock()
    settings.User.NodePublicKey[addr] = nil
    settings.User.NodeSessionKey[addr] = nil
    settings.User.NodeConnection[addr] = 0
    settings.User.Connections = utils.RemoveByElem(
        settings.User.Connections,
        addr,
    )
    settings.Mutex.Unlock()
}

func ServerTCP() {
    listen, err := net.Listen(settings.PROTOCOL_TCP, settings.IPV4_TEMPLATE + settings.User.Port) 
    utils.CheckError(err)
    defer listen.Close()

    var buffer = make([]byte, settings.BUFF_SIZE)

    for {
        conn, err := listen.Accept()
        if err != nil { 
            fmt.Println(err)
            continue 
        }

        var message = ""
        for {
            length, err := conn.Read(buffer)
            if length == 0 || err != nil { break }
            message += string(buffer[:length])
        }

        // fmt.Println(message)

        var pack settings.Package
        err = json.Unmarshal([]byte(message), &pack)
        if err != nil {
            fmt.Println(err)
            continue
        }

        if settings.User.NodeConnection[pack.From.Address] == 1 &&
           pack.Head.Header != settings.HEAD_CONNECT {
            pack = settings.Package {
                From: models.From {
                    Address: pack.From.Address,
                    Name: decrypt(pack.From.Address, pack.From.Name),
                },
                To: decrypt(pack.From.Address, pack.To),
                Head: models.Head {
                    Header: decrypt(pack.From.Address, pack.Header),
                    Mode: decrypt(pack.From.Address, pack.Mode),
                }, 
                Body: decrypt(pack.From.Address, pack.Body),
            }
        }

        switch pack.Header {
            case settings.HEAD_ARCHIVE: 
                switch pack.Mode {
                    case settings.MODE_GET_LIST: 
                        files, err := ioutil.ReadDir(settings.PATH_ARCHIVE)
                        utils.CheckError(err)
                        var files_str = ""
                        for _, file := range files {
                            files_str += file.Name() + settings.SEPARATOR
                        }
                        var new_pack = settings.Package {
                            From: models.From {
                                Address: pack.To,
                            },
                            To: pack.From.Address,
                            Head: models.Head {
                                Header: settings.HEAD_ARCHIVE,
                                Mode: settings.MODE_SAVE_LIST,
                            },
                            Body: files_str,
                        }
                        SendEncryptedPackage(new_pack)

                    case settings.MODE_SAVE_LIST: 
                        settings.Mutex.Lock()
                        settings.User.TempArchive = strings.Split(pack.Body, settings.SEPARATOR)
                        settings.Mutex.Unlock()

                    case settings.MODE_GET_FILE:
                        if !utils.FileIsExist(settings.PATH_ARCHIVE + pack.Body) {
                            goto close_connection
                        }
                        var new_pack = settings.Package {
                            From: models.From {
                                Address: pack.To,
                            },
                            To: pack.From.Address,
                            Head: models.Head {
                                Header: settings.HEAD_ARCHIVE,
                                Mode: settings.MODE_SAVE_FILE,
                            },
                            Body: pack.Body + settings.SEPARATOR + utils.ReadFile(settings.PATH_ARCHIVE + pack.Body),
                        }
                        SendEncryptedPackage(new_pack)

                    case settings.MODE_SAVE_FILE: 
                        var splited = strings.Split(pack.Body, settings.SEPARATOR)
                        var add = ""
                        if utils.FileIsExist(settings.PATH_ARCHIVE + splited[0]) {
                            add += "copy_"
                        }
                        utils.WriteFile(settings.PATH_ARCHIVE + add + splited[0], splited[1])
                }

            case settings.HEAD_PROFILE: 
                switch pack.Mode {
                    case settings.MODE_GET:
                        var new_pack = settings.Package {
                            From: models.From {
                                Address: pack.To,
                            },
                            To: pack.From.Address,
                            Head: models.Head {
                                Header: settings.HEAD_PROFILE,
                                Mode: settings.MODE_SAVE,
                            },
                            Body:   settings.User.Name + settings.SEPARATOR +
                                    settings.User.IPv4 + settings.SEPARATOR +
                                    settings.User.Port + settings.SEPARATOR +
                                    settings.User.Info,
                        }
                        SendEncryptedPackage(new_pack)

                    case settings.MODE_SAVE: 
                        settings.Mutex.Lock()
                        settings.User.TempProfile = strings.Split(pack.Body, settings.SEPARATOR)
                        settings.Mutex.Unlock()
                }

            case settings.HEAD_MESSAGE: 
                var message = fmt.Sprintf("[%s/%s]: %s\n", pack.From.Address, pack.From.Name, pack.Body)
                // var message = fmt.Sprintf("[%s/%s]: %s\n", pack.From.Address, pack.From.Name, pack.Body)
                fmt.Print(message)
                switch pack.Mode {
                    case settings.MODE_LOCAL:  
                        settings.Mutex.Lock()
                        settings.User.LocalMessages[pack.From.Address] = append(
                            settings.User.LocalMessages[pack.From.Address],
                            message,
                        ) 
                        settings.Mutex.Unlock()

                    case settings.MODE_GLOBAL: 
                        settings.Mutex.Lock()
                        settings.User.GlobalMessages = append(
                            settings.User.GlobalMessages,
                            message,
                        )
                        settings.Mutex.Unlock()
                }

            case settings.HEAD_CONNECT:
                switch pack.Mode {
                    case settings.MODE_GET: 
                        if settings.User.NodeConnection[pack.From.Address] != 1 {
                            decoded, err := hex.DecodeString(pack.Body)
                            utils.CheckError(err)

                            settings.Mutex.Lock()
                            settings.User.NodePublicKey[pack.From.Address] = encoding.DecodePublic(string(decoded))
                            settings.User.NodeSessionKey[pack.From.Address] = crypto.SessionKey(32)
                            settings.Mutex.Unlock()
                        } 

                        encrypted, err := crypto.EncryptRSA(
                            settings.User.NodePublicKey[pack.From.Address],
                            settings.User.NodeSessionKey[pack.From.Address],
                        )
                        utils.CheckError(err)

                        var new_pack = settings.Package {
                            From: models.From {
                                Address: settings.User.IPv4 + settings.User.Port,
                                Name: settings.User.Name,
                            },
                            To: pack.From.Address,
                            Head: models.Head {
                                Header: settings.HEAD_CONNECT,
                                Mode: settings.MODE_SAVE,
                            },
                            Body: hex.EncodeToString(encrypted) + settings.SEPARATOR + hex.EncodeToString([]byte(settings.User.PublicData)),
                        }

                        SendPackage(pack.From.Address, new_pack)

                        if settings.User.NodeConnection[pack.From.Address] == -1 {
                            settings.Mutex.Lock()
                            settings.User.Connections = append(
                                settings.User.Connections, 
                                pack.From.Address,
                            )
                            settings.User.NodeConnection[pack.From.Address] = 1
                            settings.Mutex.Unlock()
                        }

                    case settings.MODE_SAVE: 
                        // if settings.User.NodeConnection[pack.From.Address] == 1 {
                        //     continue
                        // }

                        var splited = strings.Split(pack.Body, settings.SEPARATOR)

                        encrypted_key, err := hex.DecodeString(splited[0])
                        utils.CheckError(err)

                        decrypted, err := crypto.DecryptRSA(
                            settings.User.PrivateKey,
                            encrypted_key,
                        )
                        utils.CheckError(err)

                        public_key, err := hex.DecodeString(splited[1])
                        utils.CheckError(err)

                        settings.Mutex.Lock()
                        settings.User.NodePublicKey[pack.From.Address] = encoding.DecodePublic(string(public_key))
                        settings.User.NodeSessionKey[pack.From.Address] = decrypted
                        settings.User.NodeConnection[pack.From.Address] = 1
                        settings.Mutex.Unlock()
                }

            case settings.HEAD_WARNING:
                switch pack.Mode {
                    case settings.MODE_SAVE: 
                        nullNode(pack.From.Address)
                        fmt.Printf("[DISCONNECTED]: %s\n", pack.From.Address)
                }

            default:
                var new_pack = settings.Package {
                    From: models.From {
                        Address: settings.User.IPv4 + settings.User.Port,
                        Name: settings.User.Name,
                    },
                    To: pack.From.Address,
                    Head: models.Head {
                        Header: settings.HEAD_WARNING,
                        Mode: settings.MODE_SAVE,
                    },
                }
                SendPackage(pack.From.Address, new_pack)
        }

close_connection:
        conn.Close()
    }
}
