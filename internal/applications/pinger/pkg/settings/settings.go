package settings

import "github.com/number571/hidden-lake/internal/utils/name"

var (
	gServiceName = name.LoadServiceName(CServiceFullName)
)

func GetServiceName() name.IServiceName {
	return gServiceName
}

const (
	CServiceFullName    = "hidden-lake-pinger"
	CServiceDescription = "ping the node to check the online status"
)

const (
	CPathYML = "hlp.yml"
)

const (
	CPingPath = "/ping"
)

const (
	CDefaultExternalAddress = "127.0.0.1:9552"
)
