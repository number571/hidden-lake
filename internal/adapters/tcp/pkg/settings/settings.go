package settings

import (
	"github.com/number571/hidden-lake/internal/adapters"
	"github.com/number571/hidden-lake/internal/utils/appname"
)

var (
	gFmtAppName = appname.LoadAppName(CAppFullName)
)

func GetFmtAppName() appname.IFmtAppName {
	return gFmtAppName
}

const (
	CAppName = "tcp"
)

const (
	CAppFullName    = adapters.CAdapterFullPrefix + "=" + CAppName
	CAppDescription = "adapts HL traffic to a custom TCP connection"
)

const (
	CPathYML = adapters.CAdapterShortPrefix + "-" + CAppName + ".yml"
	CPathDB  = adapters.CAdapterShortPrefix + "-" + CAppName + ".db"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9521"
	CDefaultInternalAddress = "127.0.0.1:9522"
)
