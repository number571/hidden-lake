package settings

func CurrentMode() string {
    if User.ModeF2F {
        return "F2F"
    } else {
        return "P2P"
    }
}
