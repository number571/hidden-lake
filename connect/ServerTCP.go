package connect

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

func ServerTCP() {
    var err error
    settings.ServerListenTCP, err = net.Listen(settings.PROTOCOL_TCP, settings.IPV4_TEMPLATE + settings.User.Port) 
    utils.CheckError(err)
    // defer settings.ServerListenTCP.Close()

    for {
        if server(settings.ServerListenTCP) == 1 {
            break
        }
    }
}

func server(listen net.Listener) int8 {
    var buffer = make([]byte, settings.BUFF_SIZE)

    conn, err := listen.Accept()
    if err != nil { 
        return 1
    }

    var message = ""
    for {
        length, err := conn.Read(buffer)
        if length == 0 || err != nil { break }
        message += string(buffer[:length])
    }

    var pack settings.PackageTCP
    err = json.Unmarshal([]byte(message), &pack)
    if err != nil {
        fmt.Println(err)
        return 0
    }

    if settings.User.NodeConnection[pack.From.Name] == 1 &&
       pack.Head.Header != settings.HEAD_CONNECT {
        pack = settings.PackageTCP {
            From: models.From {
                Address: pack.From.Address,
                Name: pack.From.Name,
            },
            To: crypto.Decrypt(settings.User.NodeSessionKey[pack.From.Name], pack.To),
            Head: models.Head {
                Header: crypto.Decrypt(settings.User.NodeSessionKey[pack.From.Name], pack.Header),
                Mode: crypto.Decrypt(settings.User.NodeSessionKey[pack.From.Name], pack.Mode),
            }, 
            Body: crypto.Decrypt(settings.User.NodeSessionKey[pack.From.Name], pack.Body),
        }
    }

    switch pack.Header {
        case settings.HEAD_ARCHIVE: 
            switch pack.Mode {
                case settings.MODE_READ_LIST: 
                    files, err := ioutil.ReadDir(settings.PATH_ARCHIVE)
                    utils.CheckError(err)
                    var files_str = ""
                    for _, file := range files {
                        files_str += file.Name() + settings.SEPARATOR
                    }
                    var new_pack = settings.PackageTCP {
                        From: models.From {
                            Name: pack.To,
                        },
                        To: pack.From.Name,
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

                case settings.MODE_READ_FILE:
                    if utils.FileIsExist(settings.PATH_ARCHIVE + pack.Body) {
                        var new_pack = settings.PackageTCP {
                            From: models.From {
                                Name: pack.To,
                            },
                            To: pack.From.Name,
                            Head: models.Head {
                                Header: settings.HEAD_ARCHIVE,
                                Mode: settings.MODE_SAVE_FILE,
                            },
                            Body: pack.Body + settings.SEPARATOR + utils.ReadFile(settings.PATH_ARCHIVE + pack.Body),
                        }
                        SendEncryptedPackage(new_pack)
                    }

                case settings.MODE_SAVE_FILE: 
                    var splited = strings.Split(pack.Body, settings.SEPARATOR)
                    var add = ""
                    if utils.FileIsExist(settings.PATH_ARCHIVE + splited[0]) {
                        add += "copy_"
                    }
                    utils.WriteFile(settings.PATH_ARCHIVE + add + splited[0], splited[1])
            }

        case settings.HEAD_MESSAGE: 
            var message = fmt.Sprintf("[%s]: %s\n", settings.User.NodeLogin[pack.From.Name], pack.Body)
            fmt.Print(message)

            switch pack.Mode {
                case settings.MODE_LOCAL:  
                    settings.Mutex.Lock()
                    _, err := settings.DataBase.Exec(
                        "INSERT INTO Local" + pack.From.Name + " (User, Body) VALUES ($1, $2)",
                        pack.From.Name, 
                        crypto.Encrypt(settings.User.Password, message),
                    )
                    settings.Messages.CurrentIdLocal[pack.From.Name]++
                    settings.Mutex.Unlock()
                    utils.CheckError(err)
                    go func() {
                        settings.Messages.NewDataExistLocal[pack.From.Name] <- true
                    }()

                case settings.MODE_GLOBAL: 
                    settings.Mutex.Lock()
                    _, err := settings.DataBase.Exec(
                        "INSERT INTO GlobalMessages (User, Body) VALUES ($1, $2)",
                        pack.From.Name, 
                        crypto.Encrypt(settings.User.Password, message),
                    )
                    settings.Messages.CurrentIdGlobal++
                    settings.Mutex.Unlock()
                    utils.CheckError(err)
                    go func() {
                        settings.Messages.NewDataExistGlobal <- true
                    }()
            }

        case settings.HEAD_CONNECT:
            switch pack.Mode {
                case settings.MODE_READ: 
                    if settings.User.NodeConnection[pack.From.Name] == 1 {
                        goto close_connection
                    }
                    
                    public_key, err := hex.DecodeString(pack.Body)
                    utils.CheckError(err)

                    var public_data = string(public_key)

                    settings.Mutex.Lock()
                    settings.User.NodeAddress[pack.From.Name] = pack.From.Address
                    settings.User.NodeLogin[pack.From.Name] = pack.From.Login
                    settings.User.NodePublicKey[pack.From.Name] = encoding.DecodePublic(string(public_key))
                    settings.User.NodeSessionKey[pack.From.Name] = crypto.SessionKey(32)
                    settings.Mutex.Unlock()

                    var encrypted_address = crypto.Encrypt(settings.User.NodeSessionKey[pack.From.Name], settings.User.IPv4 + settings.User.Port)
                    var encrypted_login = crypto.Encrypt(settings.User.NodeSessionKey[pack.From.Name], settings.User.Login)
                    var encrypted_name = crypto.Encrypt(settings.User.NodeSessionKey[pack.From.Name], settings.User.Hash)

                    var connections = strings.Join(settings.User.Connections, settings.SEPARATOR_ADDRESS)
                    var encrypted_connections = crypto.Encrypt(settings.User.NodeSessionKey[pack.From.Name], connections)

                    encrypted_session_key, err := crypto.EncryptRSA(
                        settings.User.NodeSessionKey[pack.From.Name],
                        settings.User.NodePublicKey[pack.From.Name],
                    )
                    utils.CheckError(err)

                    var new_pack = settings.PackageTCP {
                        From: models.From {
                            Address: encrypted_address,
                            Login: encrypted_login,
                            Name: encrypted_name,
                        },
                        Head: models.Head {
                            Header: settings.HEAD_CONNECT,
                            Mode: settings.MODE_SAVE,
                        },
                        Body: hex.EncodeToString(encrypted_session_key) + 
                            settings.SEPARATOR + hex.EncodeToString([]byte(settings.User.PublicData)) +
                            settings.SEPARATOR + encrypted_connections,
                    }

                    var return_code = sendAddrPackage(pack.From.Address, new_pack)

                    if return_code == settings.EXIT_SUCCESS {
                        var encrypted_login = crypto.Encrypt(settings.User.Password, pack.From.Login)
                        settings.Mutex.Lock()
                        settings.User.NodeConnection[pack.From.Name] = 1
                        settings.Messages.NewDataExistLocal[pack.From.Name] = make(chan bool)
                        settings.User.Connections = append(
                            settings.User.Connections, 
                            pack.From.Name,
                        )
                        _, err = settings.DataBase.Exec(`
CREATE TABLE IF NOT EXISTS Local` + pack.From.Name + ` (
Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
User VARCHAR(64),
Body TEXT
);
INSERT INTO Connections(User, Login, PublicKey)
SELECT '` + pack.From.Name + `', '` + encrypted_login + `', '` + public_data + `'
WHERE NOT EXISTS(SELECT 1 FROM Connections WHERE User = '` + pack.From.Name + `');
`)
                        settings.Mutex.Unlock()
                        utils.CheckError(err)
                    } else {
                        nullNode(pack.From.Name)
                    }

                case settings.MODE_SAVE: 
                    var splited = strings.Split(pack.Body, settings.SEPARATOR)

                    encrypted_session_key, err := hex.DecodeString(splited[0])
                    utils.CheckError(err)

                    session_key, err := crypto.DecryptRSA(
                        encrypted_session_key,
                        settings.User.PrivateKey,
                    )
                    utils.CheckError(err)

                    public_key, err := hex.DecodeString(splited[1])
                    utils.CheckError(err)

                    var public_data = string(public_key)

                    var address = crypto.Decrypt(session_key, pack.From.Address) 
                    var login = crypto.Decrypt(session_key, pack.From.Login)
                    var username = crypto.Decrypt(session_key, pack.From.Name)

                    var encrypted_login = crypto.Encrypt(settings.User.Password, login)

                    settings.Mutex.Lock()
                    settings.User.NodeAddress[username] = address
                    settings.User.NodeLogin[username] = login
                    settings.User.NodePublicKey[username] = encoding.DecodePublic(string(public_key))
                    settings.User.NodeSessionKey[username] = session_key
                    settings.User.NodeConnection[username] = 1
                    settings.Messages.NewDataExistLocal[username] = make(chan bool)
                    settings.User.Connections = append(
                        settings.User.Connections, 
                        username,
                    )
                    _, err = settings.DataBase.Exec(`
CREATE TABLE IF NOT EXISTS Local` + username + ` (
Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
User VARCHAR(64),
Body TEXT
);
INSERT INTO Connections(User, Login, PublicKey) 
SELECT '` + username + `', '` + encrypted_login + `', '` + public_data + `'
WHERE NOT EXISTS(SELECT 1 FROM Connections WHERE User = '` + username + `');
`)
                    
                    settings.Mutex.Unlock()
                    utils.CheckError(err)

                case settings.MODE_READ_LIST:
                    var addresses = make([]string, len(settings.User.Connections))
                    for index, username := range settings.User.Connections {
                        addresses[index] = settings.User.NodeAddress[username]
                    }

                    var connections = strings.Join(addresses, settings.SEPARATOR_ADDRESS)
                    var new_pack = settings.PackageTCP {
                        From: models.From {
                            Address: settings.User.IPv4 + settings.User.Port,
                            Name: settings.User.Hash,
                        },
                        To: pack.From.Name,
                        Head: models.Head {
                            Header: settings.HEAD_CONNECT,
                            Mode: settings.MODE_SAVE_LIST,
                        },
                        Body: connections,
                    }
                    SendEncryptedPackage(new_pack)

                case settings.MODE_SAVE_LIST:
                    var connections = strings.Split(pack.Body, settings.SEPARATOR_ADDRESS)
                    Connect(connections)
            }

        case settings.HEAD_WARNING:
            switch pack.Mode {
                case settings.MODE_SAVE:
                    nullNode(pack.From.Name)
            }

        case settings.HEAD_EMAIL:
            switch pack.Mode {
                case settings.MODE_SAVE: 
                    var splited = strings.Split(pack.Body, settings.SEPARATOR)

                    if len(splited) != 3 {
                        goto close_connection
                    }

                    settings.Mutex.Lock()
                    _, err := settings.DataBase.Exec(
                        "INSERT INTO Email (Title, Body, User, Date) VALUES ($1, $2, $3, $4)",
                        crypto.Encrypt(settings.User.Password, splited[0]),
                        crypto.Encrypt(settings.User.Password, splited[1]),
                        pack.From.Name,
                        crypto.Encrypt(settings.User.Password, splited[2]),
                    )
                    settings.Mutex.Unlock()
                    utils.CheckError(err)
            }

        default:
            // pass
    }

close_connection:
    conn.Close()

    return 0
}
