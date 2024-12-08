package handler

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/http/httptest"
// 	"strings"
// 	"testing"
// 	"time"

// 	"github.com/number571/go-peer/pkg/logger"
// 	"github.com/number571/hidden-lake/internal/service/pkg/app/config"
// 	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
// 	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
// 	testutils "github.com/number571/hidden-lake/test/utils"
// )

// func TestHandleConnectsAPI2(t *testing.T) {
// 	t.Parallel()

// 	ctx := context.Background()

// 	httpLogger := std_logger.NewStdLogger(
// 		func() std_logger.ILogging {
// 			logging, err := std_logger.LoadLogging([]string{})
// 			if err != nil {
// 				panic(err)
// 			}
// 			return logging
// 		}(),
// 		func(_ logger.ILogArg) string {
// 			return ""
// 		},
// 	)

// 	handler := HandleConfigConnectsAPI(ctx, newTsWrapper(true), httpLogger, newTsNode(true, true, true, true))
// 	if err := connectsAPIRequestOK(handler); err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if err := connectsAPIRequestPostOK(handler); err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if err := connectsAPIRequestDeleteOK(handler); err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	if err := connectsAPIRequestMethod(handler); err == nil {
// 		t.Error("request success with invalid method")
// 		return
// 	}
// 	if err := connectsAPIRequestPostConnect(handler); err == nil {
// 		t.Error("request success with invalid connect")
// 		return
// 	}

// 	handlerx := HandleConfigConnectsAPI(ctx, newTsWrapper(false), httpLogger, newTsNode(true, true, true, true))
// 	if err := connectsAPIRequestPostOK(handlerx); err == nil {
// 		t.Error("request success with invalid update editor (post)")
// 		return
// 	}
// 	if err := connectsAPIRequestDeleteOK(handlerx); err == nil {
// 		t.Error("request success with invalid update editor (delete)")
// 		return
// 	}
// }

// func connectsAPIRequestDeleteOK(handler http.HandlerFunc) error {
// 	w := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodDelete, "/", strings.NewReader("127.0.0.1:9999"))

// 	handler(w, req)
// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return errors.New("bad status code") // nolint: err113
// 	}

// 	if _, err := io.ReadAll(res.Body); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func connectsAPIRequestPostOK(handler http.HandlerFunc) error {
// 	w := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("127.0.0.1:9999"))

// 	handler(w, req)
// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return errors.New("bad status code") // nolint: err113
// 	}

// 	if _, err := io.ReadAll(res.Body); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func connectsAPIRequestOK(handler http.HandlerFunc) error {
// 	w := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodGet, "/", nil)

// 	handler(w, req)
// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return errors.New("bad status code") // nolint: err113
// 	}

// 	if _, err := io.ReadAll(res.Body); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func connectsAPIRequestPostConnect(handler http.HandlerFunc) error {
// 	w := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))

// 	handler(w, req)
// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return errors.New("bad status code") // nolint: err113
// 	}

// 	if _, err := io.ReadAll(res.Body); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func connectsAPIRequestMethod(handler http.HandlerFunc) error {
// 	w := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodPut, "/", nil)

// 	handler(w, req)
// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		return errors.New("bad status code") // nolint: err113
// 	}

// 	if _, err := io.ReadAll(res.Body); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func TestHandleConnectsAPI(t *testing.T) {
// 	t.Parallel()

// 	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 0)
// 	pathDB := fmt.Sprintf(tcPathDBTemplate, 0)

// 	wcfg, node, _, cancel, srv := testAllCreate(pathCfg, pathDB, testutils.TgAddrs[6])
// 	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

// 	client := hls_client.NewClient(
// 		hls_client.NewBuilder(),
// 		hls_client.NewRequester(
// 			"http://"+testutils.TgAddrs[6],
// 			&http.Client{Timeout: time.Minute},
// 		),
// 	)

// 	connect := "test_connect4"
// 	testGetConnects(t, client, wcfg.GetConfig())
// 	testAddConnect(t, client, connect)
// 	testDelConnect(t, client, connect)
// }

// func testGetConnects(t *testing.T, client hls_client.IClient, cfg config.IConfig) {
// 	connects, err := client.GetConnections(context.Background())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	if len(connects) != 3 {
// 		t.Error("length of connects != 3")
// 		return
// 	}

// 	for i := range connects {
// 		if connects[i] != cfg.GetConnections()[i] {
// 			t.Error("connections from config not equals with get")
// 			return
// 		}
// 	}
// }

// func testAddConnect(t *testing.T, client hls_client.IClient, connect string) {
// 	err := client.AddConnection(context.Background(), connect)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	connects, err := client.GetConnections(context.Background())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	for _, conn := range connects {
// 		if conn == connect {
// 			return
// 		}
// 	}
// 	t.Errorf("undefined connection key by '%s'", connect)
// }

// func testDelConnect(t *testing.T, client hls_client.IClient, connect string) {
// 	err := client.DelConnection(context.Background(), connect)
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	connects, err := client.GetConnections(context.Background())
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}

// 	for _, conn := range connects {
// 		if conn == connect {
// 			t.Errorf("deleted connect exists for '%s'", connect)
// 			return
// 		}
// 	}
// }
