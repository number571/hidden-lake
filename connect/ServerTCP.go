package connect

import (
    "fmt"
    "net"
    "time"
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

var temp_packets = make(map[string]bool) 

func ServerTCP() {
    var err error

    settings.ServerListenTCP, err = net.Listen(settings.PROTOCOL_TCP, settings.IPV4_TEMPLATE + settings.User.Port) 
    utils.CheckError(err)

    for {
        conn, err := settings.ServerListenTCP.Accept()
        if err != nil { break }
        go server(conn)
    }
}

func server(conn net.Conn) {
    var buffer = make([]byte, settings.BUFF_SIZE)
    var message string

    for {
        length, err := conn.Read(buffer)
        if length == 0 || err != nil { break }
        message += string(buffer[:length])
    }

    conn.Close()

    var pack settings.PackageTCP
    err := json.Unmarshal([]byte(message), &pack)
    if err != nil {
        fmt.Println(err)
        return
    }

    var is_f2f = decryptPackage(&pack)
    packageActions(pack, is_f2f)
}

// Decrypt package if connection exists.
func decryptPackage(pack *settings.PackageTCP) bool {
    var (
        session_key []byte = nil
        is_f2f bool
    )

    // F2F
    if _, ok := settings.User.NodeAddressF2F[pack.From.Name]; ok {
        var return_code = crypto.TryDecrypt(settings.User.NodeSessionKeyF2F[pack.From.Name], pack.Header)
        if return_code == 0 { 
            session_key = settings.User.NodeSessionKeyF2F[pack.From.Name]
            is_f2f = true
        }
    }

    // P2P
    if !is_f2f && pack.Head.Header != settings.HEAD_CONNECT &&
       settings.User.NodeConnection[pack.From.Name] == 1 {
        var return_code = crypto.TryDecrypt(settings.User.NodeSessionKey[pack.From.Name], pack.Header)
        if return_code == 0 { 
            session_key = settings.User.NodeSessionKey[pack.From.Name]
        }
    }

    if session_key == nil { return false }

    *pack = settings.PackageTCP {
        From: models.From {
            Name: pack.From.Name,
            Login: crypto.Decrypt(session_key, pack.From.Login),
            Address: crypto.Decrypt(session_key, pack.From.Address),
        },
        To: crypto.Decrypt(session_key, pack.To),
        Head: models.Head {
            Header: crypto.Decrypt(session_key, pack.Header),
            Mode: crypto.Decrypt(session_key, pack.Mode),
        }, 
        Body: crypto.Decrypt(session_key, pack.Body),
    }

    return is_f2f
}

// Actions with package.
func packageActions(pack settings.PackageTCP, is_f2f bool) {
    switch pack.Header {
        case settings.HEAD_REDIRECT:
            if is_f2f { 
                redirectF2F(pack) 
            } else {
                redirect(pack)
            }
            
        case settings.HEAD_ARCHIVE: 
            switch pack.Mode {
                case settings.MODE_READ_LIST: archiveReadList(pack, is_f2f)
                case settings.MODE_SAVE_LIST: archiveSaveList(pack)
                case settings.MODE_READ_FILE: archiveReadFile(pack, is_f2f)
                case settings.MODE_SAVE_FILE: archiveSaveFile(pack)
            }

        case settings.HEAD_MESSAGE: 
            var message string
            if is_f2f {
                message = fmt.Sprintf("[%s]: %s\n", pack.From.Name, pack.Body)
            } else {
                message = fmt.Sprintf("[%s]: %s\n", settings.User.NodeLogin[pack.From.Name], pack.Body)
            }
            fmt.Print(message)
            switch pack.Mode {
                case settings.MODE_LOCAL: messageLocal(pack, message, is_f2f)
                case settings.MODE_GLOBAL: messageGlobal(pack, message, is_f2f)
            }

        case settings.HEAD_CONNECT:
            if is_f2f { return }
            switch pack.Mode {
                case settings.MODE_GLOBAL: connectGlobal(pack)
                case settings.MODE_LOCAL: connectLocal(pack)
                case settings.MODE_READ: connectRead(pack)
                case settings.MODE_SAVE: connectSave(pack)
                case settings.MODE_READ_LIST: connectReadList(pack)
                case settings.MODE_SAVE_LIST: connectSaveList(pack)
            }

        case settings.HEAD_EMAIL:
            switch pack.Mode {
                case settings.MODE_SAVE: emailSave(pack)
            }
    }
}

// Read package and send to friends.
func redirectF2F(pack settings.PackageTCP) {
    var (
        to_heads = strings.Split(pack.Head.Mode, settings.SEPARATOR)
        to = to_heads[0]
        heads = to_heads[1:]
    )

    if _, ok := temp_packets[heads[2]]; ok { return }

    settings.Mutex.Lock()
    temp_packets[heads[2]] = true
    settings.Mutex.Unlock()

    go func() {
        time.Sleep(time.Second * 15)
        settings.Mutex.Lock()
        delete(temp_packets, heads[2])
        settings.Mutex.Unlock()
    }()

    if to != settings.User.Hash {
        sendRedirectF2FPackage(pack)
        if to != "" { return }
    }

    var new_pack = settings.PackageTCP {
        From: models.From {
            Name: pack.From.Login,
        },
        To: settings.User.Hash,
        Head: models.Head {
            Header: heads[0],
            Mode: heads[1],
        },
        Body: pack.Body,
    }

    packageActions(new_pack, true)
}

// Check package for receiving or redirecting.
func redirect(pack settings.PackageTCP) {
    var hashname_heads = strings.Split(pack.Head.Mode, settings.SEPARATOR)
    decoded_hashname, err := hex.DecodeString(hashname_heads[0])
    utils.CheckError(err)

    bytes_hashname, err := crypto.DecryptRSA(decoded_hashname, settings.User.PrivateKey)
    if err != nil {
        sendRedirectPackage(pack)
        return
    }

    var (
        hashname = string(bytes_hashname)
        heads = strings.Split(
            crypto.Decrypt(settings.User.NodeSessionKey[hashname], hashname_heads[1]),
            settings.SEPARATOR,
        )
    )

    var new_pack = settings.PackageTCP {
        From: models.From {
            Name: hashname,
        },
        To: settings.User.Hash,
        Head: models.Head {
            Header: heads[0],
            Mode: heads[1],
        },
        Body: crypto.Decrypt(settings.User.NodeSessionKey[hashname], pack.Body),
    }

    packageActions(new_pack, false)
}

// Save received email in database.
func emailSave(pack settings.PackageTCP) {
    var splited = strings.Split(pack.Body, settings.SEPARATOR)
    if len(splited) != 3 { return }
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

// Save connections.
func connectSaveList(pack settings.PackageTCP) {
    var connections = strings.Split(pack.Body, settings.SEPARATOR)
    Connect(connections, false)
}

// Send connections.
func connectReadList(pack settings.PackageTCP) {
    var (
        connects = settings.MakeConnects(settings.User.NodeAddress)
        connections = strings.Join(connects, settings.SEPARATOR)
    )
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
    sendEncryptedPackage(new_pack, false)
}

// Accept connection.
func connectSave(pack settings.PackageTCP) {
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
    _, err = settings.DataBase.Exec(`
CREATE TABLE IF NOT EXISTS Local` + username + ` (
Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
User VARCHAR(64),
Mode VARCHAR(3),
Body TEXT
);
INSERT INTO Connections(User, Login, PublicKey) 
SELECT '` + username + `', '` + encrypted_login + `', '` + public_data + `'
WHERE NOT EXISTS(SELECT 1 FROM Connections WHERE User = '` + username + `');
`)
    
    settings.Mutex.Unlock()
    utils.CheckError(err)
}

// Get request and create connection.
func connectRead(pack settings.PackageTCP) {
    if settings.User.NodeConnection[pack.From.Name] == 1 {
        return
    }
    
    public_key, err := hex.DecodeString(pack.Body)
    utils.CheckError(err)

    var public_data = string(public_key)

    settings.Mutex.Lock()
    settings.User.NodeAddress[pack.From.Name] = pack.From.Address
    settings.User.NodeLogin[pack.From.Name] = pack.From.Login
    settings.User.NodePublicKey[pack.From.Name] = encoding.DecodePublic(string(public_key))
    settings.User.NodeSessionKey[pack.From.Name] = crypto.SessionKey(settings.SESSION_KEY_BYTES)
    settings.Mutex.Unlock()

    var encrypted_address = crypto.Encrypt(settings.User.NodeSessionKey[pack.From.Name], settings.User.IPv4 + settings.User.Port)
    var encrypted_login = crypto.Encrypt(settings.User.NodeSessionKey[pack.From.Name], settings.User.Login)
    var encrypted_name = crypto.Encrypt(settings.User.NodeSessionKey[pack.From.Name], settings.User.Hash)

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
            settings.SEPARATOR + hex.EncodeToString([]byte(settings.User.PublicData)),
    }

    var return_code = sendAddrPackage(pack.From.Address, new_pack)

    if return_code == settings.EXIT_SUCCESS {
        var encrypted_login = crypto.Encrypt(settings.User.Password, pack.From.Login)
        settings.Mutex.Lock()
        settings.User.NodeConnection[pack.From.Name] = 1
        settings.Messages.NewDataExistLocal[pack.From.Name] = make(chan bool)
        _, err = settings.DataBase.Exec(`
CREATE TABLE IF NOT EXISTS Local` + pack.From.Name + ` (
Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
User VARCHAR(64),
Mode VARCHAR(3),
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
}

// Connection with sending public key.
func connectLocal(pack settings.PackageTCP) {
    if settings.User.NodeConnection[pack.From.Name] == 1 {
        return
    }
    Connect([]string{pack.From.Address}, true)
}

// Connection with found public key in database.
func connectGlobal(pack settings.PackageTCP) {
    if settings.User.NodeConnection[pack.From.Name] == 1 {
        return
    }

    var row = settings.DataBase.QueryRow(
        "SELECT Login, PublicKey FROM Connections WHERE User = $1",
        pack.From.Name,
    )

    var login, public_data string
    row.Scan(
        &login,
        &public_data,
    )

    if public_data == "" {
        var new_pack = settings.PackageTCP {
            From: models.From {
                Name: settings.User.Hash,
                Address: settings.User.IPv4 + settings.User.Port,
            },
            Head: models.Head {
                Header: settings.HEAD_CONNECT,
                Mode: settings.MODE_LOCAL,
            },
        }
        sendAddrPackage(pack.From.Address, new_pack)
        return
    }

    settings.Mutex.Lock()
    settings.User.NodeAddress[pack.From.Name] = pack.From.Address
    settings.User.NodeLogin[pack.From.Name] = crypto.Decrypt(settings.User.Password, login)
    settings.User.NodePublicKey[pack.From.Name] = encoding.DecodePublic(string(public_data))
    settings.User.NodeSessionKey[pack.From.Name] = crypto.SessionKey(settings.SESSION_KEY_BYTES)
    settings.Mutex.Unlock()

    var encrypted_address = crypto.Encrypt(settings.User.NodeSessionKey[pack.From.Name], settings.User.IPv4 + settings.User.Port)
    var encrypted_login = crypto.Encrypt(settings.User.NodeSessionKey[pack.From.Name], settings.User.Login)
    var encrypted_name = crypto.Encrypt(settings.User.NodeSessionKey[pack.From.Name], settings.User.Hash)

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
            settings.SEPARATOR + hex.EncodeToString([]byte(settings.User.PublicData)),
    }

    var return_code = sendAddrPackage(pack.From.Address, new_pack)

    if return_code == settings.EXIT_SUCCESS {
        settings.Mutex.Lock()
        settings.User.NodeConnection[pack.From.Name] = 1
        settings.Messages.NewDataExistLocal[pack.From.Name] = make(chan bool)
        settings.Mutex.Unlock()
    } else {
        nullNode(pack.From.Name)
    }
}

// Send global message to all nodes.
func messageGlobal(pack settings.PackageTCP, message string, is_f2f bool) {
    var mode = "P2P"
    if is_f2f { 
        mode = "F2F"
        if _, ok := settings.User.NodeAddressF2F[pack.From.Name]; !ok {
            message = "(HiddenFriend)" + message
        }
    }
    settings.Mutex.Lock()
    _, err := settings.DataBase.Exec(
        "INSERT INTO GlobalMessages (User, Mode, Body) VALUES ($1, $2, $3)",
        pack.From.Name, 
        mode,
        crypto.Encrypt(settings.User.Password, message),
    )
    settings.Messages.CurrentIdGlobal++
    settings.Mutex.Unlock()
    utils.CheckError(err)
    go func() {
        settings.Messages.NewDataExistGlobal <- true
    }()
}

// Send local message to one node.
func messageLocal(pack settings.PackageTCP, message string, is_f2f bool) {
    var mode = "P2P"
    if is_f2f { 
        mode = "F2F"
        if _, ok := settings.User.NodeAddressF2F[pack.From.Name]; !ok {
            messageGlobal(pack, message, true)
            return
        }
    }

    settings.Mutex.Lock()
    _, err := settings.DataBase.Exec(
        "INSERT INTO Local" + pack.From.Name + " (User, Mode, Body) VALUES ($1, $2, $3)",
        pack.From.Name,
        mode,
        crypto.Encrypt(settings.User.Password, message),
    )
    settings.Messages.CurrentIdLocal[pack.From.Name]++
    settings.Mutex.Unlock()
    utils.CheckError(err)

    go func() {
        settings.Messages.NewDataExistLocal[pack.From.Name] <- true
    }()
}

// Get and save file in archive.
func archiveSaveFile(pack settings.PackageTCP) {
    var (
        splited = strings.Split(pack.Body, settings.SEPARATOR)
        add string
    )
    if utils.FileIsExist(settings.PATH_ARCHIVE + splited[0]) {
        add += "copy_"
    }

    filename, err := hex.DecodeString(splited[0])
    utils.CheckError(err)

    body, err := hex.DecodeString(splited[1])
    utils.CheckError(err)

    utils.WriteFile(settings.PATH_ARCHIVE + add + string(filename), string(body))
}

// Send file from archive.
func archiveReadFile(pack settings.PackageTCP, is_f2f bool) {
    if utils.FileIsExist(settings.PATH_ARCHIVE + pack.Body) && 
      !strings.Contains(pack.Body, "..") {
        var new_pack = settings.PackageTCP {
            From: models.From {
                Name: pack.To,
            },
            To: pack.From.Name,
            Head: models.Head {
                Header: settings.HEAD_ARCHIVE,
                Mode: settings.MODE_SAVE_FILE,
            },
            Body: hex.EncodeToString([]byte(pack.Body)) + settings.SEPARATOR + 
                hex.EncodeToString([]byte(utils.ReadFile(settings.PATH_ARCHIVE + pack.Body))),
        }
        SendPackage(new_pack, is_f2f)
    }
}

// Send list of files from archive.
func archiveReadList(pack settings.PackageTCP, is_f2f bool) {
    files, err := ioutil.ReadDir(settings.PATH_ARCHIVE)
    utils.CheckError(err)
    var list = make([]string, len(files))
    for index, file := range files {
        list[index] = file.Name()
    }

    var list_of_files = strings.Join(list, settings.SEPARATOR)
    var new_pack = settings.PackageTCP {
        From: models.From {
            Name: pack.To,
        },
        To: pack.From.Name,
        Head: models.Head {
            Header: settings.HEAD_ARCHIVE,
            Mode: settings.MODE_SAVE_LIST,
        },
        Body: list_of_files,
    }
    SendPackage(new_pack, is_f2f)
}

// Save list of files.
func archiveSaveList(pack settings.PackageTCP) {
    settings.Mutex.Lock()
    settings.User.TempArchive = strings.Split(pack.Body, settings.SEPARATOR)
    settings.Mutex.Unlock()

    fmt.Printf("| %s:\n", pack.From.Name)
    for _, file := range settings.User.TempArchive {
        if file != "" {
            fmt.Println("|", file)
        }
    }
}
