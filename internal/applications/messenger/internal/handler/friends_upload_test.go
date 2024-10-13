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
	std_logger "github.com/number571/hidden-lake/internal/modules/logger/std"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
)

func TestFriendsUploadPage(t *testing.T) {
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
	cfg := &config.SConfig{
		FSettings: &config.SConfigSettings{
			FLanguage: "ENG",
		},
	}

	hlsClient := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+cfg.GetConnection(),
			&http.Client{Timeout: (10 * time.Minute)},
		),
	)

	handler := FriendsUploadPage(ctx, httpLogger, cfg, hlsClient)

	if err := friendsUploadRequest404(handler); err == nil {
		t.Error("request success with invalid path")
		return
	}
}

func friendsUploadRequest404(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/friends/upload/undefined", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}
