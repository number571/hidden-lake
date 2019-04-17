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

var set_email struct {
    title string
    body string
}

func ClientTCP() {
    var (
        message string
        splited []string
        authorization struct {
            login string
            password string
        }
    )

    for {
        message = utils.Input()
        splited = strings.Split(message, " ")

        switch splited[0] {
            case settings.TERM_EXIT:
                os.Exit( settings.EXIT_SUCCESS )

            case settings.TERM_HELP:
                utils.PrintHelp()

            case settings.TERM_INTERFACE:
                if settings.ServerListenHTTP == nil {
                    go ClientHTTP()
                } else {
                    if err := settings.ServerListenHTTP.Shutdown(context.TODO()); err != nil {
                        utils.PrintWarning("failure shutting down")
                    }
                }
        }

        if !settings.User.Auth {
            switch splited[0] {
                case settings.TERM_LOGIN:
                    if len(splited) > 1 {
                        authorization.login = strings.Join(splited[1:], " ")
                    }

                case settings.TERM_PASSWORD:
                    if len(splited) > 1 {
                        authorization.password = strings.Join(splited[1:], " ")
                    }

                case settings.TERM_ADDRESS:
                    if len(splited) > 1 {
                        var ipv4_port = strings.Split(splited[1], ":")
                        if len(ipv4_port) != 2 {
                            utils.PrintWarning("invalid argument for ':address'")
                            continue
                        } 
                        settings.Mutex.Lock()
                        settings.User.IPv4 = ipv4_port[0]
                        settings.User.Port = ":" + ipv4_port[1]
                        settings.Mutex.Unlock()
                    }

                case settings.TERM_ENTER:
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
        }

        if settings.User.Auth {
            client(splited, message)
        }
    }
}

func client(splited []string, message string) {
    switch splited[0] {
        case settings.TERM_WHOAMI:
            fmt.Println("|", settings.User.Hash)

        case settings.TERM_LOGOUT:
            connect.Logout()

        case settings.TERM_NETWORK:
            for _, username := range settings.User.Connections {
                fmt.Println("|", username)
            }

        case settings.TERM_SEND:
            if len(splited) > 2 {
                var new_pack = settings.PackageTCP {
                    From: models.From {
                        Address: settings.User.IPv4 + settings.User.Port,
                        Name: settings.User.Hash,
                    },
                    To: splited[1],
                    Head: models.Head {
                        Header: settings.HEAD_MESSAGE,
                        Mode: settings.MODE_LOCAL,
                    }, 
                    Body: strings.Join(splited[2:], " "),
                }
                connect.SendEncryptedPackage(new_pack)
            }

        case settings.TERM_EMAIL:
            var length = len(splited)
            if length > 1 {
                switch splited[1] {
                    case "title": 
                        if length > 2 {
                            set_email.title = strings.Join(splited[2:], " ")
                        }
                    case "body": 
                        if length > 2 {
                            set_email.body = strings.Join(splited[2:], " ")
                        }
                    case "write": 
                        if length == 3 {
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
                            connect.SendEncryptedPackage(new_pack)
                        }
                    case "read": 
                        switch length {
                            case 2: 
                                var email models.Email
                                rows, err := settings.DataBase.Query("SELECT Id, Title, User, Date FROM Email")
                                utils.CheckError(err)

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
                                rows.Close()

                            case 3:
                                var email models.Email
                                rows, err := settings.DataBase.Query(
                                    "SELECT Id, Title, User, Date FROM Email WHERE User=$1", 
                                    splited[2],
                                )
                                utils.CheckError(err)

                                for rows.Next() {
                                    err = rows.Scan(
                                        &email.Id,
                                        &email.Title,
                                        &email.User,
                                        &email.Date,
                                    )
                                    utils.CheckError(err)
                                    crypto.DecryptEmail(settings.User.Password, &email)
                                    fmt.Println("|", email.Id, "|", email.Title, "|", email.User, "|", email.Date)
                                }
                                rows.Close()

                            case 4:
                                var email models.Email
                                rows, err := settings.DataBase.Query(
                                    "SELECT * FROM Email WHERE User=$1 AND Id=$2", 
                                    splited[2], 
                                    splited[3],
                                )
                                utils.CheckError(err)

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
                                rows.Close()
                        }
                }
            }

        case settings.TERM_ARCHIVE:
            if len(splited) == 1 {
                files, err := ioutil.ReadDir(settings.PATH_ARCHIVE)
                utils.CheckError(err)

                fmt.Printf("| %s:\n", settings.User.Hash)
                for _, file := range files {
                    fmt.Println("|", file.Name())
                }
                return
            }

            if len(splited) == 2 { 
                return 
            }

            switch splited[1] {
                case "list": 
                    for _, name := range splited[2:] {
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
                        connect.SendEncryptedPackage(new_pack)
                        time.Sleep(time.Second * settings.TIME_SLEEP)

                        fmt.Printf("| %s:\n", name)
                        for _, file := range settings.User.TempArchive {
                            if file != "" {
                                fmt.Println("|", file)
                            }
                        }
                    }

                case "download": 
                    if len(splited) < 4 {
                        return
                    }
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
                        connect.SendEncryptedPackage(new_pack)
                        time.Sleep(time.Second * settings.TIME_SLEEP)
                    }
            }

        case settings.TERM_HISTORY:
            if len(splited) == 1 {
                printGlobalHistory()
                return
            }
            switch splited[1] {
                case "del", "delete": 
                    if len(splited) == 2 {
                        connect.DeleteGlobalMessages()
                        return
                    }
                    connect.DeleteLocalMessages(splited[2:])
                            
                case "loc", "local":
                    if len(splited) > 2 {
                        printLocalHistory(splited[2:])
                    }
            }

        case settings.TERM_CONNECT:
            if len(splited) > 1 {
                connect.Connect(splited[1:])
            }

        case settings.TERM_DISCONNECT:
            if len(splited) == 1 {
                connect.Disconnect(settings.User.Connections)
                return
            }
            connect.Disconnect(splited[1:])

        default:
            if message == "" { 
                return 
            }
            for _, username := range settings.User.Connections {
                var new_pack = settings.PackageTCP {
                    From: models.From {
                        Address: settings.User.IPv4 + settings.User.Port,
                        Name: settings.User.Hash,
                    },
                    To: username,
                    Head: models.Head {
                        Header: settings.HEAD_MESSAGE,
                        Mode: settings.MODE_GLOBAL,
                    },
                    Body: message,
                }
                connect.SendEncryptedPackage(new_pack)
            }
    }
}

func printGlobalHistory() {
    settings.Mutex.Lock()
    rows, err := settings.DataBase.Query("SELECT Body FROM GlobalMessages ORDER BY Id")
    settings.Mutex.Unlock()

    utils.CheckError(err)

    var data string

    for rows.Next() {
        rows.Scan(&data)
        fmt.Println("|", data)
    }

    rows.Close()
}

func printLocalHistory(slice []string) {
    for _, user := range slice {
        for _, check := range settings.User.Connections {
            if check == user {
                settings.Mutex.Lock()
                rows, err := settings.DataBase.Query("SELECT Body FROM Local" + user + " WHERE ORDER BY Id")
                settings.Mutex.Unlock()

                utils.CheckError(err)

                fmt.Printf("| %s:\n", user)
                var data string

                for rows.Next() {
                    rows.Scan(&data)
                    fmt.Println("|", data)
                }

                rows.Close()
                break
            }
        }
    }
}
