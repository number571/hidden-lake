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
	CAppServiceName = "notifier"
	CAppDescription = "send and recv text messages"
)

const (
	CPathYML = CAppShortName + ".yml"
	CPathDB  = CAppShortName + ".db"
)

const (
	CNotifyPath = "/notify"
	CMsgboxPath = "/msgbox"
)

const (
	CDefaultInternalAddress = "127.0.0.1:9591"
	CDefaultExternalAddress = "127.0.0.1:9592"
)

const (
	CHandleIndexPath          = "/api/index"
	CHandleChanListenPath     = "/api/chan/listen"
	CHandleChanNotifyPath     = "/api/chan/notify"
	CHandleChanMsgboxPushPath = "/api/chan/msgbox/push"
	CHandleChanMsgboxLoadPath = "/api/chan/msgbox/load"
	CHandleChanMsgboxSizePath = "/api/chan/msgbox/size"
)

const (
	CHeaderFriendName = "Hls-Notifier-Friend-Name"
	CHeaderMessageUID = "Hls-Notifier-Message-Uid"
)
