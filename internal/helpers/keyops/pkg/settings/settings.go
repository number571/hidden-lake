package settings

import (
	"strings"

	"github.com/number571/hidden-lake/internal/helpers"
)

var (
	gAppShortNameFMT = strings.ToUpper(CAppShortName)
)

func GetAppShortNameFMT() string {
	return gAppShortNameFMT
}

const (
	CPathKey    = CAppShortName + ".key"
	CPathPubKey = CPathKey + ".pub"
	CPathEncKey = CAppShortName + ".enc"
)

const (
	CAppShortName = helpers.CHelperShortName + "-" + CAppHelperName
	CAppFullName  = helpers.CHelperDomainName + "=" + CAppHelperName
)

const (
	CAppHelperName  = "keyops"
	CAppDescription = "generate / encapsulate / decapsulate key operations"
)
