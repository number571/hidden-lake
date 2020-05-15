package handle

import (
	"github.com/number571/gopeer"
)

func getTestconn(client *gopeer.Client, pack *gopeer.Package) (set string) {
	return set
}

func setTestconn(client *gopeer.Client, pack *gopeer.Package) {
	client.Connections[pack.From.Sender.Hashname].Action <- true
}
