package settings

import "github.com/number571/hidden-lake/internal/utils/appname"

var (
	gShortAppName = appname.ToShortAppName(CAppFullName)
)

func GetShortAppName() string {
	return gShortAppName
}

const (
	CAppFullName    = "hidden-lake-kernel"
	CAppDescription = "anonymizes traffic using the QB-problem"
)

const (
	CPathKey = "hlk.key"
	CPathYML = "hlk.yml"
	CPathDB  = "hlk.db"
)

const (
	CHeaderPublicKey    = "Hl-Public-Key"
	CHeaderResponseMode = "Hl-Response-Mode"
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
	CHandleProfilePubKeyPath  = "/api/profile/pubkey"
)
