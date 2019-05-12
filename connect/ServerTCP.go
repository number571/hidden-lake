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

type ModeConn int8
const (
    _DEFAULT ModeConn = 0
    _MERGE   ModeConn = 1
)

var (
    __packet_is_exist = make(map[string]bool) 
    __check_connection = make(map[string]bool) 
)

func ServerTCP() {
    var err error

    settings.ServerListenTCP, err = net.Listen(settings.PROTOCOL, settings.IPV4_TEMPLATE + settings.User.Port) 
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

    var pack models.PackageTCP
    err := json.Unmarshal([]byte(message), &pack)
    if err != nil {
        fmt.Println(err)
        return
    }

    var mode = decryptPackage(&pack)
    packageActions(pack, mode)
}

// Decrypt package if connection exists.
func decryptPackage(pack *models.PackageTCP) settings.ModeNet {
    var (
        session_key []byte = nil
        mode settings.ModeNet
    )

    // F2F
    if _, ok := settings.Node.Address.F2F[pack.From.Name]; ok {
        var return_code = crypto.TryDecrypt(settings.Node.SessionKey.F2F[pack.From.Name], pack.Head.Title)
        if return_code == 0 { 
            session_key = settings.Node.SessionKey.F2F[pack.From.Name]
            mode = settings.F2F_mode
        }
    }

    // P2P
    if mode != settings.F2F_mode && settings.Node.Connection[pack.From.Name] == 1 &&
        pack.Head.Title != settings.HEAD_CONNECT {
        var return_code = crypto.TryDecrypt(settings.Node.SessionKey.P2P[pack.From.Name], pack.Head.Title)
        if return_code == 0 {
            session_key = settings.Node.SessionKey.P2P[pack.From.Name]
        }
    }

    if session_key == nil { return settings.P2P_mode }

    *pack = models.PackageTCP {
        From: models.From {
            Name: pack.From.Name,
            Login: crypto.Decrypt(session_key, pack.From.Login),
            Address: crypto.Decrypt(session_key, pack.From.Address),
        },
        To: crypto.Decrypt(session_key, pack.To),
        Head: models.Head {
            Title: crypto.Decrypt(session_key, pack.Head.Title),
            Mode: crypto.Decrypt(session_key, pack.Head.Mode),
        }, 
        Body: crypto.Decrypt(session_key, pack.Body),
    }

    // data, _ := json.MarshalIndent(pack, "", "\t")
    // fmt.Println(string(data))

    return mode
}

// Actions with package.
func packageActions(pack models.PackageTCP, mode settings.ModeNet) {
    switch pack.Head.Title {
        case settings.HEAD_REDIRECT:
            switch mode {
                case settings.P2P_mode: redirectP2P(pack)
                case settings.F2F_mode: redirectF2F(pack) 
            }
            
        case settings.HEAD_ARCHIVE: 
            switch pack.Head.Mode {
                case settings.MODE_READ_LIST: archiveReadList(pack, mode)
                case settings.MODE_SAVE_LIST: archiveSaveList(pack)
                case settings.MODE_READ_FILE: archiveReadFile(pack, mode)
                case settings.MODE_SAVE_FILE: archiveSaveFile(pack)
            }

        case settings.HEAD_MESSAGE: 
            var message string
            switch mode {
                case settings.P2P_mode: message = fmt.Sprintf("[%s]: %s\n", settings.Node.Login[pack.From.Name], pack.Body)
                case settings.F2F_mode: message = fmt.Sprintf("[%s]: %s\n", pack.From.Name, pack.Body)
            }
            fmt.Print(message)
            switch pack.Head.Mode {
                case settings.MODE_LOCAL: messageLocal(pack, message, mode)
                case settings.MODE_GLOBAL: messageGlobal(pack, message, mode)
            }

        case settings.HEAD_CONNECT:
            if mode == settings.F2F_mode { return }
            switch pack.Head.Mode {
                // Connect with sending public keys.
                case settings.MODE_READ: connectRead(pack, _DEFAULT)
                case settings.MODE_SAVE: connectSave(pack, _DEFAULT)

                // Get connections after merge.
                case settings.MODE_GLOBAL: connectGlobal(pack)

                // Merge connection.
                case settings.MODE_READ_GLOBAL: connectRead(pack, _MERGE)
                case settings.MODE_SAVE_GLOBAL: connectSave(pack, _MERGE)

                // Attempt to connect without sending public keys.
                case settings.MODE_READ_LOCAL: connectReadLocal(pack)
                case settings.MODE_SAVE_LOCAL: connectSaveLocal(pack)

                // Check connection.
                case settings.MODE_READ_CHECK: connectReadCheck(pack)
                case settings.MODE_SAVE_CHECK: connectSaveCheck(pack)
            }

        case settings.HEAD_EMAIL:
            switch pack.Head.Mode {
                case settings.MODE_SAVE: emailSave(pack)
            }
    }
}

