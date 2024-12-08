package settings

import "github.com/number571/hidden-lake/internal/utils/name"

var GServiceName = name.LoadServiceName(CServiceFullName)

const (
	CServiceFullName    = "hidden-lake-service"
	CServiceDescription = "anonymizes traffic using the QB-problem"
)

const (
	CPathKey = "hls.key"
	CPathYML = "hls.yml"
	CPathDB  = "hls.db"
)

const (
	CHeaderPublicKey    = "Hl-Service-Public-Key"
	CHeaderResponseMode = "Hl-Service-Response-Mode"
)

const (
	CHeaderResponseModeON  = "on" // default
	CHeaderResponseModeOFF = "off"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9571"
	CDefaultHTTPAddress     = "127.0.0.1:9572"
)

const (
	CHandleIndexPath          = "/api/index"
	CHandleConfigSettingsPath = "/api/config/settings"
	CHandleConfigConnectsPath = "/api/config/connects"
	CHandleConfigFriendsPath  = "/api/config/friends"
	CHandleNetworkOnlinePath  = "/api/network/online"
	CHandleNetworkRequestPath = "/api/network/request"
	CHandleServicePubKeyPath  = "/api/service/pubkey"
)
