package settings

import (
	"github.com/number571/hidden-lake/internal/utils/name"
)

var (
	GServiceName = name.LoadServiceName(CServiceFullName)
)

const (
	CServiceAdapterScheme = "http"
)

const (
	CServiceFullName    = "hidden-lake-adapter=" + CServiceAdapterScheme
	CServiceDescription = "adapts HL traffic to a custom HTTP connection"
)

const (
	CPathYML = "hla_" + CServiceAdapterScheme + ".yml"
	CPathDB  = "hla_" + CServiceAdapterScheme + ".db"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9511"
	CDefaultInternalAddress = "127.0.0.1:9512"
)
