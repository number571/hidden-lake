package settings

func CurrentHash() string {
    if User.ModeF2F {
        return User.Hash.F2F
    } else {
        return User.Hash.P2P
    }
}
