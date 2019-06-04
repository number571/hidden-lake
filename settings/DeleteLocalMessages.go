package settings

import (
    // "net"
    "../utils"
    // "../models"
)

func DeleteLocalMessages(slice []string) {
	// var (
 //        node_address map[string]string
 //    )
 //    switch User.Mode {
 //        case models.P2P_mode: node_address = Node.Address.P2P
 //        case models.F2F_mode: node_address = Node.Address.F2F
 //    }

    Mutex.Lock()
    for _, user := range slice {
        // if _, ok := node_address[user]; ok {
        _, err := DataBase.Exec("DELETE FROM Local" + user)
        utils.CheckError(err)
        // }
    }
    Mutex.Unlock()
}
