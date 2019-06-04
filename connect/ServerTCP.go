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

type ConnectionType int8
const (
    _DEFAULT ConnectionType = 0
    _MERGE   ConnectionType = 1
)

var (
    // check packet in exists by F2F nodes.
    __packet_is_exist = make(map[string]bool) 

    // check node connection in real time.
    __check_connection = make(map[string]bool) 
)

// Raising server on set port.
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

// Read received package and mode from client or p2p node.
func server(conn net.Conn) {
    var (
        buffer = make([]byte, settings.BUFF_SIZE)
        message string
    )

    for {
        length, err := conn.Read(buffer)
        if length == 0 || err != nil { break }
        message += string(buffer[:length])

        // CLIENT-SERVER
        if strings.HasSuffix(message, settings.END_BLOCK) {
            message = strings.TrimSuffix(message, settings.END_BLOCK)
            serve(conn, message)
            message = ""
        }
    }

    // P2P/F2F
    conn.Close()
    serve(nil, message)
}

// Unmarshal and decrypt package.
func serve(conn net.Conn, message string) {
    var pack models.PackageTCP
    err := json.Unmarshal([]byte(message), &pack)
    if err != nil {
        // fmt.Println(err)
        return
    }

    var is_arch_cs = false
    if conn != nil {
        settings.Node.Address.C_S[pack.From.Hash] = conn
        is_arch_cs = true
    }

    var mode = decryptPackage(&pack, is_arch_cs)
    packageActions(pack, mode)
}

// Decrypt package if connection exists.
func decryptPackage(pack *models.PackageTCP, is_arch_cs bool) models.ModeNet {
    var (
        session_key []byte  = nil
        mode models.ModeNet
    )

    if is_arch_cs { 
        mode = models.C_S_mode
    } else {
        mode = models.P2P_mode 
    }

    // F2F
    if _, ok := settings.Node.Address.F2F[pack.From.Hash]; ok {
        mode = models.F2F_mode
        if crypto.TryDecrypt(settings.Node.SessionKey.F2F[pack.From.Hash], pack.Head.Title) == 0 {
            session_key = settings.Node.SessionKey.F2F[pack.From.Hash]
        }

    // P2P
    } else if settings.Node.ConnectionMode[pack.From.Hash] == models.CONN && pack.Head.Title != settings.HEAD_CONNECT {
        if crypto.TryDecrypt(settings.Node.SessionKey.P2P[pack.From.Hash], pack.Head.Title) == 0 {
            session_key = settings.Node.SessionKey.P2P[pack.From.Hash]
        }
    } 

    // if settings.User.Mode == models.C_S_mode {
    //     fmt.Println(mode)
    //     data, _ := json.MarshalIndent(pack, "", "\t")
    //     fmt.Println(string(data))
    // }

    if session_key == nil { return mode }

    // Decrypt
    *pack = models.PackageTCP {
        From: models.From {
            Hash: pack.From.Hash,
            Address: crypto.Decrypt(session_key, pack.From.Address),
        },
        To: models.To {
            Hash: crypto.Decrypt(session_key, pack.To.Hash),
            Address: crypto.Decrypt(session_key, pack.To.Address),
        },
        Head: models.Head {
            Title: crypto.Decrypt(session_key, pack.Head.Title),
            Mode: crypto.Decrypt(session_key, pack.Head.Mode),
        }, 
        Body: crypto.Decrypt(session_key, pack.Body),
    }

    // if settings.User.Mode == models.C_S_mode {
    //     data, _ := json.MarshalIndent(pack, "", "\t")
    //     fmt.Println(string(data))
    // }

    return mode
}

