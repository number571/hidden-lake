package app

import (
	"errors"
	"strings"

	"github.com/number571/go-peer/pkg/types"
	"github.com/number571/hidden-lake/internal/utils/build"
	"github.com/number571/hidden-lake/internal/utils/flag"
)

// initApp work with the raw data = read files, read args
func InitApp(pArgs []string, pFlags flag.IFlags) (types.IRunner, error) {
	inputPath := strings.TrimSuffix(pFlags.Get("-p").GetStringValue(pArgs), "/")
	_, err := build.SetBuildByPath(inputPath)
	if err != nil {
		return nil, errors.Join(ErrSetNetworks, err)
	}
	networkKey := pFlags.Get("-n").GetStringValue(pArgs)
	return NewApp(networkKey, inputPath), nil
}
