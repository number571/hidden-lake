package app

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/composite/pkg/app/config"
	hlc_settings "github.com/number571/hidden-lake/internal/composite/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/build"
	"github.com/number571/hidden-lake/internal/utils/flag"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"

	hla_tcp_app "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app"
	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"

	hla_http_app "github.com/number571/hidden-lake/internal/adapters/http/pkg/app"
	hla_http_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"

	hlm_app "github.com/number571/hidden-lake/internal/applications/messenger/pkg/app"
	hlm_settings "github.com/number571/hidden-lake/internal/applications/messenger/pkg/settings"

	hlf_app "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/app"
	hlf_settings "github.com/number571/hidden-lake/internal/applications/filesharer/pkg/settings"

	hlr_app "github.com/number571/hidden-lake/internal/applications/remoter/pkg/app"
	hlr_settings "github.com/number571/hidden-lake/internal/applications/remoter/pkg/settings"

	hlp_app "github.com/number571/hidden-lake/internal/applications/pinger/pkg/app"
	hlp_settings "github.com/number571/hidden-lake/internal/applications/pinger/pkg/settings"

	hls_app "github.com/number571/hidden-lake/internal/service/pkg/app"
	hls_settings "github.com/number571/hidden-lake/internal/service/pkg/settings"
)

func InitApp(pArgs []string, pFlags flag.IFlags) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(pFlags.Get("-p").GetStringValue(pArgs), "/")
	okLoaded, err := build.SetBuildByPath(inputPath)
	if err != nil {
		return nil, errors.Join(ErrSetNetworks, err)
	}

	cfgPath := filepath.Join(inputPath, hlc_settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil, pFlags.Get("-n").GetStringValue(pArgs))
	if err != nil {
		return nil, errors.Join(ErrInitConfig, err)
	}

	runners, err := getRunners(cfg, pArgs, pFlags)
	if err != nil {
		return nil, errors.Join(ErrGetRunners, err)
	}

	stdfLogger := std_logger.NewStdLogger(cfg.GetLogging(), std_logger.GetLogFunc())
	build.LogLoadedBuildFiles(hlc_settings.GetServiceName(), stdfLogger, okLoaded)

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
		case hlp_settings.CServiceFullName:
			runner, err = hlp_app.InitApp(pArgs, pFlags)
		case hla_tcp_settings.CServiceFullName:
			runner, err = hla_tcp_app.InitApp(pArgs, pFlags)
		case hla_http_settings.CServiceFullName:
			runner, err = hla_http_app.InitApp(pArgs, pFlags)
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
