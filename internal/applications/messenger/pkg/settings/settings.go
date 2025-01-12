package settings

import "github.com/number571/hidden-lake/internal/utils/name"

var (
	GServiceName = name.LoadServiceName(CServiceFullName)
)

const (
	CServiceFullName    = "hidden-lake-messenger"
	CServiceDescription = "messenger with a web interface"
)

const (
	CPathYML = "hlm.yml"
	CPathDB  = "hlm.db"
)

const (
	CStaticPath = "/static/"
	CPushPath   = "/push"
)

const (
	CDefaultInternalAddress = "127.0.0.1:9591"
	CDefaultExternalAddress = "127.0.0.1:9592"
)

const (
	CDefaultMessagesCapacity = (2 << 10) // count
	CDefaultLanguage         = ""        // ENG
)

const (
	CHandleIndexPath         = "/"
	CHandleAboutPath         = "/about"
	CHandleFaviconPath       = "/favicon.ico"
	CHandleSettingsPath      = "/settings"
	CHandleFriendsPath       = "/friends"
	CHandleFriendsChatPath   = "/friends/chat"
	CHandleFriendsUploadPath = "/friends/upload"
	CHandleFriendsChatWSPath = "/friends/chat/ws"
)

const (
	CIsText = 0x01
	CIsFile = 0x02
)
