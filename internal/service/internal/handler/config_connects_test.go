package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/number571/go-peer/pkg/logger"
	net_message "github.com/number571/go-peer/pkg/network/message"
	hls_client "github.com/number571/hidden-lake/internal/service/pkg/client"
	"github.com/number571/hidden-lake/internal/service/pkg/config"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	"github.com/number571/hidden-lake/pkg/adapters/http/client"
	testutils "github.com/number571/hidden-lake/test/utils"
)

const (
	tcConnection = "tcp://localhost:1111"
)

func TestHandleConnectsAPI2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

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

	epClients := []client.IClient{
		client.NewClient(&tsRequester{}),
	}

	handler := HandleConfigConnectsAPI(ctx, httpLogger, epClients)
	if err := connectsAPIRequestOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := connectsAPIRequestPostOK(handler); err != nil {
		t.Error(err)
		return
	}
	if err := connectsAPIRequestDeleteOK(handler); err != nil {
		t.Error(err)
		return
	}

	if err := connectsAPIRequestMethod(handler); err == nil {
		t.Error("request success with invalid method")
		return
	}
	if err := connectsAPIRequestPostConnect(handler); err == nil {
		t.Error("request success with invalid connect")
		return
	}

	epClientsx := []client.IClient{
		client.NewClient(&tsRequester{fWithFail: true}),
	}

	handlerx := HandleConfigConnectsAPI(ctx, httpLogger, epClientsx)
	if err := connectsAPIRequestOK(handlerx); err == nil {
		t.Error("request success with invalid get connections")
		return
	}
	if err := connectsAPIRequestPostOK(handlerx); err == nil {
		t.Error("request success with invalid update editor (post)")
		return
	}
	if err := connectsAPIRequestDeleteOK(handlerx); err == nil {
		t.Error("request success with invalid update editor (delete)")
		return
	}
}

func connectsAPIRequestDeleteOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/", strings.NewReader("127.0.0.1:9999"))

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

func connectsAPIRequestPostOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("127.0.0.1:9999"))

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

func connectsAPIRequestOK(handler http.HandlerFunc) error {
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

func connectsAPIRequestPostConnect(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))

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

func connectsAPIRequestMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/", nil)

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

func TestHandleConnectsAPI(t *testing.T) {
	t.Parallel()

	pathCfg := fmt.Sprintf(tcPathConfigTemplate, 0)
	pathDB := fmt.Sprintf(tcPathDBTemplate, 0)

	_, node, _, cancel, srv := testAllCreate(pathCfg, pathDB, testutils.TgAddrs[6])
	defer testAllFree(node, cancel, srv, pathCfg, pathDB)

	client := hls_client.NewClient(
		hls_client.NewBuilder(),
		hls_client.NewRequester(
			testutils.TgAddrs[6],
			&http.Client{Timeout: time.Minute},
		),
	)

	testGetConnects(t, client)
	testAddConnect(t, client)
	testDelConnect(t, client)
}

func testGetConnects(t *testing.T, client hls_client.IClient) {
	connects, err := client.GetConnections(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if len(connects) != 1 || connects[0] != tcConnection {
		t.Error("len(connects) != 1 || connects[0] != tcConnection")
		return
	}
}

func testAddConnect(t *testing.T, client hls_client.IClient) {
	if err := client.AddConnection(context.Background(), "tcp://aaa"); err != nil {
		t.Error(err)
		return
	}
}

func testDelConnect(t *testing.T, client hls_client.IClient) {
	if err := client.DelConnection(context.Background(), "tcp://bbb"); err != nil {
		t.Error(err)
		return
	}
}

var (
	_ client.IRequester = &tsRequester{}
)

type tsRequester struct {
	fWithFail bool
}

func (p *tsRequester) GetIndex(context.Context) (string, error)                    { return "", nil }
func (p *tsRequester) GetSettings(context.Context) (config.IConfigSettings, error) { return nil, nil }
func (p *tsRequester) GetOnlines(context.Context) ([]string, error) {
	if p.fWithFail {
		return nil, errors.New("some error") // nolint: err113
	}
	return []string{tcConnection}, nil
}
func (p *tsRequester) DelOnline(context.Context, string) error {
	if p.fWithFail {
		return errors.New("some error") // nolint: err113
	}
	return nil
}
func (p *tsRequester) GetConnections(context.Context) ([]string, error) {
	if p.fWithFail {
		return nil, errors.New("some error") // nolint: err113
	}
	return []string{tcConnection}, nil
}
func (p *tsRequester) AddConnection(context.Context, string) error {
	if p.fWithFail {
		return errors.New("some error") // nolint: err113
	}
	return nil
}
func (p *tsRequester) DelConnection(context.Context, string) error {
	if p.fWithFail {
		return errors.New("some error") // nolint: err113
	}
	return nil
}
func (p *tsRequester) ProduceMessage(context.Context, net_message.IMessage) error { return nil }
