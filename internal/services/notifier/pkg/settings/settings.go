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
	CAppDescription = "send and recv notification bytes"
)

const (
	CPathYML = CAppShortName + ".yml"
)

const (
	CNotifyPath = "/notify"
)

const (
	CDefaultInternalAddress = "127.0.0.1:9591"
	CDefaultExternalAddress = "127.0.0.1:9592"
)

const (
	CHandleIndexPath      = "/api/index"
	CHandleChanListenPath = "/api/chan/listen"
	CHandleChanNotifyPath = "/api/chan/notify"
)

const (
	CHeaderFriendName = "Hls-Notifier-Friend-Name"
)
