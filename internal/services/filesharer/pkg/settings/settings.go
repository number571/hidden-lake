package settings

import (
	"github.com/number571/hidden-lake/internal/services"
)

const (
	CServiceName    = "filesharer"
	CAppShortName   = services.CServiceShortName + "-" + CServiceName
	CAppFullName    = services.CServiceDomainName + "=" + CServiceName
	CAppDescription = "file sharing with a web interface"
)

const (
	CPathYML = CAppShortName + ".yml"
	CPathSTG = CAppShortName + ".stg"
	CPathTMP = CAppShortName + "-%s.tmp"
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
