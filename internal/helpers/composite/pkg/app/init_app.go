package app

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	build "github.com/number571/hidden-lake/build/environment"
	"github.com/number571/hidden-lake/internal/helpers/composite/pkg/app/config"
	hlc_settings "github.com/number571/hidden-lake/internal/helpers/composite/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
	std_logger "github.com/number571/hidden-lake/internal/utils/logger/std"

	hlk_app "github.com/number571/hidden-lake/internal/kernel/pkg/app"
	hlk_settings "github.com/number571/hidden-lake/internal/kernel/pkg/settings"

	hla_tcp_app "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/app"
	hla_tcp_settings "github.com/number571/hidden-lake/internal/adapters/tcp/pkg/settings"

	hla_http_app "github.com/number571/hidden-lake/internal/adapters/http/pkg/app"
	hla_http_settings "github.com/number571/hidden-lake/internal/adapters/http/pkg/settings"

	hla_https_app "github.com/number571/hidden-lake/internal/adapters/https/pkg/app"
	hla_https_settings "github.com/number571/hidden-lake/internal/adapters/https/pkg/settings"

	hla_meshtastic_app "github.com/number571/hidden-lake/internal/adapters/meshtastic/pkg/app"
	hla_meshtastic_settings "github.com/number571/hidden-lake/internal/adapters/meshtastic/pkg/settings"

	hls_notifier_app "github.com/number571/hidden-lake/internal/services/notifier/pkg/app"
	hls_notifier_settings "github.com/number571/hidden-lake/internal/services/notifier/pkg/settings"

	hls_filesharer_app "github.com/number571/hidden-lake/internal/services/filesharer/pkg/app"
	hls_filesharer_settings "github.com/number571/hidden-lake/internal/services/filesharer/pkg/settings"

	hls_pinger_app "github.com/number571/hidden-lake/internal/services/pinger/pkg/app"
	hls_pinger_settings "github.com/number571/hidden-lake/internal/services/pinger/pkg/settings"
)

var (
	mapInitApp = map[string]func([]string, flag.IFlags) (types.IRunner, error){
		hlk_settings.CAppShortName:            hlk_app.InitApp,
		hla_tcp_settings.CAppShortName:        hla_tcp_app.InitApp,
		hla_http_settings.CAppShortName:       hla_http_app.InitApp,
		hla_https_settings.CAppShortName:      hla_https_app.InitApp,
		hla_meshtastic_settings.CAppShortName: hla_meshtastic_app.InitApp,
		hls_notifier_settings.CAppShortName:   hls_notifier_app.InitApp,
		hls_filesharer_settings.CAppShortName: hls_filesharer_app.InitApp,
		hls_pinger_settings.CAppShortName:     hls_pinger_app.InitApp,
	}
)

func InitApp(pArgs []string, pFlags flag.IFlags) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(pFlags.Get("-p").GetStringValue(pArgs), "/")
	if err := os.MkdirAll(inputPath, 0700); err != nil {
		return nil, errors.Join(ErrMkdirPath, err)
	}

	okLoaded, err := build.SetBuildByPath(inputPath)
	if err != nil {
		return nil, errors.Join(ErrSetBuild, err)
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
	build.LogLoadedBuildFiles(hlc_settings.GetAppShortNameFMT(), stdfLogger, okLoaded)

	return NewApp(cfg, runners), nil
}

func getRunners(pCfg config.IConfig, pArgs []string, pFlags flag.IFlags) ([]types.IRunner, error) {
	var (
		applications = pCfg.GetApplications()
		runners      = make([]types.IRunner, 0, len(applications))
		mapsdupl     = make(map[string]struct{}, len(applications))
	)

	for _, app := range applications {
		if _, ok := mapsdupl[app]; ok {
			return nil, ErrHasDuplicates
		}
		mapsdupl[app] = struct{}{}

		initApp, ok := mapInitApp[app]
		if !ok {
			return nil, ErrUnknownService
		}

		runner, err := initApp(pArgs, pFlags)
		if err != nil {
			return nil, errors.Join(ErrInitApp, err)
		}

		runners = append(runners, runner)
	}

	return runners, nil
}
