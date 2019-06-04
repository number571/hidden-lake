package connect

import (
	"net"
	"encoding/json"
	"../utils"
    "../models"
	"../settings"
)

// Send package by address node.
func sendPackageByAddr(to string, pack models.PackageTCP) int8 {
    conn, err := net.Dial(settings.PROTOCOL, to)
    if err != nil {
        return settings.EXIT_FAILED
    }

    data, err := json.Marshal(pack)
    utils.CheckError(err)

    conn.Write(data)
    conn.Close()

    return settings.EXIT_SUCCESS
}
