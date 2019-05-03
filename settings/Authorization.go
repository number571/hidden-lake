package settings

import (
    "strings"
    "crypto/md5"
    "encoding/hex"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "../utils"
    "../crypto"
    "../encoding"
)

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

    initPrivateKey()
    initSettings()
    initConnects()

    initDataBase()
    initConnectsF2F()

    return 0
}

func initConnectsF2F() {
    rows, err := DataBase.Query("SELECT User, Address, SessionKey FROM ConnectionsF2F")
    utils.CheckError(err)
    defer rows.Close()

    var user, address, encoded_session_key string

    for rows.Next() {
        rows.Scan(&user, &address, &encoded_session_key)

        address = crypto.Decrypt(User.Password, address)
        session_key, err := hex.DecodeString(crypto.Decrypt(User.Password, encoded_session_key))
        utils.CheckError(err)

        Messages.NewDataExistLocal[user] = make(chan bool)
        User.NodeAddressF2F[user] = address
        User.NodeSessionKeyF2F[user] = session_key
    }
}

func initConnects() {
    if !utils.FileIsExist(FILE_CONNECTS) {
        utils.WriteFile(FILE_CONNECTS, crypto.Encrypt(User.Password, ""))
    } else {
        var data = crypto.Decrypt(User.Password, utils.ReadFile(FILE_CONNECTS))
        User.DefaultConnections = strings.Split(data, "\r\n")
    }
}

func initDataBase() {
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
    Mode VARCHAR(4),
    Body TEXT
);

CREATE TABLE IF NOT EXISTS Connections (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    User VARCHAR(32) UNIQUE,
    Login VARCHAR(64),
    PublicKey VARCHAR(1024) UNIQUE
);

CREATE TABLE IF NOT EXISTS ConnectionsF2F (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    User VARCHAR(32) UNIQUE,
    Address VARCHAR(32),
    SessionKey VARCHAR(64)
);

CREATE TABLE IF NOT EXISTS HiddenFriends (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    User VARCHAR(32) UNIQUE
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

func initSettings() {
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

func initPrivateKey() {
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
