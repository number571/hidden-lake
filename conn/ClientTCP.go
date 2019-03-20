package conn

import (
    "os"
    "fmt"
    "time"
    "strings"
    "io/ioutil"
    "encoding/hex"
    "../utils"
    "../models"
    "../settings"
)

func ClientTCP() {
    Connect(settings.User.WhiteList)

    go findConnects(10)

    for {
        var message = utils.Input()
        var splited = strings.Split(message, " ")

        switch splited[0] {
            case settings.TERM_WHOAMI:
                fmt.Println("|", settings.User.Name)

            case settings.TERM_EXIT:
                os.Exit( settings.EXIT_SUCCESS )

            case settings.TERM_HELP:
                utils.PrintHelp()

            case settings.TERM_NETWORK:
                for _, username := range settings.User.Connections {
                    fmt.Println("|", username)
                }

            case settings.TERM_SEND:
                if len(splited) > 2 {
                    var new_pack = settings.Package {
                        From: models.From {
                            Address: settings.User.IPv4 + settings.User.Port,
                            Name: settings.User.Name,
                        },
                        To: splited[1],
                        Head: models.Head {
                            Header: settings.HEAD_MESSAGE,
                            Mode: settings.MODE_LOCAL,
                        }, 
                        Body: strings.Join(splited[2:], " "),
                    }
                    SendEncryptedPackage(new_pack)
                }

            case settings.TERM_REFRESH:
                Connect(settings.User.WhiteList)

            case settings.TERM_ARCHIVE:
                if len(splited) == 1 {
                    files, err := ioutil.ReadDir(settings.PATH_ARCHIVE)
                    utils.CheckError(err)

                    fmt.Printf("| %s:\n", settings.User.Name)
                    for _, file := range files {
                        fmt.Println("|", file.Name())
                    }
                    continue
                }

                if len(splited) < 3 {
                    continue
                }

                switch splited[1] {
                    case "list": 
                        for _, name := range splited[2:] {
                            var new_pack = settings.Package {
                                From: models.From {
                                    Name: settings.User.Name,
                                },
                                To: name,
                                Head: models.Head {
                                    Header: settings.HEAD_ARCHIVE,
                                    Mode: settings.MODE_READ_LIST,
                                }, 
                            }
                            SendEncryptedPackage(new_pack)
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
                            continue
                        }
                        for _, filename := range splited[3:] {
                            var new_pack = settings.Package {
                                From: models.From {
                                    Name: settings.User.Name,
                                },
                                To: splited[2],
                                Head: models.Head {
                                    Header: settings.HEAD_ARCHIVE,
                                    Mode: settings.MODE_READ_FILE,
                                },
                                Body: filename,
                            }
                            SendEncryptedPackage(new_pack)
                            time.Sleep(time.Second * settings.TIME_SLEEP)
                        }
                }

            case settings.TERM_HISTORY:
                if len(splited) == 1 {
                    printGlobalHistory()
                    continue
                }
                switch splited[1] {
                    case "del", "delete": 
                        if len(splited) == 2 {
                            DeleteGlobalMessages()
                            return
                        }
                        DeleteLocalMessages(splited[2:])
                                
                    case "loc", "local":
                        if len(splited) > 2 {
                            printLocalHistory(splited[2:])
                        }
                }

            case settings.TERM_CONNECT:
                if len(splited) > 1 {
                    Connect(splited[1:])
                }

            case settings.TERM_DISCONNECT:
                if len(splited) == 1 {
                    Disconnect(settings.User.Connections)
                    continue
                }
                Disconnect(splited[1:])

            default:
                if message == "" {
                    continue
                }

                for _, username := range settings.User.Connections {
                    var new_pack = settings.Package {
                        From: models.From {
                            Address: settings.User.IPv4 + settings.User.Port,
                            Name: settings.User.Name,
                        },
                        To: username,
                        Head: models.Head {
                            Header: settings.HEAD_MESSAGE,
                            Mode: settings.MODE_GLOBAL,
                        },
                        Body: message,
                    }
                    SendEncryptedPackage(new_pack)
                }
        }
    }
}

func findConnects(seconds time.Duration) {
    for {
        for _, username := range settings.User.Connections {
            var new_pack = settings.Package {
                From: models.From {
                    Address: settings.User.IPv4 + settings.User.Port,
                    Name: settings.User.Name,
                },
                To: username,
                Head: models.Head {
                    Header: settings.HEAD_CONNECT,
                    Mode: settings.MODE_READ_LIST,
                },
            }
            SendEncryptedPackage(new_pack)
        }
        time.Sleep(seconds * time.Second)
    }
}

func printGlobalHistory() {
    for _, mes := range settings.User.GlobalMessages {
        fmt.Println("|", mes)
    }
}

func printLocalHistory(slice []string) {
    for _, check := range slice {
        for _, username := range settings.User.Connections {
            if username == check {
                fmt.Printf("| %s:\n", username)
                for _, mes := range settings.User.LocalMessages[username] {
                    fmt.Println("|", mes)
                }
                break
            }
        }
    }
}

func DeleteGlobalMessages() {
    settings.Mutex.Lock()
    settings.User.GlobalMessages = []string{}
    settings.Mutex.Unlock()
}

func DeleteLocalMessages(slice []string) {
    settings.Mutex.Lock()
    for _, check := range slice {
        for _, username := range settings.User.Connections {
            if username == check {
                settings.User.LocalMessages[username] = []string{}
                break
            }
        }
    }
    settings.Mutex.Unlock()
}

func Connect(slice []string) {
    next:
    for _, addr := range slice {
        var address = settings.User.IPv4 + settings.User.Port

        if addr == address {
            continue
        }

        for _, username := range settings.User.Connections {
            if addr == settings.User.NodeAddress[username] {
                continue next
            }
        }

        var new_pack = settings.Package {
            From: models.From {
                Address: address,
                Name: settings.User.Name,
            },
            Head: models.Head {
                Header: settings.HEAD_CONNECT,
                Mode: settings.MODE_READ,
            },
            Body: hex.EncodeToString([]byte(settings.User.PublicData)),
        }

        sendAddrPackage(addr, new_pack)
    }
}

func Disconnect(slice []string) {
    for _, addr := range slice {
        fmt.Println("|", addr)
        var new_pack = settings.Package {
            From: models.From {
                Address: settings.User.IPv4 + settings.User.Port,
                Name: settings.User.Name,
            },
            To: addr,
            Head: models.Head {
                Header: settings.HEAD_WARNING,
                Mode: settings.MODE_SAVE,
            },
        }
        SendEncryptedPackage(new_pack)
        nullNode(addr)
    }
}
