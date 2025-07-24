package settings

import (
	"github.com/number571/hidden-lake/internal/adapters"
	"github.com/number571/hidden-lake/pkg/utils/appname"
)

var (
	gFmtAppName = appname.LoadAppName(CAppFullName)
)

func GetFmtAppName() appname.IFmtAppName {
	return gFmtAppName
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
