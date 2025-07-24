package settings

import (
	"github.com/number571/hidden-lake/internal/services"
	"github.com/number571/hidden-lake/internal/utils/appname"
)

var (
	gFmtAppName = appname.LoadAppName(CAppFullName)
)

func GetFmtAppName() appname.IFmtAppName {
	return gFmtAppName
}

const (
	CAppName = "filesharer"
)

const (
	CAppFullName    = services.CServiceFullPrefix + "=" + CAppName
	CAppDescription = "file sharing with a web interface"
)

const (
	CPathYML = services.CServiceShortPrefix + "-" + CAppName + ".yml"
	CPathSTG = services.CServiceShortPrefix + "-" + CAppName + ".stg"
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
