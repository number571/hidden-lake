package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/number571/go-peer/pkg/crypto/asymmetric"
	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/go-peer/pkg/message/layer1"
	"github.com/number571/hidden-lake/internal/applications/notifier/pkg/app/config"
	hln_client "github.com/number571/hidden-lake/internal/applications/notifier/pkg/client"
	hln_settings "github.com/number571/hidden-lake/internal/applications/notifier/pkg/settings"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/alias"
	"github.com/number571/hidden-lake/internal/utils/api"
	"github.com/number571/hidden-lake/internal/utils/layer3"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
)

func HandleIncomingRedirectHTTP(
	pCtx context.Context,
	pConfig config.IConfig,
	pLogger logger.ILogger,
	pHLSClient hls_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		pW.Header().Set(hls_settings.CHeaderResponseMode, hls_settings.CHeaderResponseModeOFF)

		logBuilder := http_logger.NewLogBuilder(hln_settings.GServiceName.Short(), pR)

		if pR.Method != http.MethodPost {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogMethod))
			_ = api.Response(pW, http.StatusMethodNotAllowed, "failed: incorrect method")
			return
		}

		msgBytes, err := io.ReadAll(pR.Body)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage(http_logger.CLogDecodeBody))
			_ = api.Response(pW, http.StatusNotAcceptable, "failed: decode body")
			return
		}

		rawMsg, err := layer1.LoadMessage(
			layer1.NewSettings(&layer1.SSettings{
				FWorkSizeBits: pConfig.GetSettings().GetWorkSizeBits(),
			}),
			msgBytes,
		)
		if err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("decode_message"))
			_ = api.Response(pW, http.StatusNotAcceptable, "failed: decode message")
			return
		}

		if _, err := layer3.ExtractMessage(rawMsg); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("decode_message_body"))
			_ = api.Response(pW, http.StatusNotAcceptable, "failed: decode message body")
			return
		}

		friends, err := pHLSClient.GetFriends(pCtx)
		if err != nil || len(friends) == 0 {
			pLogger.PushWarn(logBuilder.WithMessage("get_friends"))
			_ = api.Response(pW, http.StatusNotAcceptable, "failed: get friends")
			return
		}

		fSender := asymmetric.LoadPubKey(pR.Header.Get(hls_settings.CHeaderPublicKey))
		if fSender == nil {
			pLogger.PushErro(logBuilder.WithMessage("load_pubkey"))
			_ = api.Response(pW, http.StatusForbidden, "failed: load public key")
			return
		}

		aliasName := alias.GetAliasByPubKey(friends, fSender)
		if aliasName == "" {
			pLogger.PushErro(logBuilder.WithMessage("find_alias"))
			_ = api.Response(pW, http.StatusForbidden, "failed: find alias")
			return
		}

		hlnClient := hln_client.NewClient(
			hln_client.NewBuilder(),
			hln_client.NewRequester(pHLSClient),
		)
		if err := hlnClient.Redirect(pCtx, alias.GetAliasesList(friends), rawMsg); err != nil {
			pLogger.PushWarn(logBuilder.WithMessage("redirect"))
			_ = api.Response(pW, http.StatusBadGateway, "failed: redirect")
			return
		}

		pLogger.PushInfo(logBuilder.WithMessage("redirect"))
		_ = api.Response(pW, http.StatusOK, http_logger.CLogSuccess)
	}
}
