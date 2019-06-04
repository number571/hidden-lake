package settings

import (
	"../models"
)

func CurrentMode() string {
	return GetStrMode(User.Mode)
}

func GetStrMode(mode models.ModeNet) string {
	switch mode {
        case models.P2P_mode: return "P2P"
        case models.F2F_mode: return "F2F"
    	default: return "C-S"
    }
}
