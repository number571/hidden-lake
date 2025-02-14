package handler

import (
	"net/http"
	"slices"
	"strings"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/notifier/pkg/app/config"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/notifier/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/api"
	http_logger "github.com/number571/hidden-lake/internal/utils/logger/http"
	"github.com/number571/hidden-lake/internal/webui"
)

type sChannels struct {
	*sTemplate
	FChannels        []string
	FChannelsBaseURL string
}

func ChannelsPage(
	pLogger logger.ILogger,
	pWrapper config.IWrapper,
) http.HandlerFunc {
	return func(pW http.ResponseWriter, pR *http.Request) {
		cfg := pWrapper.GetConfig()
		logBuilder := http_logger.NewLogBuilder(hlm_settings.GServiceName.Short(), pR)

		if pR.URL.Path != "/channels" {
			NotFoundPage(pLogger, cfg)(pW, pR)
			return
		}

		if err := pR.ParseForm(); err != nil {
			ErrorPage(pLogger, cfg, "parse_form", "parse form")(pW, pR)
			return
		}

		switch pR.FormValue("method") {
		case http.MethodPost:
			key := strings.TrimSpace(pR.FormValue("key")) // may be nil
			if key == "" {
				ErrorPage(pLogger, cfg, "key_nil", "key is nil")(pW, pR)
				return
			}

			channels := uniqAppendToSlice(cfg.GetChannels(), key)
			if err := pWrapper.GetEditor().UpdateChannels(channels); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("update_channels"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: update channels")
				return
			}
		case http.MethodDelete:
			key := strings.TrimSpace(pR.FormValue("key"))
			if key == "" {
				ErrorPage(pLogger, cfg, "get_key", "key is nil")(pW, pR)
				return
			}
			channels := slices.DeleteFunc(
				cfg.GetChannels(),
				func(p string) bool { return p == key },
			)
			if err := pWrapper.GetEditor().UpdateChannels(channels); err != nil {
				pLogger.PushWarn(logBuilder.WithMessage("delete_channels"))
				_ = api.Response(pW, http.StatusInternalServerError, "failed: delete channels")
				return
			}
		}

		result := new(sChannels)
		result.sTemplate = getTemplate(cfg)
		result.FChannels = cfg.GetChannels()
		result.FChannelsBaseURL = "/channels/chat"

		pLogger.PushInfo(logBuilder.WithMessage(http_logger.CLogSuccess))
		_ = webui.MustParseTemplate("index.html", "notifier/channels.html").Execute(pW, result)
	}
}

func uniqAppendToSlice(pSlice []string, pStr string) []string {
	if slices.Contains(pSlice, pStr) {
		return pSlice
	}
	return append(pSlice, pStr)
}
