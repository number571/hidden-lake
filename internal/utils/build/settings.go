package build

import (
	"os"
	"path/filepath"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/build"
)

func setSettings(pInputPath, pFilename string) error {
	settPath := filepath.Join(pInputPath, pFilename)
	settVal, err := os.ReadFile(settPath) //nolint:gosec
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return os.WriteFile(settPath, build.GSettingsVal, 0600)
	}
	settingsYAML := &build.SSettings{}
	if err := encoding.DeserializeYAML(settVal, settingsYAML); err != nil {
		return err
	}
	return build.SetSettings(*settingsYAML)
}
