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
    for {
        var message = utils.Input()
        var splited = strings.Split(message, " ")

        switch splited[0] {
            case settings.TERM_EXIT:
                os.Exit( settings.EXIT_SUCCESS )

            case settings.TERM_HELP:
                utils.PrintHelp()

            case settings.TERM_NETWORK:
                for _, addr := range settings.User.Connections {
                    fmt.Println("|", addr)
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

            case settings.TERM_ARCHIVE:
                if len(splited) == 1 {
                    files, err := ioutil.ReadDir(settings.PATH_ARCHIVE)
                    utils.CheckError(err)

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
                        for _, addr := range splited[2:] {
                            var new_pack = settings.Package {
                                From: models.From {
                                    Address: settings.User.IPv4 + settings.User.Port,
                                    Name: settings.User.Name,
                                },
                                To: addr,
                                Head: models.Head {
                                    Header: settings.HEAD_ARCHIVE,
                                    Mode: settings.MODE_GET_LIST,
                                }, 
                            }
                            SendEncryptedPackage(new_pack)
                            time.Sleep(time.Second * settings.TIME_SLEEP)

                            fmt.Printf("| %s:\n", addr)
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
                                    Address: settings.User.IPv4 + settings.User.Port,
                                    Name: settings.User.Name,
                                },
                                To: splited[2],
                                Head: models.Head {
                                    Header: settings.HEAD_ARCHIVE,
                                    Mode: settings.MODE_GET_FILE,
                                },
                                Body: filename,
                            }

                            SendEncryptedPackage(new_pack)
                            time.Sleep(time.Second * settings.TIME_SLEEP)
                        }
                }

            case settings.TERM_HISTORY:
                if len(splited) == 1 {
                    PrintGlobalHistory()
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
                            PrintLocalHistory(splited[2:])
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

                for _, addr := range settings.User.Connections {
                    var new_pack = settings.Package {
                        From: models.From {
                            Address: settings.User.IPv4 + settings.User.Port,
                            Name: settings.User.Name,
                        },
                        To: addr,
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

func DeleteGlobalMessages() {
    settings.Mutex.Lock()
    settings.User.GlobalMessages = []string{}
    settings.Mutex.Unlock()
}

func DeleteLocalMessages(slice []string) {
    settings.Mutex.Lock()
    for _, check := range slice {
        for _, addr := range settings.User.Connections {
            if addr == check {
                settings.User.LocalMessages[addr] = []string{}
                break
            }
        }
    }
    settings.Mutex.Unlock()
}

func PrintGlobalHistory() {
    for _, mes := range settings.User.GlobalMessages {
        fmt.Println("|", mes)
    }
}

func PrintLocalHistory(slice []string) {
    for _, check := range slice {
        for _, addr := range settings.User.Connections {
            if addr == check {
                fmt.Printf("| %s:\n", addr)
                for _, mes := range settings.User.LocalMessages[addr] {
                    fmt.Println("|", mes)
                }
                break
            }
        }
    }
}

func Connect(slice []string) {
    next:
    for _, addr := range slice {
        if addr == settings.User.IPv4 + settings.User.Port {
            continue
        }

        for _, value := range settings.User.Connections {
            if addr == value {
                continue next
            }
        }

        settings.Mutex.Lock()
        settings.User.Connections = append(
            settings.User.Connections, 
            addr,
        )
        settings.Mutex.Unlock()

        var new_pack = settings.Package {
            From: models.From {
                Address: settings.User.IPv4 + settings.User.Port,
                Name: settings.User.Name,
            },
            To: addr,
            Head: models.Head {
                Header: settings.HEAD_CONNECT,
                Mode: settings.MODE_GET,
            },
            Body: hex.EncodeToString([]byte(settings.User.PublicData)),
        }
        SendPackage(addr, new_pack)
    }
}

func Disconnect(slice []string) {
    settings.Mutex.Lock()
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
    settings.Mutex.Unlock()
}
