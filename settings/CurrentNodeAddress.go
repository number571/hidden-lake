package settings

import (
	"../models"
)

func CurrentNodeAddress() map[string][]byte {
    switch User.Mode {
        case models.P2P_mode: return Node.SessionKey.P2P
        case models.F2F_mode: return Node.SessionKey.F2F
    	default: return Node.SessionKey.P2P
    	// map[string][]byte{
    	// 	Node.ConnServer.Hash: []byte(""),
    	// }
    }
}
