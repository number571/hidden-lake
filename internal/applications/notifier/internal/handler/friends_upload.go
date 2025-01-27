package handler

import (
	"context"
	"net/http"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/notifier/internal/utils"
	"github.com/number571/hidden-lake/internal/applications/notifier/pkg/app/config"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/webui"

	hlm_settings "github.com/number571/hidden-lake/internal/applications/notifier/pkg/settings"
)

type sUploadFile struct {
	*sTemplate
	FMessageLimit uint64
}

func FriendsUploadPage(
	pCtx context.Context,
	pLogger logger.ILogger,
	pCfg config.IConfig,
	pHlsClient hls_client.IClient,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		logBuilder := http_logger.NewLogBuilder(hlm_settings.GServiceName.Short(), pR)

		if pR.URL.Path != "/friends/upload" {
			NotFoundPage(pLogger, pCfg)(pW, pR)
			return
		}

		msgLimit, err := utils.GetMessageLimit(pCtx, pHlsClient)
		if err != nil {
			ErrorPage(pLogger, pCfg, "get_message_size", "get message size (limit)")(pW, pR)
			return
		}

		res := &sUploadFile{
			sTemplate:     getTemplate(pCfg),
			FMessageLimit: msgLimit,
		}

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = webui.MustParseTemplate("index.html", "notifier/upload.html").Execute(pW, res)
	}
}
