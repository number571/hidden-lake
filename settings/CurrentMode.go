package settings

func CurrentMode() string {
    if User.ModeF2F {
        return "F2F"
    } else {
        return "P2P"
    }
}

func CurrentNodeAddress() map[string]string {
    if User.ModeF2F {
        return User.NodeAddressF2F
    } else {
        return User.NodeAddress
    }
}
