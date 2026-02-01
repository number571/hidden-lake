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
	CAppDescription = "view storage and download files"
)

const (
	CPathYML = CAppShortName + ".yml"
	CPathSTG = CAppShortName + ".stg"
	CPathTMP = CAppShortName + "-%s.tmp"
)

const (
	CPathSharingSTG       = CPathSTG + "/sharing"
	CPathPrivateSTG       = CPathSTG + "/private"
	CPathSharingPublicSTG = CPathSharingSTG + "/public"
)

const (
	CDefaultInternalAddress = "127.0.0.1:9541"
	CDefaultExternalAddress = "127.0.0.1:9542"
)

const (
	CDefaultRetryNum   = 2  // count
	CDefaultPageOffset = 10 // count
)

const (
	CInfoPath = "/info"
	CListPath = "/list"
	CLoadPath = "/load"
)

const (
	CHandleIndexPath          = "/api/index"
	CHandleRemoteListPath     = "/api/remote/list"
	CHandleRemoteFilePath     = "/api/remote/file"
	CHandleRemoteFileInfoPath = "/api/remote/file/info"
	CHandleLocalListPath      = "/api/local/list"
	CHandleLocalFilePath      = "/api/local/file"
	CHandleLocalFileInfoPath  = "/api/local/file/info"
)

const (
	CHeaderInProcess = "Hls-Filesharer-In-Process"
	CHeaderFileHash  = "Hls-Filesharer-File-Hash"
)

const (
	CHeaderProcessModeY = "+" // default
	CHeaderProcessModeN = "-"
)
