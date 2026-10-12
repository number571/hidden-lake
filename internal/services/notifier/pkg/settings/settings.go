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
	CPushPath = "/push"
)

const (
	CDefaultInternalAddress = "127.0.0.1:9591"
	CDefaultExternalAddress = "127.0.0.1:9592"
)

const (
	CHandleIndexPath      = "/api/index"
	CHandleChanPushPath   = "/api/chan/push"
	CHandleChanLoadPath   = "/api/chan/load"
	CHandleChanSizePath   = "/api/chan/size"
	CHandleChanListenPath = "/api/chan/listen"
)
