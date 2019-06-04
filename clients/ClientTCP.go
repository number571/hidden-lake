package clients

import (
    "os"
    "fmt"
    "time"
    "context"
    "strings"
    "io/ioutil"
    "../utils"
    "../crypto"
    "../models"
    "../connect"
    "../settings"
)

type auth struct {
    login string
    password string
}

var set_email struct {
    title string
    body string
}

func ClientTCP() {
    var (
        message string
        splited []string
        authorization auth
    )

    for {
        message = utils.Input()
        splited = strings.Split(message, " ")

        switch splited[0] {
            case settings.TERM_MODE: switchModeNet(); continue
            case settings.TERM_HELP: utils.PrintHelp(); continue
            case settings.TERM_INTERFACE: turnInterface(); continue
            case settings.TERM_EXIT: os.Exit(settings.EXIT_SUCCESS)
        }

        if !settings.User.Auth {
            switch splited[0] {
                case settings.TERM_LOGIN: setLogin(&authorization, splited)
                case settings.TERM_PASSWORD: setPassword(&authorization, splited)
                case settings.TERM_ADDRESS: setAddress(splited)
                case settings.TERM_ENTER: pressEnter(authorization)
            }
        } else {
            switch splited[0] {
                case settings.TERM_WHOAMI: whoami()
                case settings.TERM_LOGOUT: connect.Logout()
                case settings.TERM_NETWORK: network(splited)
                case settings.TERM_SEND: sendLocalMessage(splited)
                case settings.TERM_EMAIL: emailAction(splited)
                case settings.TERM_ARCHIVE: archiveAction(splited)
                case settings.TERM_HISTORY: historyAction(splited)
                case settings.TERM_CONNECT: connectTo(splited)
                case settings.TERM_DISCONNECT: disconnectFrom(splited)
                default: sendGlobalMessage(message)
            }
        }
    }
}

func whoami() {
    fmt.Println("| Hash:", settings.CurrentHash())
    fmt.Println("| Mode:", settings.CurrentMode())
}

// Actions with archives.
func archiveAction(splited []string) {
    switch len(splited) {
        case 1: listArchive()
        case 2: listNodeArchive(splited)
        case 3:
            switch splited[1] {
                case "download": downloadNodeFiles(splited)
            }
    }
}

// Download files from node archive.
func downloadNodeFiles(splited []string) {
    if len(splited) < 4 { return }
    var hashname = settings.CurrentHash()
    for _, filename := range splited[3:] {
        var new_pack = models.PackageTCP {
            From: models.From {
                Hash: hashname,
            },
            To: models.To {
                Hash: splited[2],
            },
            Head: models.Head {
                Title: settings.HEAD_ARCHIVE,
                Mode: settings.MODE_READ_FILE,
            },
            Body: filename,
        }
        connect.SendPackage(new_pack, settings.User.Mode)
    }
}

// Print list of files in nodes archive.
func listNodeArchive(splited []string) {
    var hashname = settings.CurrentHash()
    for _, name := range splited[1:] {
        var new_pack = models.PackageTCP {
            From: models.From {
                Hash: hashname,
            },
            To: models.To {
                Hash: name,
            },
            Head: models.Head {
                Title: settings.HEAD_ARCHIVE,
                Mode: settings.MODE_READ_LIST,
            }, 
        }
        connect.SendPackage(new_pack, settings.User.Mode)
    }
}

// Print list of files in archive.
func listArchive() {
    files, err := ioutil.ReadDir(settings.PATH_ARCHIVE)
    utils.CheckError(err)

    fmt.Printf("| %s:\n", settings.CurrentHash())
    for _, file := range files {
        fmt.Println("|", file.Name())
    }
}

// Actions with history of messages.
func historyAction(splited []string) {
    var length = len(splited)
    if length == 1 {
        printGlobalHistory()
        return
    }
    switch splited[1] {
        case "del", "delete": historyDelete(splited, length)
        case "loc", "local": historyLocal(splited, length)
    }
}

// Delete global or local messages.
func historyDelete(splited []string, length int) {
    if length == 2 {
        settings.DeleteGlobalMessages()
        return
    }
    settings.DeleteLocalMessages(splited[2:])
}

