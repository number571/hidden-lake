package settings

import "github.com/number571/hidden-lake/internal/utils/name"

var (
	gAppName = name.LoadAppName(CServiceFullName)
)

func GetAppName() name.IAppName {
	return gAppName
}

const (
	CServiceFullName    = "hidden-lake-composite"
	CServiceDescription = "runs many HL applications as one application"
)

const (
	CPathYML = "hlc.yml"
)