// Actions with package.
func packageActions(pack models.PackageTCP, mode models.ModeNet) {
    switch pack.Head.Title {
        case settings.HEAD_REDIRECT:
            switch mode {
                case models.P2P_mode: redirectP2P(pack)
                case models.F2F_mode: redirectF2F(pack) 
                case models.C_S_mode: redirectC_S(pack)
            }

        case settings.HEAD_ARCHIVE: 
            switch pack.Head.Mode {
                case settings.MODE_READ_LIST: archiveReadList(pack, mode)
                case settings.MODE_SAVE_LIST: archiveSaveList(pack)
                case settings.MODE_READ_FILE: archiveReadFile(pack, mode)
                case settings.MODE_SAVE_FILE: archiveSaveFile(pack)
            }

        case settings.HEAD_MESSAGE:
            var message = readMessage(pack, mode)
            fmt.Print(message)
            switch pack.Head.Mode {
                case settings.MODE_LOCAL: messageLocal(pack, message, mode)
                case settings.MODE_GLOBAL: 
                    messageGlobal(pack, message, mode)
                    if mode == models.P2P_mode && settings.User.Mode != models.C_S_mode {
                        redirectToClients(pack)
                    }
            }

        case settings.HEAD_CONNECT:
            if mode == models.F2F_mode { return }
            switch pack.Head.Mode {
                // Connect with sending public keys.
                case settings.MODE_READ: connectRead(pack, mode, _DEFAULT)
                case settings.MODE_SAVE: connectSave(pack, mode, _DEFAULT)

                // Get connections after merge.
                case settings.MODE_GLOBAL: connectGlobal(pack)

                // Merge connection.
                case settings.MODE_READ_GLOBAL: connectRead(pack, mode, _MERGE)
                case settings.MODE_SAVE_GLOBAL: connectSave(pack, mode, _MERGE)

                // Attempt to connect without sending public keys.
                case settings.MODE_READ_LOCAL: connectReadLocal(pack, mode)
                case settings.MODE_SAVE_LOCAL: connectSaveLocal(pack)

                // Check connection in real time.
                case settings.MODE_READ_CHECK: connectReadCheck(pack)
                case settings.MODE_SAVE_CHECK: connectSaveCheck(pack)
            }

        case settings.HEAD_EMAIL:
            switch pack.Head.Mode {
                case settings.MODE_SAVE: emailSave(pack)
            }
    }
}

// Create message by mode, author and package body.
func readMessage(pack models.PackageTCP, mode models.ModeNet) string {
    var from = pack.From.Hash
    if pack.From.Address != "" && mode != models.F2F_mode {
        // from = pack.From.Hash + "->" + pack.From.Address
        if mode == models.C_S_mode {
            from = pack.From.Address
        } else if mode == models.P2P_mode {
            from = pack.From.Hash + "->" + pack.From.Address
        }
    }
    return fmt.Sprintf("(%s)[%s]: %s\n", settings.GetStrMode(mode), from, pack.Body)
}

// Send redirect-package from client to other nodes/clients.
func redirectC_S(pack models.PackageTCP) {
    var heads = strings.Split(pack.Head.Mode, settings.SEPARATOR)
    var new_pack = models.PackageTCP {
        From: models.From {
            Hash: settings.User.Hash.P2P,
            Address: pack.From.Address,
        },
        To: models.To {
            Hash: pack.To.Address,
        },
        Head: models.Head {
            Title: heads[0],
            Mode: heads[1],
        },
        Body: pack.Body,
    }

    // data, _ := json.MarshalIndent(pack, "", "\t")
    // fmt.Println(string(data))

    if new_pack.To.Hash == settings.User.Hash.P2P || 
      (new_pack.Head.Title == settings.HEAD_CONNECT && new_pack.To.Hash == "") {
        readActionPackage(new_pack)
        return
    }
    emulateClient(new_pack)
}

// Server requests by client.
func emulateClient(pack models.PackageTCP) {
    var (
        flag = false
        mode = models.P2P_mode
    )
    if _, ok := settings.Node.Address.C_S[pack.To.Hash]; ok {
        mode = models.C_S_mode
    }

    switch pack.Head.Title {
        case settings.HEAD_MESSAGE:
            switch pack.Head.Mode {
                case settings.MODE_GLOBAL: flag = true
                    // for username := range settings.Node.Address.P2P {
                    //     pack.To.Hash = username
                    //     SendEncryptedPackage(pack, models.P2P_mode)
                    // }
                    // for username := range settings.Node.Address.C_S {
                    //     if username == pack.From.Address { continue }
                    //     pack.To.Hash = username
                    //     SendEncryptedPackage(pack, models.C_S_mode)
                    // }
                    SendPackage(pack, models.P2P_mode)
                    readActionPackage(pack)
            }
    }

    if !flag {
        SendPackage(pack, mode)
    }
}

// Redirect message from p2p node to clients.
func redirectToClients(pack models.PackageTCP) {
    var new_pack = models.PackageTCP {
        From: models.From {
            Address: pack.From.Hash,
            Hash: settings.User.Hash.P2P,
        },
        Head: models.Head {
            Title: settings.HEAD_MESSAGE,
            Mode: settings.MODE_GLOBAL,
        },
        Body: pack.Body,
    }
    SendPackage(new_pack, models.C_S_mode)
    // for username := range settings.Node.Address.C_S {
    //     new_pack.To.Hash = username
    //     SendEncryptedPackage(new_pack, models.C_S_mode)
    // }
}

