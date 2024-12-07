package settings

import "github.com/number571/hidden-lake/internal/utils/name"

var (
	GServiceName = name.LoadServiceName(CServiceFullName)
)

const (
	CServiceFullName    = "hidden-lake-loader"
	CServiceDescription = "distributes the stored traffic between nodes"
)

const (
	CPathYML = "hll.yml"
)

const (
	CDefaultHTTPAddress = "127.0.0.1:9561"
)

const (
	CHandleIndexPath           = "/api/index"
	CHandleNetworkTransferPath = "/api/network/transfer"
	CHandleConfigSettings      = "/api/config/settings"
)
