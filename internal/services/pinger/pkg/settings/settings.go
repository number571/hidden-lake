package settings

import (
	"github.com/number571/hidden-lake/internal/services"
	"github.com/number571/hidden-lake/internal/utils/appname"
)

var (
	gAppName = appname.LoadAppName(CServiceFullName)
)

func GetAppName() appname.IAppName {
	return gAppName
}

const (
	CServiceName = "pinger"
)

const (
	CServiceFullName    = services.CServiceFullPrefix + "=" + CServiceName
	CServiceDescription = "ping the node to check the online status"
)

const (
	CPathYML = services.CServiceShortPrefix + "-" + CServiceName + ".yml"
)

const (
	CPingPath = "/ping"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9552"
)
