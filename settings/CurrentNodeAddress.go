package settings

func CurrentNodeAddress() map[string]string {
    if User.ModeF2F {
        return Node.Address.F2F
    } else {
        return Node.Address.P2P
    }
}