// Print local messages.
func historyLocal(splited []string, length int) {
    if length > 2 {
        printLocalHistory(splited[2:])
    }
}

// Connect to nodes.
func connectTo(splited []string) {
    switch settings.User.Mode {
        case models.P2P_mode: 
            if len(splited) < 2 { return }
            connect.ConnectP2PMerge(splited[1])
        case models.F2F_mode:
            if len(splited) < 4 { return }
            connect.ConnectF2F(splited[1], splited[2], strings.Join(splited[3:], " "))
        case models.C_S_mode:
            if len(splited) < 2 { return }
            settings.Node.ConnServer.Addr = connect.GetConnection(splited[1])
            connect.ConnectArchCS(settings.Node.ConnServer.Addr, true)
    }
}

// Disconnect from nodes.
func disconnectFrom(splited []string) {
    if len(splited) < 2 { return }
    connect.DisconnectF2F(splited[1])
}

// Actions with email.
func emailAction(splited []string) {
    var length = len(splited)
    if length > 1 {
        switch splited[1] {
            case "title": emailSetTitle(splited, length)
            case "body": emailSetBody(splited, length)
            case "write": emailWrite(splited, length)
            case "read": emailRead(splited, length)
            case "print": emailPrint(splited, length)
        }
    }
}

// Send email to one node. 
func emailWrite(splited []string, length int) {
    if length != 3 { return }
    var new_pack = models.PackageTCP {
        From: models.From {
            Hash: settings.CurrentHash(),
        },
        To: models.To {
            Hash: splited[2],
        },
        Head: models.Head {
            Title: settings.HEAD_EMAIL,
            Mode: settings.MODE_SAVE,
        }, 
        Body: 
            set_email.title + settings.SEPARATOR +
            set_email.body + settings.SEPARATOR +
            time.Now().Format(time.RFC850),
    }
    connect.SendPackage(new_pack, settings.User.Mode)
}

// Read email.
func emailRead(splited []string, length int) {
    switch length {
        case 2: emailReadAll(splited)
        case 3: emailReadAllByUser(splited)
        case 4: emailReadByUserAndId(splited)
    }
}

// Read list of emails by all nodes.
func emailReadAll(splited []string) {
    var (
        email models.Email
        err error
    )

    rows, err := settings.DataBase.Query("SELECT Id, Title, User, Date FROM Email")
    utils.CheckError(err)
    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(
            &email.Id,
            &email.Title,
            &email.User,
            &email.Date,
        )
        utils.CheckError(err)
        crypto.DecryptEmail(settings.User.Password, &email)
        fmt.Println("|", email.Id, "|", email.Title, "|", email.User, "|", email.Date, "|")
    }
}

// Read list of emails by one node.
func emailReadAllByUser(splited []string) {
    var (
        email models.Email
        err error
    )

    rows, err := settings.DataBase.Query(
        "SELECT Id, Title, User, Date FROM Email WHERE User=$1", 
        splited[2],
    )
    utils.CheckError(err)
    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(
            &email.Id,
            &email.Title,
            &email.User,
            &email.Date,
        )
        utils.CheckError(err)
        crypto.DecryptEmail(settings.User.Password, &email)
        fmt.Println("|", email.Id, "|", email.Title, "|", email.User, "|", email.Date, "|")
    }
}

// Read selected email by user and id.
func emailReadByUserAndId(splited []string) {
    var (
        email models.Email
        err error
    )

    rows, err := settings.DataBase.Query(
        "SELECT * FROM Email WHERE User=$1 AND Id=$2", 
        splited[2], 
        splited[3],
    )
    utils.CheckError(err)
    defer rows.Close()

    for rows.Next() {
        err = rows.Scan(
            &email.Id,
            &email.Title,
            &email.Body,
            &email.User,
            &email.Date,
        )
        utils.CheckError(err)
        crypto.DecryptEmail(settings.User.Password, &email)
        fmt.Println("--------------------------")
        fmt.Println("| Title:", email.Title, "|")
        fmt.Println("--------------------------")
        fmt.Println("| Body:", email.Body, "|")
        fmt.Println("--------------------------")
        fmt.Println("| Author:", email.User, "|")
        fmt.Println("--------------------------")
        fmt.Println("| Date:", email.Date, "|")
        fmt.Println("--------------------------")
    }
}

