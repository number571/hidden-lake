package settings

import "github.com/number571/hidden-lake/internal/utils/name"

var (
	gServiceName = name.LoadServiceName(CServiceFullName)
)

func GetServiceName() name.IServiceName {
	return gServiceName
}

const (
	CServiceFullName    = "hidden-lake-service=filesharer"
	CServiceDescription = "file sharing with a web interface"
)

const (
	CPathYML = "hls_filesharer.yml"
	CPathSTG = "hls_filesharer.stg"
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
