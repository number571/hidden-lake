package connect 

import (
	"net"
	"../settings"
)

func GetConnection(address string) net.Conn {
	conn, err := net.Dial(settings.PROTOCOL, address)
	if err != nil {
		return nil
	}
	return conn
}
