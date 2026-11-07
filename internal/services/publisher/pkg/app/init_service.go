package app

import (
	"context"
	"net/http"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/services/publisher/internal/handler"
	"github.com/number571/hidden-lake/internal/services/publisher/internal/handler/incoming"
	"github.com/number571/hidden-lake/internal/services/publisher/internal/handler/process"
	hls_publisher_settings "github.com/number571/hidden-lake/internal/services/publisher/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/broker"
	hlk_client "github.com/number571/hidden-lake/pkg/api/kernel/client"
)

func (p *sApp) initExternalServiceHTTP(
	pCtx context.Context,
	pMsgBroker broker.IDataBroker,
) {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_publisher_settings.CPushPath,
		incoming.HandleIncomingPushHTTP(pCtx, p.fHTTPLogger, p.fConfig, p.fDatabase, pMsgBroker),
	) // POST

	mux.HandleFunc(
		hls_publisher_settings.CLoadPath,
		incoming.HandleIncomingLoadHTTP(pCtx, p.fHTTPLogger, p.fConfig, p.fDatabase, pMsgBroker),
	) // GET

	mux.HandleFunc(
		hls_publisher_settings.CSizePath,
		incoming.HandleIncomingSizeHTTP(pCtx, p.fHTTPLogger, p.fConfig, p.fDatabase, pMsgBroker),
	) // GET

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
	pProcessManager process.IPublishProcessManager,
) {
	mux := http.NewServeMux()

	mux.HandleFunc(
		hls_publisher_settings.CHandleIndexPath,
		handler.HandleIndexAPI(p.fHTTPLogger),
	) // GET

	mux.HandleFunc(
		hls_publisher_settings.CHandleMsgboxSizePath,
		handler.HandleMsgboxSizeAPI(pCtx, p.fHTTPLogger, p.fConfig, pHlkClient, p.fDatabase),
	) // GET

	mux.HandleFunc(
		hls_publisher_settings.CHandleMsgboxLoadPath,
		handler.HandleMsgboxLoadAPI(pCtx, p.fHTTPLogger, p.fConfig, pHlkClient, p.fDatabase),
	) // GET

	mux.HandleFunc(
		hls_publisher_settings.CHandleMsgboxPushPath,
		handler.HandleMsgboxPushAPI(pCtx, p.fHTTPLogger, p.fConfig, pHlkClient, p.fDatabase, pProcessManager),
	) // POST

	mux.HandleFunc(
		hls_publisher_settings.CHandleMsgboxListenPath,
		handler.HandleMsgboxListenAPI(pCtx, p.fHTTPLogger, pMsgBroker),
	) // GET

	buildSettings := build.GetSettings()
	p.fIntServiceHTTP = &http.Server{ // nolint: gosec
		Addr:         p.fConfig.GetAddress().GetInternal(),
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpCallbackTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpCallbackTimeout(),
		WriteTimeout: buildSettings.GetHttpCallbackTimeout(),
	}
}
