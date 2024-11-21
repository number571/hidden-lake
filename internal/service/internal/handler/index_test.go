package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	pkg_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	"github.com/number571/hidden-lake/pkg/request"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestErrorsAPI(t *testing.T) {
	t.Parallel()

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+testutils.TcUnknownHost,
			&http.Client{Timeout: time.Second},
		),
	)

	if err := client.AddConnection(context.Background(), ""); err == nil {
		t.Error("success add connection with unknown host")
		return
	}

	if err := client.DelConnection(context.Background(), ""); err == nil {
		t.Error("success del connection with unknown host")
		return
	}

	if err := client.AddFriend(context.Background(), "", tgPrivKey1.GetPubKey()); err == nil {
		t.Error("success add friend with unknown host")
		return
	}

	if err := client.DelFriend(context.Background(), ""); err == nil {
		t.Error("success del friend with unknown host")
		return
	}

	if err := client.SendRequest(context.Background(), "", request.NewRequestBuilder().Build()); err == nil {
		t.Error("success send request with unknown host")
		return
	}

	if _, err := client.FetchRequest(context.Background(), "", request.NewRequestBuilder().Build()); err == nil {
		t.Error("success fetch request with unknown host")
		return
	}

	if _, err := client.GetIndex(context.Background()); err == nil {
		t.Error("success get index with unknown host")
		return
	}

	if _, err := client.GetConnections(context.Background()); err == nil {
		t.Error("success get connections with unknown host")
		return
	}

	if _, err := client.GetFriends(context.Background()); err == nil {
		t.Error("success get friends with unknown host")
		return
	}

	if _, err := client.GetOnlines(context.Background()); err == nil {
		t.Error("success get onlines with unknown host")
		return
	}

	if _, err := client.GetPubKey(context.Background()); err == nil {
		t.Error("success get pub key with unknown host")
		return
	}

	if _, err := client.GetSettings(context.Background()); err == nil {
		t.Error("success get settings with unknown host")
		return
	}

	if err := client.DelOnline(context.Background(), "test"); err == nil {
		t.Error("success del online key with unknown host")
		return
	}
}

func TestHandleIndexAPI2(t *testing.T) {
	t.Parallel()

	httpLogger := std_logger.NewStdLogger(
		func() std_logger.ILogging {
			logging, err := std_logger.LoadLogging([]string{})
			if err != nil {
				panic(err)
			}
			return logging
		}(),
		func(_ logger.ILogArg) string {
			return ""
		},
	)

	handler := HandleIndexAPI(httpLogger)
	if err := indexAPIRequestOK(handler); err != nil {
		t.Error(err)
		return
	}
}

func indexAPIRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func TestHandleIndexAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[15]
	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 3)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 3)

	_, node, _, cancel, srv := testAllCreate(pathCfg, pathDB, addr)
	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			"http://"+addr,
			&http.Client{Timeout: time.Minute},
		),
	)

	title, err := client.GetIndex(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if title != pkg_settings.CServiceFullName {
		t.Error("incorrect title pattern")
		return
	}
}
