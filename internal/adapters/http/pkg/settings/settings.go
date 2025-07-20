package settings

import (
	"github.com/number571/hidden-lake/internal/adapters"
	"github.com/number571/hidden-lake/internal/utils/name"
)

var (
	gServiceName = name.LoadServiceName(CServiceFullName)
)

func GetServiceName() name.IServiceName {
	return gServiceName
}

const (
	CServiceName = "http"
)

const (
	CServiceFullName    = adapters.CServicePrefix + "=" + CServiceName
	CServiceDescription = "adapts HL traffic to a custom HTTP connection"
)

const (
	CPathYML = "hla-" + CServiceName + ".yml"
	CPathDB  = "hla-" + CServiceName + ".db"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9511"
	CDefaultInternalAddress = "127.0.0.1:9512"
)
