// nolint: err113
package build

import (
	_ "embed"
	"errors"
	"net/url"
	"sync"

	"github.com/number571/go-peer/pkg/encoding"
)

const (
	CDefaultNetwork = "__default_network__"
)

var (
	//go:embed networks.yml
	gNetworksVal []byte
	gNetworksMtx sync.RWMutex
	gNetworks    = make(map[string]SNetwork, 16)
)

func init() {
	networksYAML := &SNetworksYAML{}
	if err := encoding.DeserializeYAML(gNetworksVal, networksYAML); err != nil {
		panic(err)
	}
	networksYAML.FNetworks[CDefaultNetwork] = networksYAML.FDefault
	if err := SetNetworks(networksYAML.FNetworks); err != nil {
		panic(err) // build network should be always correct
	}
}

type SNetworksYAML struct {
	FDefault  SNetwork            `yaml:"default"`
	FNetworks map[string]SNetwork `yaml:"networks"`
}

type SNetwork struct {
	FMessageSizeBytes uint64       `yaml:"message_size_bytes"`
	FWorkSizeBits     uint64       `yaml:"work_size_bits"`
	FConnections      SConnections `yaml:"connections"`
}

type SConnections []string

func (p SConnections) GetByScheme(pScheme string) []string {
	connects := make([]string, 0, len(p))
	for _, c := range p {
		u, err := url.Parse(c)
		if err != nil || u.Scheme != pScheme {
			continue
		}
		connects = append(connects, u.Host)
	}
	return connects
}

func GetNetwork(k string) (SNetwork, bool) {
	gNetworksMtx.RLock()
	v, ok := gNetworks[k]
	gNetworksMtx.RUnlock()
	return v, ok
}

func GetNetworks() map[string]SNetwork {
	gNetworksMtx.RLock()
	defer gNetworksMtx.RUnlock()

	r := make(map[string]SNetwork, len(gNetworks))
	for k, v := range gNetworks {
		r[k] = v
	}
	return r
}

func SetNetworks(networksMap map[string]SNetwork) error {
	if _, ok := networksMap[CDefaultNetwork]; !ok {
		return errors.New("default network not found in map")
	}
	for _, v := range networksMap {
		if err := v.validate(); err != nil {
			return err
		}
	}
	gNetworksMtx.Lock()
	gNetworks = networksMap
	gNetworksMtx.Unlock()
	return nil
}

func (p SNetwork) validate() error {
	switch { // nolint: gocritic, staticcheck
	case p.FMessageSizeBytes == 0:
		return errors.New("message_size_bytes = 0")
	}
	for _, c := range p.FConnections {
		u, err := url.Parse(c)
		if err != nil || u.Scheme == "" {
			return errors.New("parse url conn")
		}
	}
	return nil
}
