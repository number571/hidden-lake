package settings

import (
	"strings"

	"github.com/number571/hidden-lake/internal/services"
)

var (
	gAppShortNameFMT = strings.ToUpper(CAppShortName)
)

func GetAppShortNameFMT() string {
	return gAppShortNameFMT
}

const (
	CServiceName    = "remoter"
	CAppShortName   = services.CServiceShortName + "-" + CServiceName
	CAppFullName    = services.CServiceDomainName + "=" + CServiceName
	CAppDescription = "executes remote access commands"
)

const (
	CPathYML = CAppShortName + ".yml"
)

const (
	CHeaderPassword = "Password" // nolint: gosec
	CExecSeparator  = "[@remoter-separator]"
)

const (
	CExecPath = "/exec"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9532"
	CDefaultExecTimeout     = 5_000 // 5s
)