// Read package from client.
func readActionPackage(pack models.PackageTCP) {
    pack.From.Hash = pack.From.Address
    pack.To.Hash = settings.User.Hash.P2P
    packageActions(pack, models.C_S_mode)
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
            Hash: hashname,
            Address: pack.From.Address,
        },
        To: models.To {
            Hash: settings.User.Hash.P2P,
            Address: pack.To.Address,
        },
        Head: models.Head {
            Title: heads[0],
            Mode: heads[1],
        },
        Body: crypto.Decrypt(settings.Node.SessionKey.P2P[hashname], pack.Body),
    }

    if _, ok := settings.Node.Address.C_S[pack.From.Address]; ok {
        new_pack.To.Hash = pack.From.Address
        new_pack.From.Hash = settings.User.Hash.P2P
        new_pack.From.Address = pack.From.Hash
        SendPackage(new_pack, models.C_S_mode)
        return
    }

    packageActions(new_pack, models.P2P_mode)
}

// After merge connect to all nodes in list.
func connectGlobal(pack models.PackageTCP) {
    var connects = strings.Split(pack.Body, settings.SEPARATOR)
    connectP2P(connects, true)
}

// Get request and create connection.
func connectRead(pack models.PackageTCP, mode models.ModeNet, conn_type ConnectionType) {
    if settings.Node.ConnectionMode[pack.From.Hash] == models.CONN { 
        return
    }

    var connects []string
    if conn_type == _MERGE && mode != models.C_S_mode {
        connects = settings.MakeAddresses(settings.Node.Address.P2P)
    }

    public_key, err := hex.DecodeString(pack.Body)
    utils.CheckError(err)

    var public_data = string(public_key)

    // encoding.DecodePublic(string(public_key))
    settings.Mutex.Lock()
    settings.Node.PublicKey[pack.From.Hash] = encoding.DecodePublic(public_data) 
    settings.Node.SessionKey.P2P[pack.From.Hash] = crypto.SessionKey(settings.SESSION_KEY_BYTES)
    if mode == models.P2P_mode {
        settings.Node.Address.P2P[pack.From.Hash] = pack.From.Address
    }
    settings.Mutex.Unlock()

    // fmt.Println(settings.Node.SessionKey.P2P[pack.From.Hash])

    var encrypted_address = crypto.Encrypt(
        settings.Node.SessionKey.P2P[pack.From.Hash], 
        settings.User.IPv4 + settings.User.Port,
    )

    encrypted_session_key, err := crypto.EncryptRSA(
        settings.Node.SessionKey.P2P[pack.From.Hash],
        settings.Node.PublicKey[pack.From.Hash],
    )
    utils.CheckError(err)

    var new_pack = models.PackageTCP {
        From: models.From {
            Address: encrypted_address,
            Hash: settings.User.Hash.P2P,
        },
        Head: models.Head {
            Title: settings.HEAD_CONNECT,
            Mode: settings.MODE_SAVE,
        },
        Body: hex.EncodeToString(encrypted_session_key) + 
            settings.SEPARATOR + hex.EncodeToString([]byte(settings.User.Public.Data.P2P)) + 
            settings.SEPARATOR + strings.Join(connects, settings.SEPARATOR),
    }

    if conn_type == _MERGE && mode != models.C_S_mode {
        new_pack.Head.Mode = settings.MODE_SAVE_GLOBAL
    }

    var return_code int8
    if mode == models.C_S_mode {
        return_code = sendPackageByArchCS(settings.Node.Address.C_S[pack.From.Hash], new_pack)
    } else {
        return_code = sendPackageByAddr(pack.From.Address, new_pack)
    }

    if return_code == settings.EXIT_SUCCESS {
        settings.Mutex.Lock()
        settings.Node.ConnectionMode[pack.From.Hash] = models.CONN
        settings.Messages.NewDataExistLocal[pack.From.Hash] = make(chan bool)
        _, err = settings.DataBase.Exec(`
CREATE TABLE IF NOT EXISTS Local` + pack.From.Hash + ` (
Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
User VARCHAR(64),
Mode VARCHAR(3),
Body TEXT
);
INSERT INTO Connections(User, PublicKey)
SELECT '` + pack.From.Hash + `', '` + public_data + `'
WHERE NOT EXISTS(SELECT 1 FROM Connections WHERE User = '` + pack.From.Hash + `');
`)
        settings.Mutex.Unlock()
        utils.CheckError(err)
    } else {
        nullNode(pack.From.Hash)
    }
}

