package app

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/composite/pkg/app/config"
	"github.com/number571/hidden-lake/internal/composite/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"

	hla_tcp_app "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app"
	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"

	hlm_app "github.com/number571/hidden-lake/internal/applications/messenger/pkg/app"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"

	hlf_app "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/app"
	hlf_settings "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"

	hlr_app "github.com/number571/hidden-lake/internal/applications/remoter/pkg/app"
	hlr_settings "github.com/number571/hidden-lake/internal/applications/remoter/pkg/settings"

	hls_app "github.com/number571/hidden-lake/internal/service/pkg/app"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

func InitApp(pArgs []string, pFlags flag.IFlags) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(pFlags.Get("path").GetStringValue(pArgs), "/")

	cfgPath := filepath.Join(inputPath, settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil)
	if err != nil {
		return nil, errors.Join(ErrInitConfig, err)
	}

	runners, err := getRunners(cfg, pArgs, pFlags)
	if err != nil {
		return nil, errors.Join(ErrGetRunners, err)
	}

	return NewApp(cfg, runners), nil
}

func getRunners(pCfg config.IConfig, pArgs []string, pFlags flag.IFlags) ([]types.IRunner, error) {
	var (
		services = pCfg.GetServices()
		runners  = make([]types.IRunner, 0, len(services))
		mapsdupl = make(map[string]struct{}, len(services))
	)

	var (
		runner types.IRunner
		err    error
	)

	for _, sName := range services {
		if _, ok := mapsdupl[sName]; ok {
			return nil, ErrHasDuplicates
		}
		mapsdupl[sName] = struct{}{}

		switch sName {
		case hls_settings.CServiceFullName:
			runner, err = hls_app.InitApp(pArgs, pFlags)
		case hlm_settings.CServiceFullName:
			runner, err = hlm_app.InitApp(pArgs, pFlags)
		case hlf_settings.CServiceFullName:
			runner, err = hlf_app.InitApp(pArgs, pFlags)
		case hlr_settings.CServiceFullName:
			runner, err = hlr_app.InitApp(pArgs, pFlags)
		case hla_tcp_settings.CServiceFullName:
			runner, err = hla_tcp_app.InitApp(pArgs, pFlags)
		default:
			return nil, ErrUnknownService
		}
		if err != nil {
			return nil, err
		}

		runners = append(runners, runner)
	}

	return runners, nil
}
