package settings

import (
	"encoding/json"
	"github.com/number571/gopeer"
	"github.com/number571/hiddenlake/models"
	"github.com/number571/hiddenlake/utils"
	"os"
)

func InitializeCFG(cfgname string) {
	if !utils.FileIsExist(cfgname) {
		cfg := newConfig()
		cfgJSON, err := json.MarshalIndent(cfg, "", "\t")
		if err != nil {
			panic("can't encode config")
		}
		utils.WriteFile(cfgname, string(cfgJSON))
	}
	cfgJSON := utils.ReadFile(cfgname)
	err := json.Unmarshal([]byte(cfgJSON), &CFG)
	if err != nil {
		panic("can't decode config")
	}
	os.Mkdir(PATH_TLS, 0777)
	os.Mkdir(PATH_ARCHIVE, 0777)
	if !utils.FileIsExist(CFG.Tls.Crt) && !utils.FileIsExist(CFG.Tls.Key) {
		key, cert := gopeer.GenerateCertificate(gopeer.Get("SERVER_NAME").(string), gopeer.Get("KEY_SIZE").(uint16))
		utils.WriteFile(CFG.Tls.Crt, cert)
		utils.WriteFile(CFG.Tls.Key, key)
	}
}

func newConfig() *models.Config {
	return &models.Config{
		Tls: models.Tls{
			Crt: PATH_TLS + "cert.crt",
			Key: PATH_TLS + "cert.key",
		},
		Http: models.Http{
			Ipv4: "localhost",
			Port: ":7545",
		},
		Tcp: models.Tcp{
			Ipv4: "localhost",
			Port: ":8080",
		},
	}
}
