package settings

import (
	"github.com/number571/hidden-lake/internal/services"
	"github.com/number571/hidden-lake/internal/utils/appname"
)

var (
	gShortAppName = appname.ToShortAppName(CAppFullName)
)

func GetShortAppName() string {
	return gShortAppName
}

const (
	CAppName = "pinger"
)

const (
	CAppFullName    = services.CServiceDomainName + "=" + CAppName
	CAppDescription = "ping the node to check the online status"
)

const (
	CPathYML = services.CServiceShortName + "-" + CAppName + ".yml"
)

const (
	CPingPath = "/ping"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9552"
)
