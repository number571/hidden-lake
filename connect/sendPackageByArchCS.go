package connect

import (
	"net"
	"encoding/json"
	"../utils"
    "../models"
	"../settings"
)

// Send package by address node.
func sendPackageByArchCS(conn net.Conn, pack models.PackageTCP) int8 {
    if conn == nil {
        return settings.EXIT_FAILED
    }
    data, err := json.Marshal(pack)
    utils.CheckError(err)
    conn.Write(append(data, []byte(settings.END_BLOCK)...))
    return settings.EXIT_SUCCESS
}
