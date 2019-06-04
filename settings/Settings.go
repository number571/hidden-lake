package settings

import (
    "net"
    "sync"
    "net/http"
    "crypto/rsa"
    "database/sql"
    "../models"
)

var (
    // ViewCLI = true
    NeedF2FMode bool = false
    GoroutinesIsRun bool = false
    ConnectionIsRead bool = false
    ServerListenTCP net.Listener
    ServerListenHTTP *http.Server
    Mutex sync.Mutex
    DataBase *sql.DB
)

var User = models.UserNode {
    Mode: models.C_S_mode,
}

var Node = models.Node {
    PublicKey:  make(map[string]*rsa.PublicKey),
    ConnectionMode: make(map[string]models.ModeConn),
    SessionKey: models.SessionKey {
        P2P: make(map[string][]byte),
        F2F: make(map[string][]byte),
    },
    Address: models.Address {
        P2P: make(map[string]string),
        F2F: make(map[string]string),
        C_S: make(map[string]net.Conn),
    },
}

var Messages = models.Messages {
    NewDataExistGlobal: make(chan bool),
    NewDataExistLocal: make(map[string]chan bool),
    CurrentIdLocal: make(map[string]uint16),
}

const (
    HEAD_EMAIL    = "[EMAIL]"
    HEAD_ARCHIVE  = "[ARCHIVE]"
    HEAD_MESSAGE  = "[MESSAGE]"
    HEAD_CONNECT  = "[CONNECT]"
    HEAD_WARNING  = "[WARNING]"
    HEAD_REDIRECT = "[REDIRECT]"

    MODE_READ   = "[READ]"
    MODE_SAVE   = "[SAVE]"

    MODE_LIST   = "[LIST]"
    MODE_FILE   = "[FILE]"
    MODE_CHECK  = "[CHECK]"
    MODE_LOCAL  = "[LOCAL]"
    MODE_GLOBAL = "[GLOBAL]"

    MODE_READ_CHECK = MODE_READ + MODE_CHECK
    MODE_SAVE_CHECK = MODE_SAVE + MODE_CHECK

    MODE_READ_GLOBAL = MODE_READ + MODE_GLOBAL
    MODE_SAVE_GLOBAL = MODE_SAVE + MODE_GLOBAL

    MODE_READ_LOCAL = MODE_READ + MODE_LOCAL
    MODE_SAVE_LOCAL = MODE_SAVE + MODE_LOCAL

    MODE_READ_LIST = MODE_READ + MODE_LIST
    MODE_SAVE_LIST = MODE_SAVE + MODE_LIST

    MODE_READ_FILE = MODE_READ + MODE_FILE
    MODE_SAVE_FILE = MODE_SAVE + MODE_FILE

    SEPARATOR = "[SEPARATOR]"
)

const (
    PROTOCOL  = "tcp"
    PORT_HTTP = ":7545"
    IPV4_HTTP = "127.0.0.1"
    END_BLOCK = "[END-BLOCK]"

    DATABASE_NAME = "database.db"
    IPV4_TEMPLATE = "0.0.0.0"

    ASYMMETRIC_KEY_BITS = 2048

    QUAN_OF_ROUTING_NODES = 3
    DYNAMIC_ROUTING = false

    SESSION_KEY_BYTES = 32
    ROUTING_KEY_BYTES = 16

    TIME_SLEEP = 1
    BUFF_SIZE = 512
)

const (
    PATH_ARCHIVE = "_Archive/"
    PATH_VIEWS  = "views/"
    PATH_STATIC = "static/"
)

const (
    TERM_MODE           = ":mode"
    TERM_EXIT           = ":exit"
    TERM_HELP           = ":help"
    TERM_SEND           = ":send"
    TERM_EMAIL          = ":email"
    TERM_WHOAMI         = ":whoami"
    TERM_ARCHIVE        = ":archive"
    TERM_HISTORY        = ":history"
    TERM_NETWORK        = ":network"
    TERM_CONNECT        = ":connect"
    TERM_DISCONNECT     = ":disconnect"
    TERM_LOGIN          = ":login"
    TERM_PASSWORD       = ":password"
    TERM_LOGOUT         = ":logout"
    TERM_ADDRESS        = ":address"
    TERM_ENTER          = ":enter"
    TERM_INTERFACE      = ":interface"
)

const (
    EXIT_SUCCESS = 0
    EXIT_FAILED  = 1
)