// Print selected emails data.
func emailPrint(splited []string, length int) {
    if length == 2 {
        fmt.Println("| Title:", set_email.title, "|")
        fmt.Println("| Body:", set_email.body, "|")
        return
    }

    switch splited[2] {
        case "title": fmt.Println("| Title:", set_email.title, "|")
        case "body": fmt.Println("| Body:", set_email.body, "|")
    }
}

// Set title in email.
func emailSetTitle(splited []string, length int) {
    if length > 2 {
        set_email.title = strings.Join(splited[2:], " ")
    }
}

// Set main text in email.
func emailSetBody(splited []string, length int) {
    if length > 2 {
        set_email.body = strings.Join(splited[2:], " ")
    }
}

// Send global message to all nodes.
func sendGlobalMessage(message string) {
    if message == "" { return }
    var new_pack = models.PackageTCP {
        From: models.From {
            Hash: settings.CurrentHash(),
        },
        Head: models.Head {
            Title: settings.HEAD_MESSAGE,
            Mode: settings.MODE_GLOBAL,
        },
        Body: message,
    }

    // switch settings.User.Mode {
    //     case models.P2P_mode:
    //         // for username := range settings.Node.Address.P2P {
    //         //     new_pack.To.Hash = username
    //         //     connect.SendEncryptedPackage(new_pack, models.P2P_mode)
    //         // }
    //         // for username := range settings.Node.Address.C_S {
    //         //     new_pack.To.Hash = username
    //         //     connect.SendEncryptedPackage(new_pack, models.C_S_mode)
    //         // }
    //         connect.SendPackage(new_pack, models.P2P_mode)
    //     case models.F2F_mode:
    //         connect.SendPackage(new_pack, models.F2F_mode)
    //         // connect.CreateRedirectF2FPackage(&new_pack, "")
    //         // for username := range settings.Node.Address.F2F {
    //         //     new_pack.To.Hash = username
    //         //     connect.SendEncryptedPackage(new_pack, models.F2F_mode)
    //         // }
    //     case models.C_S_mode:
    //         // if settings.Node.ConnServer.Addr == nil { return }
    //         connect.SendPackage(new_pack, models.C_S_mode)
    // }

    connect.SendPackage(new_pack, settings.User.Mode)
}

// Send local message to one node.
func sendLocalMessage(splited []string) {
    if len(splited) < 3 { return }

    var hashname = settings.CurrentHash()
    if splited[1] == hashname { return }

    var new_pack = models.PackageTCP {
        From: models.From {
            Hash: hashname,
        },
        To: models.To {
            Hash: splited[1],
        },
        Head: models.Head {
            Title: settings.HEAD_MESSAGE,
            Mode: settings.MODE_LOCAL,
        }, 
        Body: strings.Join(splited[2:], " "),
    }

    // C-S
    if settings.User.Mode == models.C_S_mode {
        // if settings.Node.ConnServer.Addr == nil { return }
        connect.SendPackage(new_pack, models.C_S_mode)
        return
    }

    // F2F
    if settings.User.Mode == models.F2F_mode {
        connect.SendPackage(new_pack, models.F2F_mode)
        // connect.CreateRedirectF2FPackage(&new_pack, splited[1])
        // for username := range settings.Node.Address.F2F {
        //     new_pack.To.Hash = username
        //     connect.SendPackage(new_pack, models.F2F_mode)
        // }
        return
    }

    // P2P
    var hashnames = strings.Split(splited[1], "->")
    if len(hashnames) < 2 {
        // hashnames = append(hashnames, "")
        // var hashname = pack.To.Hash
        // if hashname == "" { to = settings.Node.ConnServer.Hash }
    } else {
        new_pack.From.Address = hashnames[1]
        // new_pack.Body = crypto.Encrypt(
        //     settings.Node.SessionKey.P2P[hashnames[1]], 
        //     new_pack.Body,
        // )
    }

    new_pack.To.Hash = hashnames[0]
    connect.SendPackage(new_pack, models.P2P_mode)
}

