package build

import (
	"os"
	"path/filepath"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/build"
)

func setNetworks(pInputPath, pFilename string) error {
	netPath := filepath.Join(pInputPath, pFilename)
	netVal, err := os.ReadFile(netPath) //nolint:gosec
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		return os.WriteFile(netPath, build.GNetworksVal, 0600)
	}
	networksYAML := &build.SNetworksYAML{}
	if err := encoding.DeserializeYAML(netVal, networksYAML); err != nil {
		return err
	}
	networksYAML.FNetworks[build.CDefaultNetwork] = networksYAML.FSettings
	return build.SetNetworks(networksYAML.FNetworks)
}
