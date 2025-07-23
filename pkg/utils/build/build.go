package build

import (
	"fmt"

	"github.com/number571/go-peer/pkg/logger"
)

const (
	cFileSettings = "hl-settings.yml"
	cFileNetworks = "hl-networks.yml"
)

func SetBuildByPath(pInputPath string) ([2]bool, error) {
	var (
		oks [2]bool
		err error
	)
	oks[0], err = setSettings(pInputPath, cFileSettings)
	if err != nil {
		return oks, err
	}
	oks[1], err = setNetworks(pInputPath, cFileNetworks)
	if err != nil {
		return oks, err
	}
	return oks, nil
}

func LogLoadedBuildFiles(appName string, logger logger.ILogger, oks [2]bool) {
	files := [2]string{cFileSettings, cFileNetworks}
	for i := range len(oks) {
		if oks[i] {
			logger.PushInfo(fmt.Sprintf("%s load %s build file;", appName, files[i]))
		}
	}
}
