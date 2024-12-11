package settings

import "github.com/number571/hidden-lake/internal/utils/name"

var (
	GServiceName = name.LoadServiceName(CServiceFullName)
)

const (
	CServiceFullName    = "hidden-lake-encryptor"
	CServiceDescription = "encrypts and decrypts messages"
)

const (
	CPathKey = "hle.key"
	CPathYML = "hle.yml"
)

const (
	CDefaultInternalAddress = "127.0.0.1:9551"
)

const (
	CHandleIndexPath          = "/api/index"
	CHandleMessageEncryptPath = "/api/message/encrypt"
	CHandleMessageDecryptPath = "/api/message/decrypt"
	CHandleServicePubKeyPath  = "/api/service/pubkey"
	CHandleConfigSettingsPath = "/api/config/settings"
	CHandleConfigFriendsPath  = "/api/config/friends"
)
