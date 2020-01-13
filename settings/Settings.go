package settings

import (
	"../models"
	"database/sql"
	"github.com/number571/gopeer"
	"time"
)

var (
	Listener *gopeer.Listener
	Tokens   = make(map[string]string)
	Users    = make(map[string]*models.User)
	CFG      *models.Config
	DB       *sql.DB
)

const (
	PATH_VIEWS  = "views/"
	PATH_STATIC = "static/"
	DB_NAME     = "database.db"
	CFG_NAME    = "config.cfg"
)

const (
	TITLE_MESSAGE = "[TITLE-MESSAGE]"
)

var (
	OPTION_GET = gopeer.Get("OPTION_GET").(string)
)

const (
	// len(base64(sha256(x))) = 44
	LEN_BASE64_SHA256 = 44
	CHECK_DURING      = 12 * time.Hour
	LIFETIME_TOKEN    = 24 * time.Hour
)
