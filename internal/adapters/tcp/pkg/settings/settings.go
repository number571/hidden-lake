package settings

import (
	"github.com/number571/hidden-lake/internal/adapters"
	"github.com/number571/hidden-lake/internal/utils/appname"
)

var (
	gShortAppName = appname.ToShortAppName(CAppFullName)
)

func GetShortAppName() string {
	return gShortAppName
}

const (
	CAppName = "tcp"
)

const (
	CAppFullName    = adapters.CAdapterDomainName + "=" + CAppName
	CAppDescription = "adapts HL traffic to a custom TCP connection"
)

const (
	CPathYML = adapters.CAdapterShortName + "-" + CAppName + ".yml"
	CPathDB  = adapters.CAdapterShortName + "-" + CAppName + ".db"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9521"
	CDefaultInternalAddress = "127.0.0.1:9522"
)
