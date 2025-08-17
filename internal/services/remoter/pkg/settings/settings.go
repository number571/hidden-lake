package settings

import (
	"github.com/number571/hidden-lake/internal/services"
	"github.com/number571/hidden-lake/internal/utils/appname"
)

var (
	gShortAppName = appname.ToShortAppName(CAppFullName)
)

func GetShortAppName() string {
	return gShortAppName
}

const (
	CAppName = "remoter"
)

const (
	CAppFullName    = services.CServiceFullPrefix + "=" + CAppName
	CAppDescription = "executes remote access commands"
)

const (
	CPathYML        = services.CServiceShortPrefix + "-" + CAppName + ".yml"
	CHeaderPassword = "Hl-Remoter-Password" // nolint: gosec
)

const (
	CExecPath      = "/exec"
	CExecSeparator = "[@remoter-separator]"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9532"
	CDefaultExecTimeout     = 5_000 // 5s
)
