package settings

import (
    "os"
    "fmt"
    "strings"
    "io/ioutil"
    "../utils"
    "../crypto"
)

func Initialization(args []string) bool {
    var (
        flag_address bool
        flag_interface bool
        flag_login bool
        flag_password bool
        flag_delete bool
        flag_help bool
        auth struct {
            login string
            password string
        }
    )

    for _, value := range args[1:] {
        switch value {
            // Arguments without parameters
            case "-i", "--interface": 
                flag_interface = true
                continue
            case "-f", "--f2f":
                User.ModeF2F = true
                continue
            case "-h", "--help":
                if !flag_help {
                    utils.PrintHelp()
                    flag_help = true
                }
                continue
            case "-d", "--delete":
                if !flag_delete {
                    deleteAllData()
                    flag_delete = true
                }
                continue
            case "-dd", "--delete-database":
                if !flag_delete {
                    deleteDatabase()
                    flag_delete = true
                }
                continue
            case "-da", "--delete-archive":
                if !flag_delete {
                    deleteArchive()
                    flag_delete = true
                }
                continue

            // Arguments with parameters
            case "-a", "--address": 
                flag_address = true
                continue
            case "-l", "--login":
                flag_login = true
                continue
            case "-p", "--password":
                flag_password = true
                continue
        }

        switch {
            case flag_address: 
                var ipv4_port = strings.Split(value, ":")
                if len(ipv4_port) != 2 {
                    utils.PrintError("invalid argument for '--address'")
                }
                User.IPv4 = ipv4_port[0]
                User.Port = ":" + ipv4_port[1]
                flag_address = false

            case flag_login:
                flag_login = false
                auth.login = value

            case flag_password:
                flag_password = false
                auth.password = value
        }
    }

    os.Mkdir(PATH_ARCHIVE, 0777)
    switch Authorization(auth.login, auth.password) {
        case 1: // pass
        case 2: utils.PrintError("length of login > 64 bytes")
        case 3: utils.PrintError("passwords hash undefined")
        case 4: utils.PrintError("login or password is wrong")
        default: fmt.Println("[Success]: Authorization")
    }

    switch User.Port {
        case ":": utils.PrintError("port undefined in argument '--address'")
        default: // pass
    }

    return flag_interface
}

func deleteAllData() {
    deleteDatabase()
    deleteArchive()
}

func deleteDatabase() {
    if !utils.FileIsExist(DATABASE_NAME) {
        return
    }
    overwriteAndDelete(DATABASE_NAME, len(utils.ReadFile(DATABASE_NAME)))
}

func deleteArchive() {
    if !utils.FileIsExist(PATH_ARCHIVE) {
        return
    }
    files, err := ioutil.ReadDir(PATH_ARCHIVE)
    utils.CheckError(err)
    for _, file := range files {
        name := file.Name()
        overwriteAndDelete(PATH_ARCHIVE + name, len(utils.ReadFile(PATH_ARCHIVE + name)))
    }
}

func overwriteAndDelete(file string, length int) {
    overwrite(file, length, 16)
    os.Remove(file)
}

func overwrite(file string, length, iter int) {
    for i := 0; i < iter; i++ {
        utils.WriteFile(file, string(crypto.SessionKey(length)))
    }
}
