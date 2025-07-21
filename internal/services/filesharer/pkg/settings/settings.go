package settings

import (
	"github.com/number571/hidden-lake/internal/services"
	"github.com/number571/hidden-lake/internal/utils/name"
)

var (
	gServiceName = name.LoadServiceName(CServiceFullName)
)

func GetServiceName() name.IServiceName {
	return gServiceName
}

const (
	CServiceName = "filesharer"
)

const (
	CServiceFullName    = services.CServiceFullPrefix + "=" + CServiceName
	CServiceDescription = "file sharing with a web interface"
)

const (
	CPathYML = services.CServiceShortPrefix + "-" + CServiceName + ".yml"
	CPathSTG = services.CServiceShortPrefix + "-" + CServiceName + ".stg"
)

const (
	CDefaultInternalAddress = "127.0.0.1:9541"
	CDefaultExternalAddress = "127.0.0.1:9542"
)

const (
	CDefaultRetryNum   = 2  // count
	CDefaultPageOffset = 10 // count
	CDefaultLanguage   = "" // ENG
)

const (
	CHandleIndexPath          = "/"
	CHandleAboutPath          = "/about"
	CHandleFaviconPath        = "/favicon.ico"
	CHandleSettingsPath       = "/settings"
	CHandleFriendsPath        = "/friends"
	CHandleFriendsStoragePath = "/friends/storage"
	CStaticPath               = "/static/"
)

const (
	CListPath = "/list"
	CLoadPath = "/load"
)
