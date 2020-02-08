package db

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func GetAllClients(user *models.User) []models.Client {
	var (
		clients    []models.Client
		hashname   string
		address    string
		public     string
		publicrecv string 
	)
	id := GetUserId(user.Auth.Hashpasw)
	if id < 0 {
		return nil
	}
	rows, err := settings.DB.Query(
		"SELECT Hashname, Address, PublicKey, PublicRecv FROM Client WHERE IdUser=$1",
		id,
	)
	if err != nil {
		panic("query 'getallclients' failed")
	}
	for rows.Next() {
		err := rows.Scan(&hashname, &address, &public, &publicrecv)
		if err != nil {
			return nil
		}
		clients = append(clients, models.Client{
			Hashname: hashname,
			Address:  string(gopeer.DecryptAES(user.Auth.Pasw, gopeer.Base64Decode(address))),
			Public:   gopeer.ParsePublic(public),
			PublicRecv: gopeer.ParsePublic(publicrecv),
		})
	}
	return clients
}