func connectGlobal(pack models.PackageTCP) {
    var connects = strings.Split(pack.Body, settings.SEPARATOR)
    connectP2P(connects, true)
}

// Get request and create connection.
func connectRead(pack models.PackageTCP, mode ModeConn) {
    if settings.Node.Connection[pack.From.Name] == 1 {
        return
    }

    var connects []string
    if mode == _MERGE {
        connects = settings.MakeAddresses(settings.Node.Address.P2P)
    }

    public_key, err := hex.DecodeString(pack.Body)
    utils.CheckError(err)

    var public_data = string(public_key)

    settings.Mutex.Lock()
    settings.Node.Address.P2P[pack.From.Name] = pack.From.Address
    settings.Node.Login[pack.From.Name] = pack.From.Login
    settings.Node.PublicKey[pack.From.Name] = encoding.DecodePublic(string(public_key))
    settings.Node.SessionKey.P2P[pack.From.Name] = crypto.SessionKey(settings.SESSION_KEY_BYTES)
    settings.Mutex.Unlock()

    var encrypted_address = crypto.Encrypt(settings.Node.SessionKey.P2P[pack.From.Name], settings.User.IPv4 + settings.User.Port)
    var encrypted_login = crypto.Encrypt(settings.Node.SessionKey.P2P[pack.From.Name], settings.User.Login)
    var encrypted_name = crypto.Encrypt(settings.Node.SessionKey.P2P[pack.From.Name], settings.User.Hash.P2P)

    encrypted_session_key, err := crypto.EncryptRSA(
        settings.Node.SessionKey.P2P[pack.From.Name],
        settings.Node.PublicKey[pack.From.Name],
    )
    utils.CheckError(err)

    var new_pack = models.PackageTCP {
        From: models.From {
            Address: encrypted_address,
            Login: encrypted_login,
            Name: encrypted_name,
        },
        Head: models.Head {
            Title: settings.HEAD_CONNECT,
            Mode: settings.MODE_SAVE,
        },
        Body: hex.EncodeToString(encrypted_session_key) + 
            settings.SEPARATOR + hex.EncodeToString([]byte(settings.User.Public.Data.P2P)) + 
            settings.SEPARATOR + strings.Join(connects, settings.SEPARATOR),
    }

    if mode == _MERGE {
        new_pack.Head.Mode = settings.MODE_SAVE_GLOBAL
    }

    var return_code = sendAddrPackage(pack.From.Address, new_pack)

    if return_code == settings.EXIT_SUCCESS {
        var encrypted_login = crypto.Encrypt(settings.User.Password, pack.From.Login)
        settings.Mutex.Lock()
        settings.Node.Connection[pack.From.Name] = 1
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

// Accept connection.
func connectSave(pack models.PackageTCP, mode ModeConn) {
    var (
        splited = strings.Split(pack.Body, settings.SEPARATOR)
        connected_nodes = settings.Node.Address.P2P
        connects = splited[2:]
    )

    if mode == _MERGE {
        go connectP2P(connects, true)
    }

    encrypted_session_key, err := hex.DecodeString(splited[0])
    utils.CheckError(err)

    session_key, err := crypto.DecryptRSA(
        encrypted_session_key,
        settings.User.Private.Key.P2P,
    )
    utils.CheckError(err)

    public_key, err := hex.DecodeString(splited[1])
    utils.CheckError(err)

    var public_data = string(public_key)

    var address = crypto.Decrypt(session_key, pack.From.Address) 
    var login = crypto.Decrypt(session_key, pack.From.Login)
    var username = crypto.Decrypt(session_key, pack.From.Name)

    var encrypted_login = crypto.Encrypt(settings.User.Password, login)

    if mode == _MERGE {
        go redirectConnect(connected_nodes, append(connects, address))
    }

    settings.Mutex.Lock()
    settings.Node.Address.P2P[username] = address
    settings.Node.Login[username] = login
    settings.Node.PublicKey[username] = encoding.DecodePublic(string(public_key))
    settings.Node.SessionKey.P2P[username] = session_key
    settings.Node.Connection[username] = 1
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

// Connection with sending public key.
func connectSaveLocal(pack models.PackageTCP) {
    if settings.Node.Connection[pack.From.Name] == 1 {
        return
    }
    connectP2P([]string{pack.From.Address}, false)
}

// Connection with found public key in database.
func connectReadLocal(pack models.PackageTCP) {
    if settings.Node.Connection[pack.From.Name] == 1 {
        return
    }

    var row = settings.DataBase.QueryRow(
        "SELECT Login, PublicKey FROM Connections WHERE User = $1",
        pack.From.Name,
    )

    var login, public_data string
    row.Scan(&login, &public_data)

    if public_data == "" {
        var new_pack = models.PackageTCP {
            From: models.From {
                Name: settings.User.Hash.P2P,
                Address: settings.User.IPv4 + settings.User.Port,
            },
            Head: models.Head {
                Title: settings.HEAD_CONNECT,
                Mode: settings.MODE_SAVE_LOCAL,
            },
        }
        sendAddrPackage(pack.From.Address, new_pack)
        return
    }

    settings.Mutex.Lock()
    settings.Node.Address.P2P[pack.From.Name] = pack.From.Address
    settings.Node.Login[pack.From.Name] = crypto.Decrypt(settings.User.Password, login)
    settings.Node.PublicKey[pack.From.Name] = encoding.DecodePublic(string(public_data))
    settings.Node.SessionKey.P2P[pack.From.Name] = crypto.SessionKey(settings.SESSION_KEY_BYTES)
    settings.Mutex.Unlock()

    var encrypted_address = crypto.Encrypt(settings.Node.SessionKey.P2P[pack.From.Name], settings.User.IPv4 + settings.User.Port)
    var encrypted_login = crypto.Encrypt(settings.Node.SessionKey.P2P[pack.From.Name], settings.User.Login)
    var encrypted_name = crypto.Encrypt(settings.Node.SessionKey.P2P[pack.From.Name], settings.User.Hash.P2P)

    encrypted_session_key, err := crypto.EncryptRSA(
        settings.Node.SessionKey.P2P[pack.From.Name],
        settings.Node.PublicKey[pack.From.Name],
    )
    utils.CheckError(err)

    var new_pack = models.PackageTCP {
        From: models.From {
            Address: encrypted_address,
            Login: encrypted_login,
            Name: encrypted_name,
        },
        Head: models.Head {
            Title: settings.HEAD_CONNECT,
            Mode: settings.MODE_SAVE,
        },
        Body: hex.EncodeToString(encrypted_session_key) + 
            settings.SEPARATOR + hex.EncodeToString([]byte(settings.User.Public.Data.P2P)),
    }

    var return_code = sendAddrPackage(pack.From.Address, new_pack)

    if return_code == settings.EXIT_SUCCESS {
        settings.Mutex.Lock()
        settings.Node.Connection[pack.From.Name] = 1
        settings.Messages.NewDataExistLocal[pack.From.Name] = make(chan bool)
        settings.Mutex.Unlock()
    } else {
        nullNode(pack.From.Name)
    }
}

func connectReadCheck(pack models.PackageTCP) {
    var new_pack = models.PackageTCP {
        From: models.From {
            Name: settings.User.Hash.P2P,
        },
        To: pack.From.Name,
        Head: models.Head {
            Title: settings.HEAD_CONNECT,
            Mode: settings.MODE_SAVE_CHECK,
        },
    }
    sendEncryptedPackage(new_pack, settings.P2P_mode)
}

func connectSaveCheck(pack models.PackageTCP) {
    __check_connection[pack.From.Name] = true
}

// Read package and/or send to friends.
func redirectF2F(pack models.PackageTCP) {
    var (
        to_heads = strings.Split(pack.Head.Mode, settings.SEPARATOR)
        to = to_heads[0]
        heads = to_heads[1:]
    )

    if _, ok := __packet_is_exist[heads[2]]; ok { return }

    settings.Mutex.Lock()
    __packet_is_exist[heads[2]] = true
    settings.Mutex.Unlock()

    go func() {
        time.Sleep(time.Second * 15)
        settings.Mutex.Lock()
        delete(__packet_is_exist, heads[2])
        settings.Mutex.Unlock()
    }()

    if to != settings.User.Hash.F2F {
        sendRedirectF2FPackage(pack)
        if to != "" { return }
    }

    var new_pack = models.PackageTCP {
        From: models.From {
            Name: pack.From.Login,
        },
        To: settings.User.Hash.F2F,
        Head: models.Head {
            Title: heads[0],
            Mode: heads[1],
        },
        Body: pack.Body,
    }

    packageActions(new_pack, settings.F2F_mode)
}

// Check package for receiving or redirecting.
func redirectP2P(pack models.PackageTCP) {
    var hashname_heads = strings.Split(pack.Head.Mode, settings.SEPARATOR)
    decoded_hashname, err := hex.DecodeString(hashname_heads[0])
    utils.CheckError(err)

    bytes_hashname, err := crypto.DecryptRSA(decoded_hashname, settings.User.Private.Key.P2P)
    if err != nil {
        sendRedirectP2PPackage(pack)
        return
    }

    var (
        hashname = string(bytes_hashname)
        heads = strings.Split(
            crypto.Decrypt(settings.Node.SessionKey.P2P[hashname], hashname_heads[1]),
            settings.SEPARATOR,
        )
    )

    var new_pack = models.PackageTCP {
        From: models.From {
            Name: hashname,
        },
        To: settings.User.Hash.P2P,
        Head: models.Head {
            Title: heads[0],
            Mode: heads[1],
        },
        Body: crypto.Decrypt(settings.Node.SessionKey.P2P[hashname], pack.Body),
    }

    packageActions(new_pack, settings.P2P_mode)
}

// Save received email in database.
func emailSave(pack models.PackageTCP) {
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

// Send global message to all nodes.
func messageGlobal(pack models.PackageTCP, message string, mode settings.ModeNet) {
    var mode_name = "P2P"
    if mode == settings.F2F_mode { 
        mode_name = "F2F"
        if _, ok := settings.Node.Address.F2F[pack.From.Name]; !ok {
            message = "(HiddenFriend)" + message
        }
    }
    settings.Mutex.Lock()
    _, err := settings.DataBase.Exec(
        "INSERT INTO GlobalMessages (User, Mode, Body) VALUES ($1, $2, $3)",
        pack.From.Name, 
        mode_name,
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
func messageLocal(pack models.PackageTCP, message string, mode settings.ModeNet) {
    var mode_name = "P2P"
    if mode == settings.F2F_mode { 
        mode_name = "F2F"
        if _, ok := settings.Node.Address.F2F[pack.From.Name]; !ok {
            messageGlobal(pack, message, settings.F2F_mode)
            return
        }
    }

    settings.Mutex.Lock()
    _, err := settings.DataBase.Exec(
        "INSERT INTO Local" + pack.From.Name + " (User, Mode, Body) VALUES ($1, $2, $3)",
        pack.From.Name,
        mode_name,
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
func archiveSaveFile(pack models.PackageTCP) {
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
func archiveReadFile(pack models.PackageTCP, mode settings.ModeNet) {
    if utils.FileIsExist(settings.PATH_ARCHIVE + pack.Body) && 
      !strings.Contains(pack.Body, "..") {
        var new_pack = models.PackageTCP {
            From: models.From {
                Name: pack.To,
            },
            To: pack.From.Name,
            Head: models.Head {
                Title: settings.HEAD_ARCHIVE,
                Mode: settings.MODE_SAVE_FILE,
            },
            Body: hex.EncodeToString([]byte(pack.Body)) + settings.SEPARATOR + 
                hex.EncodeToString([]byte(utils.ReadFile(settings.PATH_ARCHIVE + pack.Body))),
        }
        SendPackage(new_pack, mode)
    }
}

// Send list of files from archive.
func archiveReadList(pack models.PackageTCP, mode settings.ModeNet) {
    files, err := ioutil.ReadDir(settings.PATH_ARCHIVE)
    utils.CheckError(err)
    var list = make([]string, len(files))
    for index, file := range files {
        list[index] = file.Name()
    }

    var list_of_files = strings.Join(list, settings.SEPARATOR)
    var new_pack = models.PackageTCP {
        From: models.From {
            Name: pack.To,
        },
        To: pack.From.Name,
        Head: models.Head {
            Title: settings.HEAD_ARCHIVE,
            Mode: settings.MODE_SAVE_LIST,
        },
        Body: list_of_files,
    }
    SendPackage(new_pack, mode)
}

// Save list of files.
func archiveSaveList(pack models.PackageTCP) {
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
