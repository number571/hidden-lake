package settings

import (
	"github.com/number571/hidden-lake/internal/adapters"
	"github.com/number571/hidden-lake/internal/utils/appname"
)

var (
	gAppName = appname.LoadAppName(CServiceFullName)
)

func GetAppName() appname.IAppName {
	return gAppName
}

const (
	CServiceName = "tcp"
)

const (
	CServiceFullName    = adapters.CAdapterFullPrefix + "=" + CServiceName
	CServiceDescription = "adapts HL traffic to a custom TCP connection"
)

const (
	CPathYML = adapters.CAdapterShortPrefix + "-" + CServiceName + ".yml"
	CPathDB  = adapters.CAdapterShortPrefix + "-" + CServiceName + ".db"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9521"
	CDefaultInternalAddress = "127.0.0.1:9522"
)
