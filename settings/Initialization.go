package settings

import (
    "os"
    "strings"
    "crypto/md5"
    "encoding/hex"
    "../utils"
    "../crypto"
    "../encoding"
)

func Initialization(args []string) bool {
    var (
        flag_address bool
        flag_interface bool
        flag_white_list bool
        flag_black_list bool
        flag_private_key bool
        change_config [3]string
    )

    if len(args) == 1 || (args[1] == "-h" || args[1] == "--help") {
        utils.PrintHelp()
        os.Exit( EXIT_SUCCESS )
    }

    for _, value := range args[1:] {
        switch value {
            // Arguments without parameters
            case "-i", "--interface": 
                flag_interface = true
                continue

            // Arguments with parameters
            case "-a", "--address": 
                flag_address = true
                continue
            case "-w", "--white-list":
                flag_white_list = true
                continue
            case "-b", "--black-list":
                flag_black_list = true
                continue
            case "-p", "--private-key":
                flag_private_key = true
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

            case flag_white_list:
                flag_white_list = false
                change_config[0] = value

            case flag_black_list:
                flag_black_list = false
                change_config[1] = value

            case flag_private_key:
                flag_private_key = false
                change_config[2] = value
        }
    }

    if User.Port == ":" {
        utils.PrintError("port undeclared")
    }

    os.Mkdir(PATH_CONFIG, 0777)
    os.Mkdir(PATH_KEYS, 0777)
    os.Mkdir(PATH_ARCHIVE, 0777)

    // Profile info
    if !utils.FileIsExist(PATH_CONFIG + "profile_info.cfg") {
        var default_info = "default info"
        utils.WriteFile(PATH_CONFIG + "profile_info.cgf", default_info)
        User.Info = default_info
    } else {
        User.Info = utils.ReadFile(PATH_CONFIG + "profile_info.cgf")
    }

    initDataBase()
    checkConfig(change_config)
    return flag_interface
}

func CreateDatabase(name string) {
    if !utils.FileIsExist(name) {
        utils.WriteFile(name, "")
    }
}

func initDataBase() {
    _, err := DataBase.Exec(`
CREATE TABLE IF NOT EXISTS Email (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    Title VARCHAR(128),
    Body TEXT,
    User VARCHAR(128),
    Date VARCHAR(64) NULL
);

DROP TABLE IF EXISTS GlobalMessages;
CREATE TABLE GlobalMessages (
    Id INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    User VARCHAR(128),
    Body TEXT
);
`)
    utils.CheckError(err)
}

func checkConfig(config [3]string) {
    // White List
    if config[0] != "" {
        var white_list = utils.ReadFile(config[0])
        var addresses = strings.Split(white_list, "\n")
        for _, addr := range addresses {
            User.WhiteList = append(User.WhiteList, addr)
        }
    } else {
        if !utils.FileIsExist(PATH_CONFIG + "white_list.cfg") {
            utils.WriteFile(PATH_CONFIG + "white_list.cfg", "")
        } else {
            var white_list = utils.ReadFile(PATH_CONFIG + "white_list.cfg")
            var addresses = strings.Split(white_list, "\n")
            for _, addr := range addresses {
                User.WhiteList = append(User.WhiteList, addr)
            }
        }
    }

    // Black List
    if config[1] != "" {
        var black_list = utils.ReadFile(config[1])
        var addresses = strings.Split(black_list, "\n")
        for _, addr := range addresses {
            User.BlackList = append(User.BlackList, addr)
        }
    } else {
        if !utils.FileIsExist(PATH_CONFIG + "black_list.cfg") {
            utils.WriteFile(PATH_CONFIG + "black_list.cfg", "")
        } else {
            var black_list = utils.ReadFile(PATH_CONFIG + "black_list.cfg")
            var addresses = strings.Split(black_list, "\n")
            for _, addr := range addresses {
                User.BlackList = append(User.BlackList, addr)
            }
        }
    }

    // Private Key
    if config[2] != "" {
        User.PrivateData = utils.ReadFile(config[2])
        User.PrivateKey = encoding.DecodePrivate(User.PrivateData)
        User.PublicKey = &(User.PrivateKey).PublicKey
    } else {
        if utils.FileIsExist(PATH_KEYS + "private_key") {
            User.PrivateData = utils.ReadFile(PATH_KEYS + "private_key")
            User.PrivateKey = encoding.DecodePrivate(User.PrivateData)
            User.PublicKey = &(User.PrivateKey).PublicKey
        } else {
            User.PrivateKey, User.PublicKey = crypto.GenerateKeys(2048)
            User.PrivateData = string(encoding.EncodePrivate(User.PrivateKey))
            utils.WriteFile(PATH_KEYS + "private_key", User.PrivateData)
        }
    }

    var pub_data = encoding.EncodePublic(User.PublicKey)
    var hashed = md5.Sum(pub_data)

    User.PublicData = string(pub_data)
    User.Name = hex.EncodeToString(hashed[:])

    utils.WriteFile(PATH_KEYS + "public_key", User.PublicData)
}
