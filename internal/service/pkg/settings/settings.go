package settings

const (
	CServiceName     = "HLS"
	CServiceFullName = "hidden-lake-service"
)

const (
	CNetworkMask = uint32(0x5f67705f) // bytes_prefix: _gp_
	CServiceMask = uint32(0x5f686c5f) // bytes_prefix: _hl_
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
	CDefaultTCPAddress  = "127.0.0.1:9571"
	CDefaultHTTPAddress = "127.0.0.1:9572"
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
