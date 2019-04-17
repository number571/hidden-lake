package settings

import (
    "os"
    "fmt"
    "strings"
    "io/ioutil"
    "crypto/md5"
    "encoding/hex"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "../utils"
    "../crypto"
    "../encoding"
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

    os.Mkdir(PATH_CONFIG, 0777)
    os.Mkdir(PATH_KEYS, 0777)
    os.Mkdir(PATH_PASW, 0777)
    os.Mkdir(PATH_ARCHIVE, 0777)

    switch Authorization(auth.login, auth.password) {
        case 1: // pass
        case 2: utils.PrintError("length of login > 64 bytes")
        case 3: utils.PrintError("password.hash undefined")
        case 4: utils.PrintError("login or password is wrong")
        default: fmt.Println("[Success]: Authorization")
    }

    switch User.Port {
        case ":": utils.PrintError("port undefined in argument '--address'")
        default: // pass
    }

    return flag_interface
}

func Authorization(login, password string) int8 {
    switch {
        case login == "": return 1
        case len(login) > 64: return 2
        default: // continue
    }

    var concat = login + password

    // Signup
    if !utils.FileIsExist(FILE_PRIVATE_KEY) {
        var new_pasw = createPassword(concat)
        createAsymmetricKeys(new_pasw)
        utils.WriteFile(DATABASE_NAME, "")
    }

    // Login
    if result := checkPassword(concat); result != 0 {
        return result
    }
    User.Login = login
    User.Auth = true

    checkPrivateKey()
    checkConnects()
    checkSettings()
    runDataBase()

    return 0
}

func runDataBase() {
    if !utils.FileIsExist(DATABASE_NAME) {
        utils.WriteFile(DATABASE_NAME, "")
    }

    var err error
    DataBase, err = sql.Open("sqlite3", DATABASE_NAME)
    utils.CheckError(err)

    _, err = DataBase.Exec(`
CREATE TABLE IF NOT EXISTS Email (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    Title VARCHAR(128),
    Body TEXT,
    User VARCHAR(32),
    Date VARCHAR(64) NULL
);

CREATE TABLE IF NOT EXISTS GlobalMessages (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    User VARCHAR(32),
    Body TEXT
);

CREATE TABLE IF NOT EXISTS Connections (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    User VARCHAR(32) UNIQUE,
    Login VARCHAR(64),
    PublicKey VARCHAR(1024) UNIQUE
);
    `)
    utils.CheckError(err)
}

func createPassword(pasw string) []byte {
    var new_pasw = crypto.HashSum([]byte(pasw))
    utils.WriteFile(
        FILE_PASSWORD, 
        hex.EncodeToString(crypto.HashSum(new_pasw)),
    )
    return new_pasw
}

func createAsymmetricKeys(pasw []byte) {
    Mutex.Lock()
    User.PrivateKey, User.PublicKey = crypto.GenerateKeys(2048)
    User.PrivateData = string(encoding.EncodePrivate(User.PrivateKey))
    Mutex.Unlock()
    utils.WriteFile(FILE_PRIVATE_KEY, crypto.Encrypt(
        pasw,
        User.PrivateData,
    ))
}

func checkPassword(pasw string) int8 {
    if !utils.FileIsExist(FILE_PASSWORD) {
        return 3
    }
    var hash = utils.ReadFile(FILE_PASSWORD)
    var new_pasw = crypto.HashSum([]byte(pasw))
    var hash_input = hex.EncodeToString(crypto.HashSum(new_pasw))
    if hash != hash_input {
        return 4
    }
    Mutex.Lock()
    User.Password = new_pasw
    Mutex.Unlock()
    return 0
}

func checkConnects() {
    if !utils.FileIsExist(FILE_CONNECTS) {
        utils.WriteFile(FILE_CONNECTS, crypto.Encrypt(User.Password, ""))
    } else {
        var data = crypto.Decrypt(User.Password, utils.ReadFile(FILE_CONNECTS))
        User.DefaultConnections = strings.Split(data, "\r\n")
    }
}

func checkSettings() {
    if !utils.FileIsExist(FILE_SETTINGS) {
        return
    }

    var slice = strings.Split(crypto.Decrypt(User.Password, utils.ReadFile(FILE_SETTINGS)), ":")
    if len(slice) == 2 {
        Mutex.Lock()
        User.IPv4 = slice[0]
        User.Port = ":" + slice[1]
        Mutex.Unlock()
    }
}

func checkPrivateKey() {
    if utils.FileIsExist(FILE_PRIVATE_KEY) {
        Mutex.Lock()
        User.PrivateData = crypto.Decrypt(
            User.Password,
            utils.ReadFile(FILE_PRIVATE_KEY),
        )
        User.PrivateKey = encoding.DecodePrivate(User.PrivateData)
        User.PublicKey = &(User.PrivateKey).PublicKey
        Mutex.Unlock()
    } else {
        createAsymmetricKeys(User.Password)
    }

    var pub_data = encoding.EncodePublic(User.PublicKey)
    var hashed = md5.Sum(pub_data)

    Mutex.Lock()
    User.PublicData = string(pub_data)
    User.Hash = hex.EncodeToString(hashed[:])
    Mutex.Unlock()

    utils.WriteFile(FILE_PUBLIC_KEY, crypto.Encrypt(
        User.Password,
        User.PublicData,
    ))
}

func deleteAllData() {
    deleteArchive()
    deleteConfig()
    deleteDatabase()
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

func deleteConfig() {
    if !utils.FileIsExist(PATH_CONFIG) {
        return
    }

    var list_of_deletes = []string{
        FILE_PRIVATE_KEY,
        FILE_PUBLIC_KEY,
        FILE_PASSWORD,
        FILE_CONNECTS,
        FILE_SETTINGS,
    }

    for _, value := range list_of_deletes {
        if utils.FileIsExist(value) {
            overwriteAndDelete(value, len(utils.ReadFile(value)))
        }
    }
}

func deleteDatabase() {
    if !utils.FileIsExist(DATABASE_NAME) {
        return
    }
    overwriteAndDelete(DATABASE_NAME, len(utils.ReadFile(DATABASE_NAME)))
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
