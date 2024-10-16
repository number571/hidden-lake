package settings

const (
	CServiceName     = "HLM"
	CServiceFullName = "hidden-lake-messenger"
)

const (
	CPathYML = "hlm.yml"
	CPathDB  = "hlm.db"
)

const (
	CStaticPath = "/static/"
	CPingPath   = "/ping"
	CPushPath   = "/push"
)

const (
	CDefaultInterfaceAddress = "127.0.0.1:9591"
	CDefaultIncomingAddress  = "127.0.0.1:9592"
)

const (
	CDefaultMessagesCapacity = (2 << 10) // count
	CDefaultLanguage         = ""        // ENG
)

const (
	CHandleIndexPath         = "/"
	CHandleAboutPath         = "/about"
	CHandleFaviconPath       = "/favicon.ico"
	CHandleSettingsPath      = "/settings"
	CHandleFriendsPath       = "/friends"
	CHandleFriendsChatPath   = "/friends/chat"
	CHandleFriendsUploadPath = "/friends/upload"
	CHandleFriendsChatWSPath = "/friends/chat/ws"
)

const (
	CIsText = 0x01
	CIsFile = 0x02
)
