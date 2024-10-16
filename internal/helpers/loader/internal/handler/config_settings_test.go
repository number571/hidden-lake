package handler

import (
	"context"
	"net/http"
	"testing"
	"time"

	testutils "github.com/number571/go-peer/test/utils"
	hll_client "github.com/number571/hidden-lake/internal/helpers/loader/pkg/client"
)

func TestHandleConfigSettingsAPI(t *testing.T) {
	t.Parallel()

	service := testRunService(testutils.TgAddrs[52])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hllClient := hll_client.NewClient(
		hll_client.NewRequester(
			"http://"+testutils.TgAddrs[52],
			&http.Client{Timeout: time.Second / 2},
		),
	)

	settings, err := hllClient.GetSettings(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if settings.GetNetworkKey() != testutils.TCNetworkKey {
		t.Error("incorrect network key")
		return
	}

	if settings.GetMessagesCapacity() != testutils.TCCapacity {
		t.Error("incorrect messages capacity")
		return
	}

	if settings.GetWorkSizeBits() != testutils.TCWorkSize {
		t.Error("incorrect work size bits")
		return
	}
}
