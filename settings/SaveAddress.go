package settings

import (
	"../utils"
	"../crypto"
)

func SaveAddress(ipv4, port string) {
    Mutex.Lock()
    _, err := DataBase.Exec(
        "INSERT INTO Address (IPv4, Port) VALUES ($1, $2)", 
        crypto.Encrypt(User.Password, ipv4), 
        crypto.Encrypt(User.Password, port),
    )
    utils.CheckError(err)
    Mutex.Unlock()
}
