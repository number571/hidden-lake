package handler

import (
	"context"
	"fmt"
	"os"
	"testing"

	testutils "github.com/number571/go-peer/test/utils"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

func TestHandleConfigSettingsAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[50]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, cancel, db, hltClient := testAllRun(addr)
	defer testAllFree(addr, srv, cancel, db)

	settings, err := hltClient.GetSettings(context.Background())
	if err != nil {
		t.Error(err)
		return
	}

	if settings.GetNetworkKey() != testutils.TCNetworkKey {
		t.Error("incorrect network key")
		return
	}

	if settings.GetKeySizeBits() != testutils.TcKeySize {
		t.Error("incorrect key size bits")
		return
	}

	if settings.GetWorkSizeBits() != testutils.TCWorkSize {
		t.Error("incorrect work size bits")
		return
	}

	if settings.GetMessagesCapacity() != testutils.TCCapacity {
		t.Error("incorrect messages capacity")
		return
	}

	if settings.GetMessageSizeBytes() != testutils.TCMessageSize {
		t.Error("incorrect messages size bytes")
		return
	}

	if settings.GetRandMessageSizeBytes() != hls_settings.CDefaultRandMessageSizeBytes {
		t.Error("incorrect rand message size bytes")
		return
	}
}
