// nolint: err113
package build

import (
	_ "embed"

	"github.com/number571/go-peer/pkg/encoding"
)

var (
	//go:embed version.yml
	gVersion []byte
	GVersion string
)

func init() {
	var versionYAML struct {
		FVersion string `yaml:"version"`
	}
	if err := encoding.DeserializeYAML(gVersion, &versionYAML); err != nil {
		panic(err)
	}
	GVersion = versionYAML.FVersion
}
