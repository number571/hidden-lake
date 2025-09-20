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
	CServiceName    = "pinger"
	CAppShortName   = services.CServiceShortName + "-" + CServiceName
	CAppFullName    = services.CServiceDomainName + "=" + CServiceName
	CAppDescription = "ping the node to check the online status"
)

const (
	CPathYML = CAppShortName + ".yml"
)

const (
	CPingPath = "/ping"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9552"
)
