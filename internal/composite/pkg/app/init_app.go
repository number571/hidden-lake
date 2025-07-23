package app

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/composite/pkg/app/config"
	hlc_settings "github.com/number571/hidden-lake/internal/composite/pkg/settings"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"
	"github.com/number571/hidden-lake/pkg/utils/build"
	"github.com/number571/hidden-lake/pkg/utils/flag"

	hla_tcp_app "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app"
	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"

	hla_http_app "github.com/number571/hidden-lake/internal/adapters/http/pkg/app"
	hla_http_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"

	hls_messenger_app "github.com/number571/hidden-lake/internal/services/messenger/pkg/app"
	hls_messenger_settings "github.com/number571/hidden-lake/internal/services/messenger/pkg/settings"

	hls_filesharer_app "github.com/number571/hidden-lake/internal/services/filesharer/pkg/app"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"

	hls_remoter_app "github.com/number571/hidden-lake/internal/services/remoter/pkg/app"
	hls_remoter_settings "github.com/number571/hidden-lake/internal/services/remoter/pkg/settings"

	hls_pinger_app "github.com/number571/hidden-lake/internal/services/pinger/pkg/app"
	hls_pinger_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"

	hlk_app "github.com/number571/hidden-lake/internal/kernel/pkg/app"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"
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
	build.LogLoadedBuildFiles(hlc_settings.GetFmtAppName().Short(), stdfLogger, okLoaded)

	return NewApp(cfg, runners), nil
}

func getRunners(pCfg config.IConfig, pArgs []string, pFlags flag.IFlags) ([]types.IRunner, error) {
	var (
		applications = pCfg.GetApplications()
		runners      = make([]types.IRunner, 0, len(applications))
		mapsdupl     = make(map[string]struct{}, len(applications))
	)

	var (
		runner types.IRunner
		err    error
	)

	for _, app := range applications {
		if _, ok := mapsdupl[app]; ok {
			return nil, ErrHasDuplicates
		}
		mapsdupl[app] = struct{}{}

		switch app {
		case hlk_settings.CAppFullName:
			runner, err = hlk_app.InitApp(pArgs, pFlags)
		case hls_messenger_settings.CAppFullName:
			runner, err = hls_messenger_app.InitApp(pArgs, pFlags)
		case hls_filesharer_settings.CAppFullName:
			runner, err = hls_filesharer_app.InitApp(pArgs, pFlags)
		case hls_remoter_settings.CAppFullName:
			runner, err = hls_remoter_app.InitApp(pArgs, pFlags)
		case hls_pinger_settings.CAppFullName:
			runner, err = hls_pinger_app.InitApp(pArgs, pFlags)
		case hla_tcp_settings.CAppFullName:
			runner, err = hla_tcp_app.InitApp(pArgs, pFlags)
		case hla_http_settings.CAppFullName:
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
