package settings

func CurrentModeNet() ModeNet {
	if User.ModeF2F {
		return F2F_mode
	} else {
		return P2P_mode
	}
}
