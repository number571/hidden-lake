package handler

import (
	"context"
	"net/http"
	"testing"
	"time"

	hle_client "github.com/number571/hidden-lake/internal/helpers/encryptor/pkg/client"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandleConfigSettingsAPI(t *testing.T) {
	t.Parallel()

	service := testRunService(testutils.TgAddrs[34])
	defer service.Close()

	time.Sleep(100 * time.Millisecond)
	hleClient := hle_client.NewClient(
		hle_client.NewRequester(
			"http://"+testutils.TgAddrs[34],
			&http.Client{Timeout: time.Second / 2},
			testNetworkMessageSettings(),
		),
	)

	settings, err := hleClient.GetSettings(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if settings.GetNetworkKey() != tcNetworkKey {
		t.Error("incorrect network key")
		return
	}

	if settings.GetWorkSizeBits() != tcWorkSize {
		t.Error("incorrect work size bits")
		return
	}

	if settings.GetMessageSizeBytes() != tcMessageSize {
		t.Error("incorrect message size bytes")
		return
	}
}
