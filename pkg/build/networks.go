package build

import (
	"os"
	"path/filepath"

	"github.com/number571/go-peer/pkg/encoding"
	"github.com/number571/hidden-lake/build"
)

func setNetworks(pInputPath, pFilename string) (bool, error) {
	netPath := filepath.Join(pInputPath, pFilename)
	netVal, err := os.ReadFile(netPath) //nolint:gosec
	if err != nil {
		if !os.IsNotExist(err) {
			return false, err
		}
		return false, nil
	}
	networksYAML := &build.SNetworksYAML{}
	if err := encoding.DeserializeYAML(netVal, networksYAML); err != nil {
		return false, err
	}
	networksYAML.FNetworks[build.CDefaultNetwork] = networksYAML.FDefault
	return true, build.SetNetworks(networksYAML.FNetworks)
}
