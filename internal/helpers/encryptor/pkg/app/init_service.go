package app

import (
	"net/http"
	"time"

	"github.com/number571/go-peer/pkg/client"
	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/internal/handler"
	"github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/app/config"
	hle_settings "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/settings"
)

func initMapPubKeys(pCfg config.IConfig) asymmetric.IMapPubKeys {
	f2f := asymmetric.NewMapPubKeys()
	for _, pubKey := range pCfg.GetFriends() {
		f2f.SetPubKey(pubKey)
	}
	return f2f
}

func (p *sApp) initServiceHTTP() {
	mux := http.NewServeMux()
	cfg := p.fCfgW.GetConfig()

	mapKeys := initMapPubKeys(cfg)
	client := client.NewClient(p.fPrivKey, cfg.GetSettings().GetMessageSizeBytes())

	mux.HandleFunc(hle_settings.CHandleIndexPath, handler.HandleIndexAPI(p.fHTTPLogger))
	mux.HandleFunc(hle_settings.CHandleMessageEncryptPath, handler.HandleMessageEncryptAPI(cfg, p.fHTTPLogger, client, p.fParallel))
	mux.HandleFunc(hle_settings.CHandleMessageDecryptPath, handler.HandleMessageDecryptAPI(cfg, p.fHTTPLogger, client, mapKeys))
	mux.HandleFunc(hle_settings.CHandleServicePubKeyPath, handler.HandleServicePubKeyAPI(p.fHTTPLogger, client.GetPrivKey().GetPubKey()))
	mux.HandleFunc(hle_settings.CHandleConfigSettingsPath, handler.HandleConfigSettingsAPI(cfg, p.fHTTPLogger))
	mux.HandleFunc(hle_settings.CHandleConfigFriendsPath, handler.HandleConfigFriendsAPI(p.fCfgW, p.fHTTPLogger, mapKeys))

	p.fServiceHTTP = &http.Server{
		Addr:        cfg.GetAddress().GetInternal(),
		Handler:     http.TimeoutHandler(mux, time.Minute, "timeout"),
		ReadTimeout: (5 * time.Second),
	}
}
