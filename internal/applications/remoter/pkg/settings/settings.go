package settings

import "github.com/number571/hidden-lake/internal/utils/name"

var (
	gServiceName = name.LoadServiceName(CServiceFullName)
)

func GetServiceName() name.IServiceName {
	return gServiceName
}

const (
	CServiceFullName    = "hidden-lake-remoter"
	CServiceDescription = "executes remote access commands"
)

const (
	CPathYML        = "hlr.yml"
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
