package handler

import (
	"context"
	"net/http"
	"testing"
	"time"

	hll_client "github.com/number571/hidden-lake/internal/helpers/loader/pkg/client"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandleConfigSettingsAPI(t *testing.T) {
	t.Parallel()

	service := testRunService(testutils.TgAddrs[28])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hllClient := hll_client.NewClient(
		hll_client.NewRequester(
			"http://"+testutils.TgAddrs[28],
			&http.Client{Timeout: time.Second / 2},
		),
	)

	settings, err := hllClient.GetSettings(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if settings.GetNetworkKey() != tcNetworkKey {
		t.Error("incorrect network key")
		return
	}

	if settings.GetMessagesCapacity() != tcCapacity {
		t.Error("incorrect messages capacity")
		return
	}

	if settings.GetWorkSizeBits() != tcWorkSize {
		t.Error("incorrect work size bits")
		return
	}
}
