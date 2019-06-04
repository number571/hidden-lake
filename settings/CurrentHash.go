package settings

import (
    "../models"
)

func CurrentHash() string {
    switch User.Mode {
        case models.P2P_mode: return User.Hash.P2P
        case models.F2F_mode: return User.Hash.F2F
    	default: return User.Hash.P2P
    }
}
