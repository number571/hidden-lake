package settings

import (
	"github.com/number571/hidden-lake/internal/adapters"
)

const (
	CAdapterName    = "http"
	CAppShortName   = adapters.CAdapterShortName + "-" + CAdapterName
	CAppFullName    = adapters.CAdapterFullName + "=" + CAdapterName
	CAppDescription = "adapts HL traffic to a custom HTTP connection"
)

const (
	CPathYML = CAppShortName + ".yml"
	CPathDB  = CAppShortName + ".db"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9511"
	CDefaultInternalAddress = "127.0.0.1:9512"
)
