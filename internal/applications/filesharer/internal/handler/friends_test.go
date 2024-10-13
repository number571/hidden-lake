// nolint: goerr113
package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/internal/applications/filesharer/internal/config"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestFriendsPage(t *testing.T) {
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
	handler := FriendsPage(ctx, httpLogger, &config.SConfig{
		FSettings: &config.SConfigSettings{
			FLanguage: "ENG",
		},
	}, nil)

	if err := friendsRequest404(handler); err == nil {
		t.Error("request success with invalid path")
		return
	}
}

func friendsRequest404(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/friends/undefined", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code")
	}

	return nil
}
