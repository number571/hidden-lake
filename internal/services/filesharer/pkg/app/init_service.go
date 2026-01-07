package app

import (
	"context"
	"net/http"

	"github.com/number571/hidden-lake/build"
	hlk_client "github.com/number571/hidden-lake/internal/kernel/pkg/client"
	"github.com/number571/hidden-lake/internal/services/filesharer/internal/handler"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"
)

func (p *sApp) initExternalServiceHTTP(pCtx context.Context, pHlkClient hlk_client.IClient) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		hls_filesharer_settings.CLoadPath,
		handler.HandleIncomingLoadHTTP(pCtx, p.fHTTPLogger, p.fPathTo, pHlkClient),
	) // GET

	mux.HandleFunc(
		hls_filesharer_settings.CListPath,
		handler.HandleIncomingListHTTP(p.fHTTPLogger, p.fConfig, p.fPathTo),
	) // GET

	mux.HandleFunc(
		hls_filesharer_settings.CInfoPath,
		handler.HandleIncomingInfoHTTP(p.fHTTPLogger, p.fPathTo),
	) // GET

	buildSettings := build.GetSettings()
	p.fExtServiceHTTP = &http.Server{
		Addr:         p.fConfig.GetAddress().GetExternal(),
		Handler:      http.TimeoutHandler(mux, buildSettings.GetHttpHandleTimeout(), "handle timeout"),
		ReadTimeout:  buildSettings.GetHttpReadTimeout(),
		WriteTimeout: buildSettings.GetHttpHandleTimeout(),
	}
}
