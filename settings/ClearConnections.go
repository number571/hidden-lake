package settings

import (
	"github.com/number571/gopeer"
	"time"
)

func ClearConnections(checkTime time.Duration) {
	for {
		time.Sleep(checkTime)
		sendTestPackages()
	}
}

func sendTestPackages() {
	for _, client := range Listener.Clients {
		for hash := range client.Connections {
			dest := client.Destination(hash)
			client.SendTo(dest, &gopeer.Package{
				Head: gopeer.Head{
					Title:  TITLE_TESTCONN,
					Option: gopeer.Get("OPTION_GET").(string),
				},
			})
			select {
			case <-client.Connections[hash].Chans.Action:
				// pass
			case <-time.After(time.Duration(gopeer.Get("WAITING_TIME").(uint8)) * time.Second):
				delete(client.Connections, hash)
			}
		}
	}
}