// Accept connection.
func connectSave(pack models.PackageTCP, mode models.ModeNet, conn_type ConnectionType) {
    if settings.Node.ConnectionMode[pack.From.Hash] == models.CONN {
        return
    }

    var (
        splited = strings.Split(pack.Body, settings.SEPARATOR)
        connected_nodes = settings.Node.Address.P2P
        connects = splited[2:]
    )

    if conn_type == _MERGE && mode != models.C_S_mode {
        go connectP2P(connects, true)
    }

    encrypted_session_key, err := hex.DecodeString(splited[0])
    utils.CheckError(err)

    session_key, err := crypto.DecryptRSA(
        encrypted_session_key,
        settings.User.Private.Key.P2P,
    )
    utils.CheckError(err)

    // fmt.Println(session_key)

    public_key, err := hex.DecodeString(splited[1])
    utils.CheckError(err)

    var public_data = string(public_key)
    var address string 

    if conn_type == _MERGE && mode != models.C_S_mode {
        address = crypto.Decrypt(session_key, pack.From.Address) 
        go redirectConnect(connected_nodes, append(connects, address))
    }

    settings.Mutex.Lock()
    settings.Node.PublicKey[pack.From.Hash] = encoding.DecodePublic(string(public_key))
    settings.Node.SessionKey.P2P[pack.From.Hash] = session_key
    if mode != models.C_S_mode {
        settings.Node.Address.P2P[pack.From.Hash] = address
    } else if mode == models.C_S_mode { // && settings.Node.ConnServer.Hash == ""
        settings.Node.ConnServer.Hash = pack.From.Hash
    }
    settings.Node.ConnectionMode[pack.From.Hash] = models.CONN
    settings.Messages.NewDataExistLocal[pack.From.Hash] = make(chan bool)
    _, err = settings.DataBase.Exec(`
CREATE TABLE IF NOT EXISTS Local` + pack.From.Hash + ` (
Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
User VARCHAR(64),
Mode VARCHAR(3),
Body TEXT
);
INSERT INTO Connections(User, PublicKey) 
SELECT '` + pack.From.Hash + `', '` + public_data + `'
WHERE NOT EXISTS(SELECT 1 FROM Connections WHERE User = '` + pack.From.Hash + `');
`)
    settings.Mutex.Unlock()
    utils.CheckError(err)
}

// Connection with sending public key.
func connectSaveLocal(pack models.PackageTCP) {
    if settings.Node.ConnectionMode[pack.From.Hash] == models.CONN { 
        return
    }
    if settings.User.Mode == models.C_S_mode {
        ConnectArchCS(settings.Node.ConnServer.Addr, false)
        return
    }
    connectP2P([]string{pack.From.Address}, false)
}

// Connection with found public key in database.
func connectReadLocal(pack models.PackageTCP, mode models.ModeNet) {
    if settings.Node.ConnectionMode[pack.From.Hash] == models.CONN { 
        return
    }
    var row = settings.DataBase.QueryRow(
        "SELECT PublicKey FROM Connections WHERE User = $1",
        pack.From.Hash,
    )

    var public_data string
    row.Scan(&public_data)

    if public_data == "" {
        var new_pack = models.PackageTCP {
            From: models.From {
                Hash: settings.User.Hash.P2P,
                Address: settings.User.IPv4 + settings.User.Port,
            },
            Head: models.Head {
                Title: settings.HEAD_CONNECT,
                Mode: settings.MODE_SAVE_LOCAL,
            },
        }
        if mode == models.C_S_mode {
            sendPackageByArchCS(settings.Node.Address.C_S[pack.From.Hash], new_pack)
        } else {
            sendPackageByAddr(pack.From.Address, new_pack)
        }
        return
    }

    settings.Mutex.Lock()
    settings.Node.PublicKey[pack.From.Hash] = encoding.DecodePublic(public_data)
    settings.Node.SessionKey.P2P[pack.From.Hash] = crypto.SessionKey(settings.SESSION_KEY_BYTES)
    if mode == models.P2P_mode {
        settings.Node.Address.P2P[pack.From.Hash] = pack.From.Address
    }
    settings.Mutex.Unlock()

    var encrypted_address = crypto.Encrypt(settings.Node.SessionKey.P2P[pack.From.Hash], settings.User.IPv4 + settings.User.Port)

    encrypted_session_key, err := crypto.EncryptRSA(
        settings.Node.SessionKey.P2P[pack.From.Hash],
        settings.Node.PublicKey[pack.From.Hash],
    )
    utils.CheckError(err)

    var new_pack = models.PackageTCP {
        From: models.From {
            Address: encrypted_address,
            Hash: settings.User.Hash.P2P,
        },
        Head: models.Head {
            Title: settings.HEAD_CONNECT,
            Mode: settings.MODE_SAVE,
        },
        Body: hex.EncodeToString(encrypted_session_key) + 
            settings.SEPARATOR + hex.EncodeToString([]byte(settings.User.Public.Data.P2P)),
    }

    var return_code int8
    if mode == models.C_S_mode {
        return_code = sendPackageByArchCS(settings.Node.Address.C_S[pack.From.Hash], new_pack)
    } else {
        return_code = sendPackageByAddr(pack.From.Address, new_pack)
    }

    if return_code == settings.EXIT_SUCCESS {
        settings.Mutex.Lock()
        settings.Node.ConnectionMode[pack.From.Hash] = models.CONN
        settings.Messages.NewDataExistLocal[pack.From.Hash] = make(chan bool)
        settings.Mutex.Unlock()
    } else {
        nullNode(pack.From.Hash)
    }
}

