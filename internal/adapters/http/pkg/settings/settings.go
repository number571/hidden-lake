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
	CServiceFullName    = adapters.CServiceFullPrefix + "=" + CServiceName
	CServiceDescription = "adapts HL traffic to a custom HTTP connection"
)

const (
	CPathYML = adapters.CServiceShortPrefix + "-" + CServiceName + ".yml"
	CPathDB  = adapters.CServiceShortPrefix + "-" + CServiceName + ".db"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9511"
	CDefaultInternalAddress = "127.0.0.1:9512"
)
