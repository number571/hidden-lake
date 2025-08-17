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
	CAppName = "http"
)

const (
	CAppFullName    = adapters.CAdapterFullPrefix + "=" + CAppName
	CAppDescription = "adapts HL traffic to a custom HTTP connection"
)

const (
	CPathYML = adapters.CAdapterShortPrefix + "-" + CAppName + ".yml"
	CPathDB  = adapters.CAdapterShortPrefix + "-" + CAppName + ".db"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9511"
	CDefaultInternalAddress = "127.0.0.1:9512"
)
