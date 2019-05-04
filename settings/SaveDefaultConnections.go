package settings

import (
	"../utils"
	"../crypto"
)

func SaveDefaultConnections(connections []string) {
	for _, value := range connections {
        Mutex.Lock()
        _, err := DataBase.Exec(
            "INSERT INTO DefaultConnections (Address) VALUES ($1)", 
            crypto.Encrypt(User.Password, value),
        )
        Mutex.Unlock()
        utils.CheckError(err)
    }
}
