package settings

import (
	"database/sql"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
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
	PATH_TLS     = "tls/"
	PATH_VIEWS   = "views/"
	PATH_STATIC  = "static/"
	PATH_INPUT   = "inputd/"
	PATH_ARCHIVE = PATH_STATIC + "archive/"
	DB_NAME      = PATH_INPUT + "database.db"
	CFG_NAME     = PATH_INPUT + "config.json"
	UPD_NAME     = PATH_INPUT + "updates.json"
)

// Tools | Archive
// Tools: check and sign messages
// Archive: files
const (
	TITLE_MESSAGE  = "[TITLE-MESSAGE]"
	TITLE_ARCHIVE  = "[TITLE-ARCHIVE]"
)

const (
	FILE_PART_SIZE    = 8 << 20 // 8MiB
	BUFFER_SIZE       = 2 << 20 // 1MiB
	CHECK_DURING      = 12 * time.Hour
	LIFETIME_TOKEN    = 24 * time.Hour
)
