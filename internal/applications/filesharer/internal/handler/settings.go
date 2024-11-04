package handler

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/number571/go-peer/pkg/crypto/hashing"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/filesharer/pkg/app/config"
	hlf_settings "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"
	"github.com/number571/hidden-lake/internal/applications/filesharer/web"
	"github.com/number571/hidden-lake/internal/utils/language"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"

	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
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
	pHlsClient hls_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlf_settings.CServiceName, pR)

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
			port := strings.TrimSpace(pR.FormValue("port"))

			if host == "" || port == "" {
				ErrorPage(pLogger, cfg, "get_host_port", "host or port is nil")(pW, pR)
				return
			}
			if _, err := strconv.Atoi(port); err != nil {
				ErrorPage(pLogger, cfg, "port_to_int", "port is not a number")(pW, pR)
				return
			}

			connect := fmt.Sprintf("%s:%s", host, port)
			if err := pHlsClient.AddConnection(pCtx, connect); err != nil {
				ErrorPage(pLogger, cfg, "add_connection", "add connection")(pW, pR)
				return
			}
		case http.MethodDelete:
			connect := strings.TrimSpace(pR.FormValue("address"))
			if connect == "" {
				ErrorPage(pLogger, cfg, "get_connection", "connect is nil")(pW, pR)
				return
			}

			if err := pHlsClient.DelConnection(pCtx, connect); err != nil {
				ErrorPage(pLogger, cfg, "del_connection", "delete connection")(pW, pR)
				return
			}
		}

		result, err := getSettings(pCtx, cfg, pHlsClient)
		if err != nil {
			ErrorPage(pLogger, cfg, "get_settings", "get settings")(pW, pR)
			return
		}

		t, err := template.ParseFS(
			web.GetTemplatePath(),
			"index.html",
			"settings.html",
		)
		if err != nil {
			panic("can't load hmtl files")
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = t.Execute(pW, result)
	}
}

func getSettings(
	pCtx context.Context,
	pCfg config.IConfig,
	pHlsClient hls_client.IClient,
) (*sSettings, error) {
	result := new(sSettings)
	result.sTemplate = getTemplate(pCfg)

	myPubKey, err := pHlsClient.GetPubKey(pCtx)
	if err != nil {
		return nil, errors.Join(ErrGetPublicKey, err)
	}

	myPubKeyStr := myPubKey.ToString()
	result.FPublicKey = myPubKeyStr

	// echo PubKey{...} | sha384sum
	result.FPublicKeyHash = hashing.NewHasher([]byte(myPubKeyStr)).ToString()

	gotSettings, err := pHlsClient.GetSettings(pCtx)
	if err != nil {
		return nil, errors.Join(ErrGetSettingsHLS, err)
	}
	result.FNetworkKey = gotSettings.GetNetworkKey()

	allConns, err := getAllConnections(pCtx, pHlsClient)
	if err != nil {
		return nil, errors.Join(ErrGetAllConnections, err)
	}
	result.FConnections = allConns

	return result, nil
}

func getAllConnections(
	pCtx context.Context,
	pClient hls_client.IClient,
) ([]sConnection, error) {
	conns, err := pClient.GetConnections(pCtx)
	if err != nil {
		return nil, errors.Join(ErrReadConnections, err)
	}

	onlines, err := pClient.GetOnlines(pCtx)
	if err != nil {
		return nil, errors.Join(ErrReadOnlineConnections, err)
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
