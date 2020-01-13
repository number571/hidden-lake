package settings

import (
	"../models"
	"../utils"
	"encoding/json"
)

func InitializeCFG(cfgname string) {
	if !utils.FileIsExist(cfgname) {
		cfgJSON, err := json.MarshalIndent(newConfig(), "", "\t")
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
}

func newConfig() *models.Config {
	return &models.Config{
		Host: models.Host{
			Http: models.Http{
				Ipv4: "localhost",
				Port: ":7545",
				Tls: models.Tls{
					Crt: "",
					Key: "",
				},
			},
			Tcp: models.Tcp{
				Ipv4: "localhost",
				Port: ":8080",
			},
		},
	}
}
