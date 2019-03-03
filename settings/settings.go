package settings

import (
	"sync"
	"crypto/rsa"
	"../models"
)

type Package struct {
    From models.From
    To   string
    models.Head
    Body string
}

type UserNode struct {
	models.Keys
	models.Connection
	models.ChatMessages
	models.Transportation
}

var Mutex sync.Mutex
var User = UserNode {
	Transportation: models.Transportation {
		Info: "default_info",
	},
	Keys: models.Keys {
		NodePublicKey: make(map[string]*rsa.PublicKey),
		NodeSessionKey: make(map[string][]byte),
		NodeConnection: make(map[string]int8),
	},
	ChatMessages: models.ChatMessages {
		LocalMessages: make(map[string][]string),
	},
}

const (
	HEAD_ARCHIVE = "[ARCHIVE]"
	HEAD_PROFILE = "[PROFILE]"
	HEAD_MESSAGE = "[MESSAGE]"
	HEAD_CONNECT = "[CONNECT]"
	HEAD_WARNING = "[WARNING]"

	MODE_GET = "[GET]"
	MODE_SAVE = "[SAVE]"

	MODE_LOCAL = "[LOCAL]"
	MODE_GLOBAL = "[GLOBAL]"

	MODE_GET_LIST = MODE_GET + "[LIST]"
	MODE_SAVE_LIST = MODE_SAVE + "[LIST]"

	MODE_GET_FILE = MODE_GET + "[FILE]"
	MODE_SAVE_FILE = MODE_SAVE + "[FILE]"

	SEPARATOR = "[SEPARATOR]"
)

const (
	PROTOCOL_TCP = "tcp"
	
	PORT_HTTP = ":7546"
	IPV4_HTTP = "127.0.0.1"

	IPV4_TEMPLATE = "0.0.0.0"

	TIME_SLEEP = 1
	BUFF_SIZE = 512
)

const (
	PATH_KEYS = "Keys/"
	PATH_ARCHIVE = "Archive/"
	PATH_STATIC = "static/"
	PATH_TEMPLATES = "templates/"
)

const (
	TERM_EXIT = ":exit"
	TERM_HELP = ":help"
	TERM_SEND = ":send"
	TERM_ARCHIVE = ":archive"
	TERM_HISTORY = ":history"
	TERM_NETWORK = ":network"
	TERM_CONNECT = ":connect"
	TERM_DISCONNECT = ":disconnect"
)

const (
	EXIT_SUCCESS = 0
	EXIT_FAILED  = 1
)
