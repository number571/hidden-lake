package settings

import "github.com/number571/hidden-lake/internal/utils/name"

var (
	GServiceName = name.LoadServiceName(CServiceFullName)
)

const (
	CServiceFullName    = "hidden-lake-traffic"
	CServiceDescription = "retransmits and saves encrypted traffic"
)

const (
	CPathYML = "hlt.yml"
	CPathDB  = "hlt.db"
)

const (
	CHandleIndexPath          = "/api/index"
	CHandleStoragePointerPath = "/api/storage/pointer"
	CHandleStorageHashesPath  = "/api/storage/hashes"
	CHandleNetworkMessagePath = "/api/network/message"
	CHandleConfigSettings     = "/api/config/settings"
)

const (
	CDefaultTCPAddress  = "127.0.0.1:9581"
	CDefaultHTTPAddress = "127.0.0.1:9582"
)

const (
	CDefaultConnectionAddress = "127.0.0.1:9571"
)
