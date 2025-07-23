package app

import (
	"os"
	"testing"

	"github.com/number571/hidden-lake/build"
	"github.com/number571/hidden-lake/internal/projects/chat/pkg/settings"
	"github.com/number571/hidden-lake/internal/utils/flag"
)

var (
	tgFlags = flag.NewFlagsBuilder(
		flag.NewFlagBuilder("-v", "--version").
			WithDescription("print version of service"),
		flag.NewFlagBuilder("-h", "--help").
			WithDescription("print information about service"),
		flag.NewFlagBuilder("-p", "--path").
			WithDescription("set path to config, database files").
			WithDefinedValue("."),
		flag.NewFlagBuilder("-n", "--network").
			WithDescription("set network key of connections from build").
			WithDefinedValue(build.CDefaultNetwork),
	).Build()
)

const (
	tcTestdataPath = "./testdata"
	tcPathConfig   = settings.CPathDB
)

func testDeleteFiles(prefixPath string) {
	_ = os.RemoveAll(prefixPath + tcPathConfig)
}

func TestError(t *testing.T) {
	t.Parallel()

	str := "value"
	err := &SAppError{str}
	if err.Error() != errPrefix+str {
		t.Fatal("incorrect err.Error()")
	}
}

func TestInitApp(t *testing.T) {
	t.Parallel()

	testDeleteFiles(tcTestdataPath)
	defer testDeleteFiles(tcTestdataPath)

	if _, err := InitApp([]string{"--path", tcTestdataPath}, tgFlags); err != nil {
		t.Fatal(err)
	}
}
