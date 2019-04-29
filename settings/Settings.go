package settings

import (
    "net"
    "sync"
    "net/http"
    "crypto/rsa"
    "database/sql"
    "../models"
)

type PackageTCP struct {
    From models.From
    To string
    models.Head
    Body string
}

type PackageHTTP struct {
    Exists  bool
    Head    string
    Body    string
}

type UserNode struct {
    models.Keys
    models.Messages
    models.Connection
    models.Authorization
    models.Transportation
}

var (
    GoroutinesIsRun = false
    ServerListenTCP net.Listener
    ServerListenHTTP *http.Server
    Mutex sync.Mutex
    DataBase *sql.DB
)

var User = UserNode {
    Keys: models.Keys {
        NodePublicKey:  make(map[string]*rsa.PublicKey),
        NodeSessionKey: make(map[string][]byte),
        NodeConnection: make(map[string]int8),
    },
    Connection: models.Connection {
        NodeAddress: make(map[string]string),
        NodeLogin: make(map[string]string),
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
    MODE_LOCAL  = "[LOCAL]"
    MODE_GLOBAL = "[GLOBAL]"

    OPT_LIST = "[LIST]"
    OPT_FILE = "[FILE]"

    MODE_READ_LIST = MODE_READ + OPT_LIST
    MODE_SAVE_LIST = MODE_SAVE + OPT_LIST

    MODE_READ_FILE = MODE_READ + OPT_FILE
    MODE_SAVE_FILE = MODE_SAVE + OPT_FILE

    SEPARATOR = "[SEPARATOR]"
    SEPARATOR_ADDRESS  = SEPARATOR + "[ADDRESS]"
)

const (
    PROTOCOL_TCP = "tcp"
    
    PORT_HTTP = ":7545"
    IPV4_HTTP = "127.0.0.1"

    DATABASE_NAME = "database.db"
    IPV4_TEMPLATE = "0.0.0.0"

    DYNAMIC_ROUTING_NODES = true
    QUAN_OF_ROUTING_NODES = 3

    SESSION_KEY_BYTES = 32
    ROUTING_KEY_BYTES = 16

    TIME_SLEEP = 1
    BUFF_SIZE = 512
)

const (
    PATH_CONFIG  = "_Config/"
    PATH_ARCHIVE = "_Archive/"
    PATH_KEYS    = PATH_CONFIG + "Keys/"
    PATH_PASW    = PATH_CONFIG + "Pasw/"

    PATH_VIEWS  = "views/"
    PATH_STATIC = "static/"
)

const (
    FILE_PRIVATE_KEY = PATH_KEYS + "private.key"
    FILE_PUBLIC_KEY  = PATH_KEYS + "public.key"
    FILE_PASSWORD    = PATH_PASW + "password.hash"
    FILE_SETTINGS    = PATH_CONFIG + "settings.cfg"
    FILE_CONNECTS    = PATH_CONFIG + "connects.cfg"
)

const (
    TERM_EXIT       = ":exit"
    TERM_HELP       = ":help"
    TERM_SEND       = ":send"
    TERM_EMAIL      = ":email"
    TERM_WHOAMI     = ":whoami"
    TERM_REFRESH    = ":refresh"
    TERM_ARCHIVE    = ":archive"
    TERM_HISTORY    = ":history"
    TERM_NETWORK    = ":network"
    TERM_CONNECT    = ":connect"
    TERM_LOGIN      = ":login"
    TERM_PASSWORD   = ":password"
    TERM_LOGOUT     = ":logout"
    TERM_ADDRESS    = ":address"
    TERM_ENTER      = ":enter"
    TERM_INTERFACE  = ":interface"
    TERM_DISCONNECT = ":disconnect"
)

const (
    EXIT_SUCCESS = 0
    EXIT_FAILED  = 1
)
