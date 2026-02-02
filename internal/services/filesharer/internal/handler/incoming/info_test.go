package incoming

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/number571/go-peer/pkg/logger"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
)

func TestHandleIncomingInfoHTTP(t *testing.T) {
	t.Parallel()

	logging, err := std_logger.LoadLogging([]string{})
	if err != nil {
		t.Fatal(err)
	}

	httpLogger := std_logger.NewStdLogger(
		logging,
		func(_ logger.ILogArg) string {
			return ""
		},
	)

	ctx := context.Background()
	handler := HandleIncomingInfoHTTP(ctx, httpLogger, "./testdata", newTsHLSClient(true))
	if err := incomingInfoRequestOK(handler); err != nil {
		t.Fatal(err)
	}
	if err := incomingInfoRequestInvalidPersonal(handler); err != nil {
		t.Fatal(err)
	}
	if err := incomingInfoRequestInvalidFileName(handler); err != nil {
		t.Fatal(err)
	}
	if err := incomingInfoRequestInvalidMethod(handler); err != nil {
		t.Fatal(err)
	}
	if err := incomingInfoRequestInvalidFriend(handler); err != nil {
		t.Fatal(err)
	}
	if err := incomingInfoRequestNotFoundFile(handler); err != nil {
		t.Fatal(err)
	}

	handlerX := HandleIncomingInfoHTTP(ctx, httpLogger, "./testdata", newTsHLSClient(true))
	if err := incomingInfoRequestOK(handlerX); err != nil {
		t.Fatal(err)
	}
}

func incomingInfoRequestOK(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?name=file.txt&personal=false", nil)

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

func incomingInfoRequestInvalidFriend(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	// aliasName := pR.Header.Get(hlk_settings.CHeaderSenderName) (void)
	req := httptest.NewRequest(http.MethodGet, "/?name=file.txt&personal=true", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusForbidden {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func incomingInfoRequestInvalidPersonal(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?name=file.txt&personal=qwefew", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusBadRequest {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func incomingInfoRequestInvalidFileName(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?name=&personal", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusBadRequest {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func incomingInfoRequestInvalidMethod(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusMethodNotAllowed {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}

func incomingInfoRequestNotFoundFile(handler http.HandlerFunc) error {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/?name=file111.txt&personal=false", nil)

	handler(w, req)
	res := w.Result()
	defer func() { _ = res.Body.Close() }()

	if res.StatusCode != http.StatusNotFound {
		return errors.New("bad status code") // nolint: err113
	}

	if _, err := io.ReadAll(res.Body); err != nil {
		return err
	}

	return nil
}
