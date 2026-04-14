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
	CAppAdapterName = "https"
	CAppDescription = "adapts HL traffic to a custom HTTPS connection"
)

const (
	CPathYML   = CAppShortName + ".yml"
	CPathDB    = CAppShortName + ".db"
	CPathCerts = CAppShortName + ".certs"
	CPathKey   = CAppShortName + ".key"
	CPathCert  = CAppShortName + ".cert"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9531"
	CDefaultInternalAddress = "127.0.0.1:9532"
)

const (
	CHandleAdapterProducePath = "/adapter/produce"
	CHandleAdapterConsumePath = "/adapter/consume"
)

const (
	CPasswordHeader = "Password"
)
