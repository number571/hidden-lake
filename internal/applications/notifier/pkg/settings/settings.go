package settings

import "github.com/number571/hidden-lake/internal/utils/name"

var (
	GServiceName = name.LoadServiceName(CServiceFullName)
)

const (
	CServiceFullName    = "hidden-lake-notifier"
	CServiceDescription = "broadcast message throw crowds protocol"
)

const (
	CStaticPath = "/static/"
	CPathYML    = "hln.yml"
	CPathDB     = "hln.db"
)

const (
	CFinalyzePath = "/finalyze"
	CRedirectPath = "/redirect"
)

const (
	CHandleIndexPath          = "/"
	CHandleAboutPath          = "/about"
	CHandleFaviconPath        = "/favicon.ico"
	CHandleSettingsPath       = "/settings"
	CHandleFriendsPath        = "/friends"
	CHandleChannelsPath       = "/channels"
	CHandleChannelsChatWSPath = "/channels/chat/ws"
	CHandleChannelsChatPath   = "/channels/chat"
	CHandleChannelsUploadPath = "/channels/upload"
)

const (
	CDefaultInternalAddress = "127.0.0.1:9561"
	CDefaultExternalAddress = "127.0.0.1:9562"
)

const (
	CDefaultMessagesCapacity = (2 << 10) // count
	CDefaultLanguage         = ""        // ENG
)
