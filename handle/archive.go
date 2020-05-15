package handle

import (
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/db"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/settings"
)

func getArchive(client *gopeer.Client, pack *gopeer.Package) (set string) {
	var (
		token = settings.Tokens[client.Hashname()]
		user  = settings.Users[token]
	)
	if !user.State.UsedFSH {
		return ""
	}
	if pack.Body.Data == "" {
		var result []models.File
		files := db.GetAllFiles(user)
		for _, file := range files {
			if file.Encr {
				continue
			}
			result = append(result, file)
		}
		return string(gopeer.PackJSON(result))
	}
	file := db.GetFile(user, pack.Body.Data)
	if file == nil {
		return ""
	}
	return string(gopeer.PackJSON([]models.File{*file}))
}

func setArchive(client *gopeer.Client, pack *gopeer.Package) {
	var (
		token = settings.Tokens[client.Hashname()]
		user  = settings.Users[token]
	)
	gopeer.UnpackJSON([]byte(pack.Body.Data), &user.Temp.FileList)
	client.Connections[pack.From.Sender.Hashname].Action <- true
}
