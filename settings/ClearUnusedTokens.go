package settings

import (
	"time"
)

func ClearUnusedTokens() {
	for {
		time.Sleep(CHECK_DURING)
		for token := range Users {
			CheckLifetimeToken(token)
		}
	}
}
