package settings

import (
    "crypto/md5"
    "encoding/hex"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    "../utils"
    "../crypto"
    "../models"
    "../encoding"
)

func Authorization(login, password string) int8 {
    switch {
        case login == "": return 1
        case len(login) > 64: return 2
        default: // continue
    }

    var concat = login + password
    initDataBase(DATABASE_NAME)

    // Signup
    if passwordIsNotExist() {
        var new_pasw = createPassword(concat)
        createAsymmetricKeys(new_pasw)
    }

    // Login
    if result := checkPassword(concat); result != 0 {
        return result
    }

    User.Login = login
    User.Auth = true

    initPrivateKey()
    initAddress()
    initConnects()
    initConnectsF2F()
    initMode()

    return 0
}

func initMode() {
    if User.Port == "" { 
        User.Mode = models.C_S_mode
        return 
    }

    if NeedF2FMode { 
        User.Mode = models.F2F_mode 
    } else { 
        User.Mode = models.P2P_mode 
    }
}

func initDataBase(database_name string) {
    if !utils.FileIsExist(database_name) {
        utils.WriteFile(database_name, "")
    }

    var err error
    DataBase, err = sql.Open("sqlite3", database_name)
    utils.CheckError(err)

    _, err = DataBase.Exec(`
CREATE TABLE IF NOT EXISTS Keys (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    Mode VARCHAR(4) UNIQUE,
    PrivateKey VARCHAR(4096) UNIQUE
);

CREATE TABLE IF NOT EXISTS Password (
    Hash VARCHAR(64) UNIQUE
);

CREATE TABLE IF NOT EXISTS Address (
    IPv4 VARCHAR(16) UNIQUE,
    Port VARCHAR(8) UNIQUE
);

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

CREATE TABLE IF NOT EXISTS DefaultConnections (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    Address VARCHAR(32) UNIQUE
);

CREATE TABLE IF NOT EXISTS Connections (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    User VARCHAR(32) UNIQUE,
    PublicKey VARCHAR(1024) UNIQUE
);

CREATE TABLE IF NOT EXISTS ConnectionsF2F (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    User VARCHAR(32) UNIQUE,
    Address VARCHAR(32),
    SessionKey VARCHAR(64)
);
`)
    utils.CheckError(err)
}

func passwordIsNotExist() bool {
    var (
        row = DataBase.QueryRow("SELECT Hash FROM Password")
        password string
    )
    row.Scan(&password)
    if password == "" { return true }
    return false
}

func initPrivateKey() {
    var private_key, mode string

    rows, err := DataBase.Query("SELECT PrivateKey, Mode FROM Keys ORDER BY Id")
    utils.CheckError(err)
    
    for rows.Next() {
        rows.Scan(&private_key, &mode)
        if private_key == "" {
            createAsymmetricKeys(User.Password)
            break
        } else {
            switch mode {
                case "P2P":
                    Mutex.Lock()
                    User.Private.Data.P2P = crypto.Decrypt(User.Password, private_key)
                    User.Private.Key.P2P = encoding.DecodePrivate(User.Private.Data.P2P)
                    User.Public.Key.P2P = &(User.Private.Key.P2P).PublicKey
                    Mutex.Unlock()
                case "F2F":
                    Mutex.Lock()
                    User.Private.Data.F2F = crypto.Decrypt(User.Password, private_key)
                    User.Private.Key.F2F = encoding.DecodePrivate(User.Private.Data.F2F)
                    User.Public.Key.F2F = &(User.Private.Key.F2F).PublicKey
                    Mutex.Unlock()
            }
        }
    }

    rows.Close()

    var pub_data_p2p = encoding.EncodePublic(User.Public.Key.P2P)
    var hashed_p2p = md5.Sum(pub_data_p2p)

    var pub_data_f2f = encoding.EncodePublic(User.Public.Key.F2F)
    var hashed_f2f = md5.Sum(pub_data_f2f)

    Mutex.Lock()
    User.Public.Data.P2P = string(pub_data_p2p)
    User.Public.Data.F2F = string(pub_data_f2f)
    User.Hash.P2P = hex.EncodeToString(hashed_p2p[:])
    User.Hash.F2F = hex.EncodeToString(hashed_f2f[:])
    Mutex.Unlock()
}

func initAddress() {
    var (
        row = DataBase.QueryRow("SELECT IPv4, Port FROM Address")
        ipv4, port string
    )
    row.Scan(&ipv4, &port)
    if port == "" { 
        if User.Port != "" { SaveAddress(User.IPv4, User.Port) }
        return
    }
    
    Mutex.Lock()
    User.IPv4 = crypto.Decrypt(User.Password, ipv4)
    User.Port = crypto.Decrypt(User.Password, port)
    Mutex.Unlock()
}

func initConnects() {
    rows, err := DataBase.Query("SELECT Address FROM DefaultConnections")
    utils.CheckError(err)
    defer rows.Close()

    var addresses []string
    var address string
    for rows.Next() {
        rows.Scan(&address)
        addresses = append(addresses, crypto.Decrypt(User.Password, address))
    }

    Mutex.Lock()
    User.DefaultConnections = addresses
    Mutex.Unlock()
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
        Node.Address.F2F[user] = address
        Node.SessionKey.F2F[user] = session_key
    }
}

func createPassword(pasw string) []byte {
    var new_pasw = crypto.HashSum([]byte(pasw))
    Mutex.Lock()
    _, err := DataBase.Exec(
        "INSERT INTO Password (Hash) VALUES ($1)", 
        hex.EncodeToString(crypto.HashSum(new_pasw)),
    )
    Mutex.Unlock()
    utils.CheckError(err)
    return new_pasw
}

func createAsymmetricKeys(pasw []byte) {
    Mutex.Lock()
    User.Private.Key.P2P, User.Public.Key.P2P = crypto.GenerateKeys(ASYMMETRIC_KEY_BITS)
    User.Private.Data.P2P = string(encoding.EncodePrivate(User.Private.Key.P2P))
    _, err1 := DataBase.Exec(
        "INSERT INTO Keys (PrivateKey, Mode) VALUES ($1, $2)",
        crypto.Encrypt(pasw, User.Private.Data.P2P), "P2P",
    )
    
    User.Private.Key.F2F, User.Public.Key.F2F = crypto.GenerateKeys(ASYMMETRIC_KEY_BITS)
    User.Private.Data.F2F = string(encoding.EncodePrivate(User.Private.Key.F2F))
    _, err2 := DataBase.Exec(
        "INSERT INTO Keys (PrivateKey, Mode) VALUES ($1, $2)", 
        crypto.Encrypt(pasw, User.Private.Data.F2F), "F2F",
    )
    Mutex.Unlock()

    utils.CheckError(err1)
    utils.CheckError(err2)
}

func checkPassword(pasw string) int8 {
    var (
        row = DataBase.QueryRow("SELECT Hash FROM Password")
        hash_password string
    )

    row.Scan(&hash_password)
    if hash_password == "" { return 3 }

    var (
        new_pasw = crypto.HashSum([]byte(pasw))
        hashed_password = hex.EncodeToString(crypto.HashSum(new_pasw))
    )

    if hash_password != hashed_password {
        return 4
    }

    Mutex.Lock()
    User.Password = new_pasw
    Mutex.Unlock()
    return 0
}
