package app

import (
	"context"
	"net/http"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/notifier/internal/handler"
	"github.com/number571/hidden-lake/internal/services/notifier/internal/handler/incoming"
	"github.com/number571/hidden-lake/internal/services/notifier/internal/handler/process"
	hls_notifier_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/broker"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
)

func (p *sApp) initExternalServiceHTTP(
	pCtx context.Context,
	pMsgBroker broker.IDataBroker,
) {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_notifier_settings.CNotifyPath,
		incoming.HandleIncomingNotifyHTTP(pCtx, p.fHTTPLogger, pMsgBroker),
	) // POST

	mux.HandleFunc(
		hls_notifier_settings.CMsgboxPath,
		incoming.HandleIncomingMsgboxHTTP(pCtx, p.fHTTPLogger, p.fConfig, p.fDatabase, pMsgBroker),
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
	pMsgBroker broker.IDataBroker,
	pProcessManager process.INotifyProcessManager,
) {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_notifier_settings.CHandleIndexPath,
		handler.HandleIndexAPI(p.fHTTPLogger),
	) // GET

	mux.HandleFunc(
		hls_notifier_settings.CHandleChanNotifyPath,
		handler.HandleChanNotifyAPI(pCtx, p.fHTTPLogger, pHlkClient),
	) // POST

	mux.HandleFunc(
		hls_notifier_settings.CHandleChanMsgboxPushPath,
		handler.HandleChanPushAPI(pCtx, p.fHTTPLogger, p.fConfig, pHlkClient, p.fDatabase, pProcessManager),
	) // POST

	mux.HandleFunc(
		hls_notifier_settings.CHandleChanMsgboxLoadPath,
		handler.HandleChanLoadAPI(pCtx, p.fHTTPLogger, p.fConfig, pHlkClient, p.fDatabase),
	) // GET

	mux.HandleFunc(
		hls_notifier_settings.CHandleChanMsgboxSizePath,
		handler.HandleChanSizeAPI(pCtx, p.fHTTPLogger, p.fConfig, p.fDatabase),
	) // GET

	mux.HandleFunc(
		hls_notifier_settings.CHandleChanListenPath,
		handler.HandleChanListenAPI(pCtx, p.fHTTPLogger, pMsgBroker),
	) // GET

	buildSettings := build.GetSettings()
	p.fIntServiceHTTP = &http.Server{ // nolint: gosec
		Addr:         p.fConfig.GetAddress().GetInternal(),
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpCallbackTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpCallbackTimeout(),
		WriteTimeout: buildSettings.GetHttpCallbackTimeout(),
	}
}
