package settings

import "github.com/number571/hidden-lake/internal/utils/appname"

var (
	gAppName = appname.LoadAppName(CServiceFullName)
)

func GetAppName() appname.IAppName {
	return gAppName
}

const (
	CServiceFullName    = "hidden-lake-composite"
	CServiceDescription = "runs many HL applications as one application"
)

const (
	CPathYML = "hlc.yml"
)
