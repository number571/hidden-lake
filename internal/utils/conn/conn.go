package conn

import (
	"strconv"
	"strings"

	hiddenlake "github.com/number571/hidden-lake"
)

type IAddress interface {
	GetTCP() string
}

func IsAmI(pAddr IAddress, conn hiddenlake.SConnection) bool {
	splited := strings.Split(pAddr.GetTCP(), ":")
	if len(splited) < 2 {
		return false
	}
	tcpPort, _ := strconv.Atoi(splited[1])
	if conn.FHost == "localhost" || conn.FHost == "127.0.0.1" {
		if conn.FPort == uint16(tcpPort) {
			return true
		}
	}
	return false
}
