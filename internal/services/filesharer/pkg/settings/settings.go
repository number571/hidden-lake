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
	CAppServiceName = "filesharer"
	CAppDescription = "file sharing with a web interface"
)

const (
	CPathYML = CAppShortName + ".yml"
	CPathSTG = CAppShortName + ".stg"
	CPathTMP = CAppShortName + "-%s.tmp"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9542"
)

const (
	CDefaultRetryNum   = 2  // count
	CDefaultPageOffset = 10 // count
	CDefaultLanguage   = "" // ENG
)

const (
	CInfoPath = "/info"
	CListPath = "/list"
	CLoadPath = "/load"
)

const (
	CHandleIndexPath        = "/api/index"
	CHandleFileInfoPath     = "/api/file/info"
	CHandleFileDownloadPath = "/api/file/download"
	CHandleStorageFilesPath = "/api/storage/files"
)
