package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/app/config"
	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/language"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/utils/pubkey"
	"github.com/number571/hidden-lake/internal/webui"

	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
)

type sConnection struct {
	FAddress  string
	FIsBackup bool
	FOnline   bool
}

type sSettings struct {
	*sTemplate
	FNetworkKey    string
	FPublicKey     string
	FPublicKeyHash string
	FConnections   []sConnection
}

func SettingsPage(
	pCtx context.Context,
	pLogger logger.ILogger,
	pWrapper config.IWrapper,
	pHlkClient hlk_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hls_messenger_settings.GetAppShortNameFMT(), pR)

		cfg := pWrapper.GetConfig()
		cfgEditor := pWrapper.GetEditor()

		if pR.URL.Path != "/settings" {
			NotFoundPage(pLogger, cfg)(pW, pR)
			return
		}

		if err := pR.ParseForm(); err != nil {
			ErrorPage(pLogger, cfg, "parse_form", "parse form")(pW, pR)
			return
		}

		switch pR.FormValue("method") {
		case http.MethodPut:
			strLang := strings.TrimSpace(pR.FormValue("language"))
			ilang, err := language.ToILanguage(strLang)
			if err != nil {
				ErrorPage(pLogger, cfg, "to_language", "load unknown language")(pW, pR)
				return
			}
			if err := cfgEditor.UpdateLanguage(ilang); err != nil {
				ErrorPage(pLogger, cfg, "update_language", "update language")(pW, pR)
				return
			}
		case http.MethodPost:
			host := strings.TrimSpace(pR.FormValue("host"))
			if host == "" {
				ErrorPage(pLogger, cfg, "get_host", "host is nil")(pW, pR)
				return
			}
			if err := pHlkClient.AddConnection(pCtx, host); err != nil {
				ErrorPage(pLogger, cfg, "add_connection", "add connection")(pW, pR)
				return
			}
		case http.MethodDelete:
			connect := strings.TrimSpace(pR.FormValue("address"))
			if connect == "" {
				ErrorPage(pLogger, cfg, "get_connection", "connect is nil")(pW, pR)
				return
			}

			if err := pHlkClient.DelConnection(pCtx, connect); err != nil {
				ErrorPage(pLogger, cfg, "del_connection", "delete connection")(pW, pR)
				return
			}
		}

		result, err := getSettings(pCtx, cfg, pHlkClient)
		if err != nil {
			ErrorPage(pLogger, cfg, "get_settings", "get settings")(pW, pR)
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = webui.MustParseTemplate("index.html", "settings.html").Execute(pW, result)
	}
}

func getSettings(pCtx context.Context, pCfg config.IConfig, pHlkClient hlk_client.IClient) (*sSettings, error) {
	result := new(sSettings)
	result.sTemplate = getTemplate(pCfg)

	myPubKey, err := pHlkClient.GetPubKey(pCtx)
	if err != nil {
		return nil, errors.Join(ErrGetPublicKey, err)
	}

	result.FPublicKey = myPubKey.ToString()
	result.FPublicKeyHash = pubkey.GetPubKeyHash(myPubKey)

	gotSettings, err := pHlkClient.GetSettings(pCtx)
	if err != nil {
		return nil, errors.Join(ErrGetSettings, err)
	}
	result.FNetworkKey = gotSettings.GetNetworkKey()

	allConns, err := getAllConnections(pCtx, pHlkClient)
	if err != nil {
		return nil, errors.Join(ErrGetAllConnections, err)
	}
	result.FConnections = allConns

	return result, nil
}

func getAllConnections(pCtx context.Context, pClient hlk_client.IClient) ([]sConnection, error) {
	conns, err := pClient.GetConnections(pCtx)
	if err != nil {
		return nil, ErrReadConnections
	}

	onlines, err := pClient.GetOnlines(pCtx)
	if err != nil {
		return nil, ErrReadOnlineConnections
	}

	connections := make([]sConnection, 0, len(conns))
	for _, c := range conns {
		connections = append(
			connections,
			sConnection{
				FAddress: c,
				FOnline:  getOnline(onlines, c),
			},
		)
	}

	return connections, nil
}

func getOnline(onlines []string, c string) bool {
	for _, o := range onlines {
		if o == c {
			return true
		}
	}
	return false
}
