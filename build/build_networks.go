// nolint: err113
package build

import (
	_ "embed"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/number571/go-peer/pkg/encoding"
)

const (
	CDefaultNetwork = "__default_network__"
)

var (
	//go:embed networks.yml
	gNetworks []byte
	GNetworks map[string]SNetwork
)

func init() {
	networksYAML := &SNetworksYAML{}
	if err := encoding.DeserializeYAML(gNetworks, networksYAML); err != nil {
		panic(err)
	}
	if _, ok := networksYAML.FNetworks[CDefaultNetwork]; ok {
		panic(fmt.Sprintf("network '%s' already exist", CDefaultNetwork))
	}
	GNetworks = networksYAML.FNetworks
	GNetworks[CDefaultNetwork] = networksYAML.FSettings
	for _, n := range GNetworks {
		if err := n.validate(); err != nil {
			panic(err)
		}
	}
}

type SNetworksYAML struct {
	FSettings SNetwork            `yaml:"settings"`
	FNetworks map[string]SNetwork `yaml:"networks"`
}

type SNetwork struct {
	FMessageSizeBytes uint64   `yaml:"message_size_bytes"`
	FFetchTimeoutMS   uint64   `yaml:"fetch_timeout_ms"`
	FQueuePeriodMS    uint64   `yaml:"queue_period_ms"`
	FWorkSizeBits     uint64   `yaml:"work_size_bits"`
	FConnections      []string `yaml:"connections"`
}

func (p SNetwork) validate() error {
	switch {
	case p.FMessageSizeBytes == 0:
		return errors.New("message_size_bytes = 0")
	case p.FFetchTimeoutMS == 0:
		return errors.New("fetch_timeout_ms = 0")
	case p.FQueuePeriodMS == 0:
		return errors.New("queue_period_ms = 0")
	}
	for _, c := range p.FConnections {
		u, err := url.Parse(c)
		switch {
		case err != nil:
			return err
		case u.Scheme == "":
			return errors.New("scheme = ''")
		}
	}
	return nil
}

func (p SNetwork) GetFetchTimeout() time.Duration {
	return time.Duration(p.FFetchTimeoutMS) * time.Millisecond
}

func (p SNetwork) GetQueuePeriod() time.Duration {
	return time.Duration(p.FQueuePeriodMS) * time.Millisecond
}
