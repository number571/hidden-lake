package app

import (
	"path/filepath"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/go-peer/pkg/utils"
	"github.com/number571/hidden-lake/cmd/composite/internal/config"
	"github.com/number571/hidden-lake/cmd/composite/pkg/settings"
	"github.com/number571/hidden-lake/internal/flag"

	hla_chatingar_app "github.com/number571/hidden-lake/cmd/adapters/chatingar/pkg/app"
	hla_chatingar_settings "github.com/number571/hidden-lake/cmd/adapters/chatingar/pkg/settings"

	hla_common_app "github.com/number571/hidden-lake/cmd/adapters/common/pkg/app"
	hla_common_settings "github.com/number571/hidden-lake/cmd/adapters/common/pkg/settings"

	hlm_app "github.com/number571/hidden-lake/cmd/applications/messenger/pkg/app"
	hlm_settings "github.com/number571/hidden-lake/cmd/applications/messenger/pkg/settings"

	hlf_app "github.com/number571/hidden-lake/cmd/applications/filesharer/pkg/app"
	hlf_settings "github.com/number571/hidden-lake/cmd/applications/filesharer/pkg/settings"

	hlr_app "github.com/number571/hidden-lake/cmd/applications/remoter/pkg/app"
	hlr_settings "github.com/number571/hidden-lake/cmd/applications/remoter/pkg/settings"

	hle_app "github.com/number571/hidden-lake/cmd/helpers/encryptor/pkg/app"
	hle_settings "github.com/number571/hidden-lake/cmd/helpers/encryptor/pkg/settings"

	hll_app "github.com/number571/hidden-lake/cmd/helpers/loader/pkg/app"
	hll_settings "github.com/number571/hidden-lake/cmd/helpers/loader/pkg/settings"

	hlt_app "github.com/number571/hidden-lake/cmd/helpers/traffic/pkg/app"
	hlt_settings "github.com/number571/hidden-lake/cmd/helpers/traffic/pkg/settings"

	hls_app "github.com/number571/hidden-lake/cmd/service/pkg/app"
	hls_settings "github.com/number571/hidden-lake/cmd/service/pkg/settings"
)

func InitApp(
	pArgs []string,
	pDefaultPath string,
	pDefaultParallel uint64,
) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(flag.GetFlagValue(pArgs, "path", pDefaultPath), "/")

	cfgPath := filepath.Join(inputPath, settings.CPathYML)
	cfg, err := config.InitConfig(cfgPath, nil)
	if err != nil {
		return nil, utils.MergeErrors(ErrInitConfig, err)
	}

	runners, err := getRunners(
		cfg,
		pArgs,
		pDefaultPath,
		pDefaultParallel,
	)
	if err != nil {
		return nil, utils.MergeErrors(ErrGetRunners, err)
	}

	return NewApp(cfg, runners), nil
}

func getRunners(
	pCfg config.IConfig,
	pArgs []string,
	pDefaultPath string,
	pDefaultParallel uint64,
) ([]types.IRunner, error) {
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
			runner, err = hls_app.InitApp(pArgs, pDefaultPath, pDefaultParallel)
		case hle_settings.CServiceFullName:
			runner, err = hle_app.InitApp(pArgs, pDefaultPath, pDefaultParallel)
		case hlt_settings.CServiceFullName:
			runner, err = hlt_app.InitApp(pArgs, pDefaultPath)
		case hll_settings.CServiceFullName:
			runner, err = hll_app.InitApp(pArgs, pDefaultPath)
		case hlm_settings.CServiceFullName:
			runner, err = hlm_app.InitApp(pArgs, pDefaultPath)
		case hlf_settings.CServiceFullName:
			runner, err = hlf_app.InitApp(pArgs, pDefaultPath)
		case hlr_settings.CServiceFullName:
			runner, err = hlr_app.InitApp(pArgs, pDefaultPath)
		case hla_common_settings.CServiceFullName:
			runner, err = hla_common_app.InitApp(pArgs, pDefaultPath)
		case hla_chatingar_settings.CServiceFullName:
			runner, err = hla_chatingar_app.InitApp(pArgs, pDefaultPath)
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
