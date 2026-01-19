package settings

import (
	"strings"

	"github.com/number571/hidden-lake/internal/adapters"
)

var (
	gAppShortNameFMT = strings.ToUpper(CAppShortName)
)

func GetAppShortNameFMT() string {
	return gAppShortNameFMT
}

const (
	CAppShortName = adapters.CAdapterShortName + "-" + CAppAdapterName
	CAppFullName  = adapters.CAdapterFullName + "=" + CAppAdapterName
)

const (
	CAppAdapterName = "http"
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

const (
	CHandleIndexPath          = "/api/index"
	CHandleConfigSettingsPath = "/api/config/settings"
	CHandleConfigConnectsPath = "/api/config/connects"
	CHandleNetworkOnlinePath  = "/api/network/online"
	CHandleNetworkAdapterPath = "/api/network/adapter"
)
