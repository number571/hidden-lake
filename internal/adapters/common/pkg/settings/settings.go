package settings

import "github.com/number571/hidden-lake/internal/utils/name"

var (
	GServiceName = name.LoadServiceName(CServiceFullName)
)

const (
	CServiceFullName    = "hidden-lake-adapter=common"
	CServiceDescription = "adapts HL traffic to a custom HTTP server"
)

const (
	CPathYML = "hla_common.yml"
)

const (
	CDefaultSrvAddress  = "http://127.0.0.1:6060"
	CDefaultHTTPAddress = "127.0.0.1:9531"
	CDefaultWaitTimeMS  = 1_000 // 1 second
)
