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
	CAppAdapterName = "tcp"
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
