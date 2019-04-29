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
            case settings.TERM_EXIT: os.Exit(settings.EXIT_SUCCESS)
            case settings.TERM_HELP: utils.PrintHelp()
            case settings.TERM_INTERFACE: turnInterface()
        }

        if !settings.User.Auth {
            switch splited[0] {
                case settings.TERM_LOGIN: setLogin(&authorization, splited)
                case settings.TERM_PASSWORD: setPassword(&authorization, splited)
                case settings.TERM_ADDRESS: setAddress(splited)
                case settings.TERM_ENTER: pressEnter(authorization)
            }
        } else {
            client(splited, message)
        }
    }
}

func client(splited []string, message string) {
    switch splited[0] {
        case settings.TERM_WHOAMI: fmt.Println("|", settings.User.Hash)
        case settings.TERM_LOGOUT: connect.Logout()
        case settings.TERM_NETWORK: network()
        case settings.TERM_SEND: sendLocalMessage(splited)
        case settings.TERM_EMAIL: emailAction(splited)
        case settings.TERM_ARCHIVE: archiveAction(splited)
        case settings.TERM_HISTORY: historyAction(splited)
        case settings.TERM_CONNECT: connectTo(splited)
        default: sendGlobalMessage(message)
    }
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
    for _, filename := range splited[3:] {
        var new_pack = settings.PackageTCP {
            From: models.From {
                Name: settings.User.Hash,
            },
            To: splited[2],
            Head: models.Head {
                Header: settings.HEAD_ARCHIVE,
                Mode: settings.MODE_READ_FILE,
            },
            Body: filename,
        }
        connect.CreateRedirectPackage(&new_pack)
        connect.SendInitRedirectPackage(new_pack)
        time.Sleep(time.Second * settings.TIME_SLEEP) // FIX
    }
}

// Print list of files in nodes archive.
func listNodeArchive(splited []string) {
    for _, name := range splited[1:] {
        var new_pack = settings.PackageTCP {
            From: models.From {
                Name: settings.User.Hash,
            },
            To: name,
            Head: models.Head {
                Header: settings.HEAD_ARCHIVE,
                Mode: settings.MODE_READ_LIST,
            }, 
        }
        connect.CreateRedirectPackage(&new_pack)
        connect.SendInitRedirectPackage(new_pack)
        time.Sleep(time.Second * settings.TIME_SLEEP) // FIX
        fmt.Printf("| %s:\n", name)
        for _, file := range settings.User.TempArchive {
            if file != "" {
                fmt.Println("|", file)
            }
        }
    }
}

// Print list of files in archive.
func listArchive() {
    files, err := ioutil.ReadDir(settings.PATH_ARCHIVE)
    utils.CheckError(err)

    fmt.Printf("| %s:\n", settings.User.Hash)
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
        connect.DeleteGlobalMessages()
        return
    }
    connect.DeleteLocalMessages(splited[2:])
}

// Print local messages.
func historyLocal(splited []string, length int) {
    if length > 2 {
        printLocalHistory(splited[2:])
    }
}

// Connect to nodes.
func connectTo(splited []string) {
    if len(splited) > 1 {
        connect.Connect(splited[1:], false)
    }
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
    var new_pack = settings.PackageTCP {
        From: models.From {
            Name: settings.User.Hash,
        },
        To: splited[2],
        Head: models.Head {
            Header: settings.HEAD_EMAIL,
            Mode: settings.MODE_SAVE,
        }, 
        Body: 
            set_email.title + settings.SEPARATOR +
            set_email.body + settings.SEPARATOR +
            time.Now().Format(time.RFC850),
    }
    connect.CreateRedirectPackage(&new_pack)
    connect.SendInitRedirectPackage(new_pack)
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
    for username := range settings.User.NodeAddress {
        var new_pack = settings.PackageTCP {
            From: models.From {
                Name: settings.User.Hash,
            },
            To: username,
            Head: models.Head {
                Header: settings.HEAD_MESSAGE,
                Mode: settings.MODE_GLOBAL,
            },
            Body: message,
        }
        connect.CreateRedirectPackage(&new_pack)
        connect.SendInitRedirectPackage(new_pack)
    }
}

// Send local message to one node.
func sendLocalMessage(splited []string) {
    if len(splited) > 2 {
        var new_pack = settings.PackageTCP {
            From: models.From {
                Name: settings.User.Hash,
            },
            To: splited[1],
            Head: models.Head {
                Header: settings.HEAD_MESSAGE,
                Mode: settings.MODE_LOCAL,
            }, 
            Body: strings.Join(splited[2:], " "),
        }
        connect.CreateRedirectPackage(&new_pack)
        connect.SendInitRedirectPackage(new_pack)
    }
}

// Print connections.
func network() {
    for username := range settings.User.NodeAddress {
        fmt.Println("|", username)
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
                go connect.FindConnects(10)
            }
            fmt.Println("[SUCCESS]: Authorization")
    }
}

// Turn on/off interface.
func turnInterface() {
    if settings.ServerListenHTTP == nil {
        go ClientHTTP()
    } else {
        if err := settings.ServerListenHTTP.Shutdown(context.TODO()); err != nil {
            utils.PrintWarning("failure shutting down")
        }
    }
}

// Set address ipv4:port.
func setAddress(splited []string) {
    if len(splited) > 1 {
        var ipv4_port = strings.Split(splited[1], ":")
        if len(ipv4_port) != 2 {
            utils.PrintWarning("invalid argument for ':address'")
            return
        } 
        settings.Mutex.Lock()
        settings.User.IPv4 = ipv4_port[0]
        settings.User.Port = ":" + ipv4_port[1]
        settings.Mutex.Unlock()
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
    for _, user := range slice {
        if _, ok := settings.User.NodeAddress[user]; ok {
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
