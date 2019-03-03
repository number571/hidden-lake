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

func Initialization(args []string) {
    var (
        flag_address bool
    )

    if len(args) == 1 || (args[1] == "-h" || args[1] == "--help") {
        utils.PrintHelp()
        os.Exit( EXIT_SUCCESS )
    }

    for _, value := range args[1:] {
        switch value {
            case "-a", "--address": 
                flag_address = true
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
        }
    }

    if User.Port == ":" {
        utils.PrintError("port undeclared")
    }

    os.Mkdir(PATH_KEYS, 0777)
    os.Mkdir(PATH_ARCHIVE, 0777)

    if utils.FileIsExist(PATH_KEYS + "private_key") {
        User.PrivateData = utils.ReadFile(PATH_KEYS + "private_key")
        User.PrivateKey = encoding.DecodePrivate(User.PrivateData)
        User.PublicKey = &(User.PrivateKey).PublicKey

    } else {
        User.PrivateKey, User.PublicKey = crypto.GenerateKeys(2048)
        User.PrivateData = string(encoding.EncodePrivate(User.PrivateKey))
        utils.WriteFile(PATH_KEYS + "private_key", User.PrivateData)
    }

    var pub_data = encoding.EncodePublic(User.PublicKey)
    var hashed = md5.Sum(pub_data)

    User.PublicData = string(pub_data)
    User.Name = hex.EncodeToString(hashed[:])

    utils.WriteFile(PATH_KEYS + "public_key", User.PublicData)
}
