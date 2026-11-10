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
	CAppServiceName = "publisher"
	CAppDescription = "push and load publications by options"
)

const (
	CPathYML = CAppShortName + ".yml"
	CPathDB  = CAppShortName + ".db"
)

const (
	CPushPath = "/push"
	CLoadPath = "/load"
	CSizePath = "/size"
)

const (
	CDefaultInternalAddress = "127.0.0.1:9551"
	CDefaultExternalAddress = "127.0.0.1:9552"
)

const (
	CHandleIndexPath        = "/api/index"
	CHandleMsgboxPushPath   = "/api/msgbox/push"
	CHandleMsgboxLoadPath   = "/api/msgbox/load"
	CHandleMsgboxSizePath   = "/api/msgbox/size"
	CHandleMsgboxListenPath = "/api/msgbox/listen"
)

const (
	CHeaderFriendName = "Hls-Publisher-Friend-Name"
	CHeaderMessageNID = "Hls-Publisher-Message-Nid"
)

const (
	CHeaderMessageStateOnLoad     = "M-S"
	CHeaderAnotherMessageIDOnLoad = "M-I"
)

const (
	CUnknownCreatorName = "" // null string
)
