package settings

import (
	"github.com/number571/hidden-lake/internal/utils/name"
)

var (
	GServiceName = name.LoadServiceName(CServiceFullName)
)

const (
	CServiceFullName    = "hidden-lake-adapter=tcp"
	CServiceDescription = "adapts HL traffic to a custom TCP connection"
)

const (
	CPathYML = "hla_tcp.yml"
	CPathDB  = "hla_tcp.db"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9521"
	CDefaultInternalAddress = "127.0.0.1:9522"
)
