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
	PATH_ARCHIVE = PATH_STATIC + "archive/"
	DB_NAME      = "database.db"
	CFG_NAME     = "config.cfg"
)

// Tools | Archive
// Tools: check and sign messages
// Archive: files
const (
	TITLE_MESSAGE  = "[TITLE-MESSAGE]"
	TITLE_ARCHIVE  = "[TITLE-ARCHIVE]"
	TITLE_CONNLIST = "[TITLE-CONNLIST]"
)

var (
	OPTION_GET = gopeer.Get("OPTION_GET").(string)
	IS_CLIENT  = gopeer.Get("IS_CLIENT").(string)
)

const (
	// len(base64(sha256(x))) = 44
	LEN_BASE64_SHA256 = 44
	PACKAGE_SIZE      = 16 << 20 // 16MiB
	BUFFER_SIZE       = 1 << 10  // 1KiB
	CHECK_DURING      = 12 * time.Hour
	LIFETIME_TOKEN    = 24 * time.Hour
)