// Print connections.
func network(splited []string) {
    if settings.User.Mode == models.C_S_mode {
        fmt.Println("|", settings.Node.ConnServer.Hash)
        return
    }

    if len(splited) < 2 { 
        printAllConnections() 
        return
    }

    switch splited[1] {
        case "p2p", "P2P": printP2PConnections()
        case "f2f", "F2F": printF2FConnections()
        case "c-s", "C-S": printArchCSConnections()
        default: printAllConnections()
    }
}

func printAllConnections() {
    printF2FConnections()
    printP2PConnections()
    printArchCSConnections()
}

func printF2FConnections() {
    fmt.Println("| F2F connections:")
    for username := range settings.Node.Address.F2F {
        fmt.Println("| -", username)
    }
}

func printArchCSConnections() {
    fmt.Println("| C-S connections:")
    for username := range settings.Node.Address.C_S {
        fmt.Println("| -", username)
    }
}

func printP2PConnections() {
    fmt.Println("| P2P connections:")
    for username := range settings.Node.Address.P2P {
        fmt.Println("| -", username)
    }
}

// Try to log in from login/password
func pressEnter(authorization auth) {
    switch settings.Authorization(authorization.login, authorization.password) {
        case 1: utils.PrintWarning("login is undefined")
        case 2: utils.PrintWarning("length of login > 64 bytes")
        case 3: utils.PrintWarning("password.hash undefined")
        case 4: utils.PrintWarning("login or password is wrong")
        default: 
            if !settings.GoroutinesIsRun && settings.User.Port != "" {
                settings.Mutex.Lock()
                settings.GoroutinesIsRun = true
                settings.Mutex.Unlock()
                go connect.ServerTCP()
                go connect.CheckConnects()
            }
            fmt.Println("[SUCCESS]: Authorization")
    }
}

// Switch P2P/F2F.
func switchModeNet() {
    if settings.User.Mode == models.C_S_mode { 
        goto print_mode 
    }

    settings.Mutex.Lock()
    if settings.User.Mode == models.P2P_mode {
        settings.User.Mode = models.F2F_mode
    } else {
        settings.User.Mode = models.P2P_mode
    }
    settings.Mutex.Unlock()

print_mode:
    fmt.Println("| Mode:", settings.CurrentMode())
}

// Turn on/off interface.
func turnInterface() {
    var mode string
    if settings.ServerListenHTTP == nil {
        go ClientHTTP()
        mode = "on"
    } else {
        if err := settings.ServerListenHTTP.Shutdown(context.TODO()); err != nil {
            utils.PrintWarning("failure shutting down")
        }
        mode = "off"
    }
    fmt.Println("| Interface:", mode)
}

// Set address ipv4:port.
func setAddress(splited []string) {
    if len(splited) < 2 { 
        fmt.Println("| Address:", settings.User.IPv4 + settings.User.Port)
        return
    }
    switch splited[1] {
        case "del", "delete": settings.ClearAddress()
        default: 
            var ipv4_port = strings.Split(splited[1], ":")
            if len(ipv4_port) != 2 {
                utils.PrintWarning("invalid argument for ':address'")
                return
            }
            settings.SetAddress(ipv4_port[0], ipv4_port[1])
    }
}

// Set login.
func setLogin(authorization *auth, splited []string) {
    if len(splited) > 1 {
        authorization.login = strings.Join(splited[1:], " ")
    }
}

// Set password.
func setPassword(authorization *auth, splited []string) {
    if len(splited) > 1 {
        authorization.password = strings.Join(splited[1:], " ")
    }
}

// Print messages from all nodes.
func printGlobalHistory() {
    rows, err := settings.DataBase.Query("SELECT Body FROM GlobalMessages ORDER BY Id")
    utils.CheckError(err)

    var data string

    for rows.Next() {
        rows.Scan(&data)
        fmt.Println("|", data)
    }

    rows.Close()
}

// Print local messages from nodes.
func printLocalHistory(slice []string) {
    var node_address = settings.CurrentNodeAddress()
    for _, user := range slice {
        if _, ok := node_address[user]; ok {
            rows, err := settings.DataBase.Query("SELECT Body FROM Local" + user + " WHERE ORDER BY Id")
            utils.CheckError(err)

            fmt.Printf("| %s:\n", user)
            var data string

            for rows.Next() {
                rows.Scan(&data)
                fmt.Println("|", data)
            }

            rows.Close()
        }
    }
}
