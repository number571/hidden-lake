package handle

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/settings"
)

func Actions(client *gopeer.Client, pack *gopeer.Package) {
	client.HandleAction(settings.TITLE_EMAIL,   pack, getEmail,   setEmail)
	client.HandleAction(settings.TITLE_ARCHIVE, pack, getArchive, setArchive)
	client.HandleAction(settings.TITLE_MESSAGE, pack, getMessage, setMessage)
}
