package settings

import (
	"../models"
)

func CurrentModeNet() models.ModeNet {
	return User.Mode
}