// Send alive connection confirmation.
func connectReadCheck(pack models.PackageTCP) {
    var new_pack = models.PackageTCP {
        From: models.From {
            Hash: settings.User.Hash.P2P,
        },
        To: models.To {
            Hash: pack.From.Hash,
        },
        Head: models.Head {
            Title: settings.HEAD_CONNECT,
            Mode: settings.MODE_SAVE_CHECK,
        },
    }
    SendEncryptedPackage(new_pack, models.P2P_mode)
}

// Connection exists.
func connectSaveCheck(pack models.PackageTCP) {
    __check_connection[pack.From.Hash] = true
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
            Hash: pack.From.Hash,
        },
        To: models.To {
            Hash: settings.User.Hash.F2F,
        },
        Head: models.Head {
            Title: heads[0],
            Mode: heads[1],
        },
        Body: pack.Body,
    }

    packageActions(new_pack, models.F2F_mode)
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
        pack.From.Hash,
        crypto.Encrypt(settings.User.Password, splited[2]),
    )
    settings.Mutex.Unlock()
    utils.CheckError(err)
}

// Send global message to all nodes.
func messageGlobal(pack models.PackageTCP, message string, mode models.ModeNet) {
    if mode == models.F2F_mode { 
        if _, ok := settings.Node.Address.F2F[pack.From.Hash]; !ok {
            message = "(HiddenFriend)" + message
        }
    }
    var mode_name = settings.GetStrMode(mode)
    settings.Mutex.Lock()
    _, err := settings.DataBase.Exec(
        "INSERT INTO GlobalMessages (User, Mode, Body) VALUES ($1, $2, $3)",
        pack.From.Hash, 
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
func messageLocal(pack models.PackageTCP, message string, mode models.ModeNet) {
    if mode == models.F2F_mode { 
        if _, ok := settings.Node.Address.F2F[pack.From.Hash]; !ok {
            messageGlobal(pack, message, models.F2F_mode)
            return
        }
    }
    var mode_name = settings.GetStrMode(mode)
    settings.Mutex.Lock()
    _, err := settings.DataBase.Exec(
        "INSERT INTO Local" + pack.From.Hash + " (User, Mode, Body) VALUES ($1, $2, $3)",
        pack.From.Hash,
        mode_name,
        crypto.Encrypt(settings.User.Password, message),
    )
    settings.Messages.CurrentIdLocal[pack.From.Hash]++
    settings.Mutex.Unlock()
    utils.CheckError(err)
    go func() {
        settings.Messages.NewDataExistLocal[pack.From.Hash] <- true
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
func archiveReadFile(pack models.PackageTCP, mode models.ModeNet) {
    if utils.FileIsExist(settings.PATH_ARCHIVE + pack.Body) && 
      !strings.Contains(pack.Body, "..") {
        var new_pack = models.PackageTCP {
            From: models.From {
                Hash: pack.To.Hash,
            },
            To: models.To {
                Hash: pack.From.Hash,
            },
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
func archiveReadList(pack models.PackageTCP, mode models.ModeNet) {
    files, err := ioutil.ReadDir(settings.PATH_ARCHIVE)
    utils.CheckError(err)
    var list = make([]string, len(files))
    for index, file := range files {
        list[index] = file.Name()
    }

    var list_of_files = strings.Join(list, settings.SEPARATOR)
    var new_pack = models.PackageTCP {
        From: models.From {
            Hash: pack.To.Hash,
        },
        To: models.To {
            Hash: pack.From.Hash,
        },
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

    fmt.Printf("| %s:\n", pack.From.Hash)
    for _, file := range settings.User.TempArchive {
        if file != "" {
            fmt.Println("|", file)
        }
    }
}
