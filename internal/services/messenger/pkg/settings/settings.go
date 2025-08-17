package settings

import (
	"github.com/number571/hidden-lake/internal/services"
	"github.com/number571/hidden-lake/internal/utils/appname"
)

var (
	gShortAppName = appname.ToShortAppName(CAppFullName)
)

func GetShortAppName() string {
	return gShortAppName
}

const (
	CAppName = "messenger"
)

const (
	CAppFullName    = services.CServiceDomainName + "=" + CAppName
	CAppDescription = "messenger with a web interface"
)

const (
	CPathYML = services.CServiceShortName + "-" + CAppName + ".yml"
	CPathDB  = services.CServiceShortName + "-" + CAppName + ".db"
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
