package settings

import "github.com/number571/hidden-lake/internal/utils/name"

var (
	gServiceName = name.LoadServiceName(CServiceFullName)
)

func GetServiceName() name.IServiceName {
	return gServiceName
}

const (
	CServiceFullName    = "hidden-lake-composite"
	CServiceDescription = "runs many HL services as one application"
)

const (
	CPathYML = "hlc.yml"
)
