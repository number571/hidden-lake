package settings

import (
	"github.com/number571/hidden-lake/internal/adapters"
)

const (
	CAdapterName    = "tcp"
	CAppShortName   = adapters.CAdapterShortName + "-" + CAdapterName
	CAppFullName    = adapters.CAdapterFullName + "=" + CAdapterName
	CAppDescription = "adapts HL traffic to a custom TCP connection"
)

const (
	CPathYML = CAppShortName + ".yml"
	CPathDB  = CAppShortName + ".db"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9521"
	CDefaultInternalAddress = "127.0.0.1:9522"
)
