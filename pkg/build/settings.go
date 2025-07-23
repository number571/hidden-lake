package build

import (
	"os"
	"path/filepath"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/build"
)

func setSettings(pInputPath, pFilename string) (bool, error) {
	settPath := filepath.Join(pInputPath, pFilename)
	settVal, err := os.ReadFile(settPath) //nolint:gosec
	if err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
		return false, nil
	}
	settingsYAML := &build.SSettings{}
	if err := encoding.DeserializeYAML(settVal, settingsYAML); err != nil {
		return false, err
	}
	return true, build.SetSettings(*settingsYAML)
}
