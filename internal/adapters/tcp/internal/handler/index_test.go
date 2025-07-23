package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	"github.com/number571/hidden-lake/pkg/adapters/http/client"
	testutils "github.com/number571/hidden-lake/test/utils"
)

// const (
// 	tcPathConfigTemplate = "config_test_%d.yml"
// )

func TestErrorsAPI(t *testing.T) {
	t.Parallel()

	client := client.NewClient(
		client.NewRequester(
			testutils.TcUnknownHost,
			&http.Client{Timeout: time.Second},
		),
	)

	if err := client.AddConnection(context.Background(), ""); err == nil {
		t.Fatal("success add connection with unknown host")
	}

	if err := client.DelConnection(context.Background(), ""); err == nil {
		t.Fatal("success del connection with unknown host")
	}

	if _, err := client.GetIndex(context.Background()); err == nil {
		t.Fatal("success get index with unknown host")
	}

	if _, err := client.GetConnections(context.Background()); err == nil {
		t.Fatal("success get connections with unknown host")
	}

	if _, err := client.GetOnlines(context.Background()); err == nil {
		t.Fatal("success get onlines with unknown host")
	}

	if err := client.DelOnline(context.Background(), "test"); err == nil {
		t.Fatal("success del online key with unknown host")
	}
}

func TestHandleIndexAPI2(t *testing.T) {
	t.Parallel()

	log := logger.NewLogger(
		logger.NewSettings(&logger.SSettings{}),
		func(_ logger.ILogArg) string { return "" },
	)

	handler := HandleIndexAPI(log)
	if err := indexAPIRequestOK(handler); err != nil {
		t.Fatal(err)
	}
}

func indexAPIRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusOK {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

// func TestHandleIndexAPI(t *testing.T) {
// 	t.Parallel()

// 	addr := testutils.TgAddrs[AAA]
// 	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 3)

// 	client := client.NewClient(
// 		client.NewRequester(
// 			addr,
// 			&http.Client{Timeout: time.Minute},
// 		),
// 	)

// 	title, err := client.GetIndex(context.Background())
// 	if err != nil {
// 		t.Fatal(err)
// 		return
// 	}

// 	if title != settings.CAppFullName {
// 		t.Fatal("incorrect title pattern")
// 		return
// 	}
// }
