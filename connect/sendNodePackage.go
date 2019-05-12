package connect

import (
    "net"
    "encoding/json"
    "../utils"
    "../models"
    "../settings"
)

func sendNodePackage(to string, pack models.PackageTCP, mode settings.ModeNet) int8 {
    var address string
    switch mode {
        case settings.P2P_mode: address = settings.Node.Address.P2P[to]
        case settings.F2F_mode: address = settings.Node.Address.F2F[to]
    }

    conn, err := net.Dial(settings.PROTOCOL, address)
    if err != nil {
        nullNode(to)
        return settings.EXIT_FAILED
    }

    data, err := json.Marshal(pack)
    utils.CheckError(err)

    conn.Write(data)
    conn.Close()

    return settings.EXIT_SUCCESS
}
