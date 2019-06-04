package connect 

import (
	// "fmt"
	"net"
	"strings"
	"../settings"
)

func readConnection(conn net.Conn) {
	if settings.ConnectionIsRead { return }
	settings.ConnectionIsRead = true
	var (
		buffer = make([]byte, settings.BUFF_SIZE)
		message string
	)
	for {
		length, err := conn.Read(buffer)
        if length == 0 || err != nil { break }
        message += string(buffer[:length])

        if strings.HasSuffix(message, settings.END_BLOCK) {
            message = strings.TrimSuffix(message, settings.END_BLOCK)
            serve(conn, message)
            message = ""
        }
	}
	settings.ConnectionIsRead = false
}
