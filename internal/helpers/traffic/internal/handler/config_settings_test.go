package handler

import (
	"context"
	"fmt"
	"os"
	"testing"

	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
	testutils "github.com/number571/hidden-lake/test/utils"
)

func TestHandleConfigSettingsAPI(t *testing.T) {
	t.Parallel()

	addr := testutils.TgAddrs[22]
	os.RemoveAll(fmt.Sprintf(databaseTemplate, addr))

	srv, cancel, db, hltClient := testAllRun(addr)
	defer testAllFree(addr, srv, cancel, db)

	settings, err := hltClient.GetSettings(context.Background())
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

	if settings.GetMessagesCapacity() != tcCapacity {
		t.Error("incorrect messages capacity")
		return
	}

	if settings.GetMessageSizeBytes() != tcMessageSize {
		t.Error("incorrect messages size bytes")
		return
	}

	if settings.GetRandMessageSizeBytes() != hls_settings.CDefaultRandMessageSizeBytes {
		t.Error("incorrect rand message size bytes")
		return
	}
}
