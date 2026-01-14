package settings

import (
	"strings"

	"github.com/number571/hidden-lake/internal/services"
)

var (
	gAppShortNameFMT = strings.ToUpper(CAppShortName)
)

func GetAppShortNameFMT() string {
	return gAppShortNameFMT
}

const (
	CAppShortName = services.CServiceShortName + "-" + CAppServiceName
	CAppFullName  = services.CServiceDomainName + "=" + CAppServiceName
)

const (
	CAppServiceName = "messenger"
	CAppDescription = "chat with a web interface"
)

const (
	CPathYML = CAppShortName + ".yml"
	CPathDB  = CAppShortName + ".db"
)

const (
	CPushPath = "/push"
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
	CHandleIndexPath        = "/api/index"
	CHandlePushMessagePath  = "/api/push/message"
	CHandleLoadMessagesPath = "/api/load/messages"
	CHandleListenChatPath   = "/api/listen/chat"
)
