package settings

import (
	"github.com/number571/hidden-lake/internal/adapters"
	"github.com/number571/hidden-lake/internal/utils/name"
)

var (
	GServiceName = name.LoadServiceName(CServiceFullName)
)

const (
	CServiceAdapterScheme = "tcp"
)

const (
	CServiceFullName    = adapters.CServicePrefix + "=" + CServiceAdapterScheme
	CServiceDescription = "adapts HL traffic to a custom TCP connection"
)

const (
	CPathYML = "hla_" + CServiceAdapterScheme + ".yml"
	CPathDB  = "hla_" + CServiceAdapterScheme + ".db"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9521"
	CDefaultInternalAddress = "127.0.0.1:9522"
)
