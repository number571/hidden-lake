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
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	settingsYAML := &build.SSettings{}
	if err := encoding.DeserializeYAML(settVal, settingsYAML); err != nil {
		panic(err)
	}
	return build.SetSettings(*settingsYAML)
}
