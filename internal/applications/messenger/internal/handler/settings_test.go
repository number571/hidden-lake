// nolint: goerr113
package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/messenger/internal/config"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestSettingsPage(t *testing.T) {
	t.Parallel()

	logging, err := std_logger.LoadLogging([]string{})
	if err != nil {
		t.Error(err)
		return
	}

	httpLogger := std_logger.NewStdLogger(
		logging,
		func(_ logger.ILogArg) string {
			return ""
		},
	)

	ctx := context.Background()
	cfgWrapper := config.NewWrapper(&config.SConfig{
		FSettings: &config.SConfigSettings{
			FLanguage: "ENG",
		},
	})

	hlsClient := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+cfgWrapper.GetConfig().GetConnection(),
			&http.Client{Timeout: (10 * time.Minute)},
		),
	)

	handler := SettingsPage(ctx, httpLogger, cfgWrapper, hlsClient)

	if err := settingsRequest404(handler); err == nil {
		t.Error("request success with invalid path")
		return
	}
}

func settingsRequest404(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/settings/undefined", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}
