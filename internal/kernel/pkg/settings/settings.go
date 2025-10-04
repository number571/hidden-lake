package settings

import "strings"

var (
	gAppShortNameFMT = strings.ToUpper(CAppShortName)
)

func GetAppShortNameFMT() string {
	return gAppShortNameFMT
}

const (
	CAppShortName = "hlk"
	CAppFullName  = "hidden-lake-kernel"
)

const (
	CAppDescription = "anonymizes traffic using the QB-problem"
)

const (
	CPathKey = CAppShortName + ".key"
	CPathYML = CAppShortName + ".yml"
	CPathDB  = CAppShortName + ".db"
)

const (
	CHeaderSenderFriend = "Hlk-Sender-Friend"
	CHeaderResponseMode = "Hlk-Response-Mode"
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
