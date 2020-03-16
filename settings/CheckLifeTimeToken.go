package settings

import (
	"errors"
	"github.com/number571/hiddenlake/utils"
)

func CheckLifetimeToken(token string) error {
	user := Users[token]
	tokenTime := utils.ParseTime(user.Session.Time)
	currTime := utils.ParseTime(utils.CurrentTime())
	if tokenTime.Add(LIFETIME).Before(currTime) {
		delete(Listener.Clients, user.Hashname)
		delete(Tokens, user.Hashname)
		delete(Users, token)
		return errors.New("Token lifetime is over")
	}
	return nil
}
