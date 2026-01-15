package app

import (
	"context"
	"net/http"

	"github.com/number571/hidden-lake/build"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/handler"
	"github.com/number571/hidden-lake/internal/services/messenger/internal/handler/incoming"
	"github.com/number571/hidden-lake/internal/services/messenger/pkg/message"
	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"
)

func (p *sApp) initExternalServiceHTTP(
	pCtx context.Context,
	pHlkClient hlk_client.IClient,
	pMsgBroker message.IMessageBroker,
) {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_messenger_settings.CPushPath,
		incoming.HandleIncomingPushHTTP(pCtx, p.fHTTPLogger, p.fDatabase, pMsgBroker, pHlkClient),
	) // POST

	buildSettings := build.GetSettings()
	p.fExtServiceHTTP = &http.Server{
		Addr:         p.fConfig.GetAddress().GetExternal(),
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpHandleTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpReadTimeout(),
		WriteTimeout: buildSettings.GetHttpHandleTimeout(),
	}
}

func (p *sApp) initInternalServiceHTTP(
	pCtx context.Context,
	pHlkClient hlk_client.IClient,
	pMsgBroker message.IMessageBroker,
) {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_messenger_settings.CHandleIndexPath,
		handler.HandleIndexAPI(p.fHTTPLogger),
	) // GET

	mux.HandleFunc(
		hls_messenger_settings.CHandlePushMessagePath,
		handler.HandlePushMessageAPI(pCtx, p.fHTTPLogger, p.fConfig, pHlkClient, p.fDatabase),
	) // POST

	mux.HandleFunc(
		hls_messenger_settings.CHandleLoadMessagesPath,
		handler.HandleLoadMessagesAPI(pCtx, p.fHTTPLogger, p.fConfig, pHlkClient, p.fDatabase),
	) // GET

	mux.HandleFunc(
		hls_messenger_settings.CHandleListenChatPath,
		handler.HandleListenChatAPI(pCtx, pMsgBroker),
	) // GET

	buildSettings := build.GetSettings()
	p.fIntServiceHTTP = &http.Server{ // nolint: gosec
		Addr:         p.fConfig.GetAddress().GetInternal(),
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpCallbackTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpCallbackTimeout(),
		WriteTimeout: buildSettings.GetHttpCallbackTimeout(),
	}
}
