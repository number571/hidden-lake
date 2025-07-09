package settings

import "github.com/number571/hidden-lake/internal/utils/name"

var (
	gServiceName = name.LoadServiceName(CServiceFullName)
)

func GetServiceName() name.IServiceName {
	return gServiceName
}

const (
	CServiceFullName    = "hidden-lake-service=messenger"
	CServiceDescription = "messenger with a web interface"
)

const (
	CPathYML = "hls_messenger.yml"
	CPathDB  = "hls_messenger.db"
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
