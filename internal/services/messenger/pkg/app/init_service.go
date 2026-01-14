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

	timeoutMsg := "handle timeout"
	buildSettings := build.GetSettings()

	indexHandler := http.TimeoutHandler(handler.HandleIndexAPI(p.fHTTPLogger), buildSettings.GetHttpReadTimeout(), timeoutMsg)
	mux.Handle(hls_messenger_settings.CHandleIndexPath, indexHandler) // GET

	pushMessageHandler := http.TimeoutHandler(handler.HandlePushMessageAPI(pCtx, p.fHTTPLogger, p.fConfig, pHlkClient, p.fDatabase), buildSettings.GetHttpReadTimeout(), timeoutMsg)
	mux.Handle(hls_messenger_settings.CHandlePushMessagePath, pushMessageHandler) // POST

	loadMessagesHandler := http.TimeoutHandler(handler.HandleLoadMessagesAPI(pCtx, p.fHTTPLogger, p.fConfig, pHlkClient, p.fDatabase), buildSettings.GetHttpReadTimeout(), timeoutMsg)
	mux.Handle(hls_messenger_settings.CHandleLoadMessagesPath, loadMessagesHandler) // GET

	listenMessageHandler := http.TimeoutHandler(handler.HandleListenMessageAPI(pCtx, pMsgBroker), buildSettings.GetHttpReadTimeout()<<1, timeoutMsg)
	mux.Handle(hls_messenger_settings.CHandleListenMessagePath, listenMessageHandler) // GET

	p.fIntServiceHTTP = &http.Server{
		Addr:         p.fConfig.GetAddress().GetInternal(),
		Handler:      mux,
		ReadTimeout:  buildSettings.GetHttpReadTimeout() << 2,
		WriteTimeout: buildSettings.GetHttpHandleTimeout() << 2,
	}
}
