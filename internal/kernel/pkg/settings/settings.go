package settings

import "github.com/number571/hidden-lake/pkg/utils/appname"

var (
	gAppName = appname.LoadAppName(CServiceFullName)
)

func GetAppName() appname.IAppName {
	return gAppName
}

const (
	CServiceFullName    = "hidden-lake-kernel"
	CServiceDescription = "anonymizes traffic using the QB-problem"
)

const (
	CPathKey = "hlk.key"
	CPathYML = "hlk.yml"
	CPathDB  = "hlk.db"
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
	CDefaultInternalAddress = "127.0.0.1:9572"
)

const (
	CHandleIndexPath          = "/api/index"
	CHandleConfigSettingsPath = "/api/config/settings"
	CHandleConfigConnectsPath = "/api/config/connects"
	CHandleConfigFriendsPath  = "/api/config/friends"
	CHandleNetworkOnlinePath  = "/api/network/online"
	CHandleNetworkRequestPath = "/api/network/request"
	CHandleServicePubKeyPath  = "/api/kernel/pubkey"
)
