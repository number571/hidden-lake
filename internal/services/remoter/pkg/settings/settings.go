package settings

import (
	"github.com/number571/hidden-lake/internal/services"
	"github.com/number571/hidden-lake/internal/utils/name"
)

var (
	gServiceName = name.LoadServiceName(CServiceFullName)
)

func GetServiceName() name.IServiceName {
	return gServiceName
}

const (
	CServiceName = "remoter"
)

const (
	CServiceFullName    = services.CServiceFullPrefix + "=" + CServiceName
	CServiceDescription = "executes remote access commands"
)

const (
	CPathYML        = services.CServiceShortPrefix + "-" + CServiceName + ".yml"
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
