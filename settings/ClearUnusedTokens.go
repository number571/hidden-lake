package settings

import (
	"time"
)

func ClearUnusedTokens(checkTime time.Duration) {
	for {
		time.Sleep(checkTime)
		for token := range Users {
			CheckLifetimeToken(token)
		}
	}
}
